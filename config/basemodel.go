package config

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	// import mysql, postgres, and pq
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

// used constants
const (
	MysqlAdapter    string = "mysql_adapter"
	PostgresAdapter string = "postgres_adapter"
)

type (
	// DBConfig contains db configs
	DBConfig struct {
		Adapter        string
		Host           string
		Port           string
		Username       string
		Password       string
		Database       string
		Timezone       string
		Maxlifetime    int
		IdleConnection int
		OpenConnection int
		SSL            string
		Logmode        bool
	}

	// BaseModel will be used as foundation of all models
	BaseModel struct {
		ID        uint64    `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key,column:id"`
		CreatedAt time.Time `json:"created_at" gorm:"column:created_at;DEFAULT:CURRENT_TIMESTAMP" sql:"DEFAULT:CURRENT_TIMESTAMP"`
		UpdateAt  time.Time `json:"update_at" gorm:"column:update_at" sql:"sql:"DEFAULT:CURRENT_TIMESTAMP"`
	}

	// DBFunc gorm trx function
	DBFunc func(tx *gorm.DB) error

	// PagedFindResult pagination format
	PagedFindResult struct {
		TotalData   int         `json:"total_data"`   // matched datas
		Rows        int         `json:"rows"`         // shown datas per page
		CurrentPage int         `json:"current_page"` // current page
		LastPage    int         `json:"last_page"`
		From        int         `json:"from"` // offset, starting index of data shown in current page
		To          int         `json:"to"`   // last index of data shown in current page
		Data        interface{} `json:"data"`
	}

	// CompareFilter filter used for comparing 2 values
	CompareFilter struct {
		Value1 interface{} `json:"value1"`
		Value2 interface{} `json:"value2"`
	}
)

// DB is a connected db instance
var DB *gorm.DB

// Start set basemodel config
func Start(conf DBConfig) {
	switch conf.Adapter {
	case MysqlAdapter:
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)
		if err := DBConnect("mysql", connectionString, conf); err != nil {
			panic(err)
		}
	case PostgresAdapter:
		connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database, conf.SSL)
		if err := DBConnect("postgres", connectionString, conf); err != nil {
			panic(err)
		}
	}
}

// Close DB connection
func Close() {
	DB.Close()
}

// DBConnect connects to db instance
func DBConnect(gormAdapter string, connectionString string, conf DBConfig) (err error) {
	switch conf.Adapter {
	case MysqlAdapter:
		DB, err = gorm.Open("mysql", connectionString)
	case PostgresAdapter:
		DB, err = gorm.Open(gormAdapter, connectionString)
	}
	if err != nil {
		return err
	}

	if err = DB.DB().Ping(); err != nil {
		return err
	}

	DB.LogMode(conf.Logmode)

	DB.Exec("SET time_zone = 'Asia/Jakarta'")
	DB.DB().SetConnMaxLifetime(time.Second * time.Duration(conf.Maxlifetime))
	DB.DB().SetMaxIdleConns(conf.IdleConnection)
	DB.DB().SetMaxOpenConns(conf.OpenConnection)
	DB.SingularTable(true)

	return nil
}

// WithinTransaction sql execute helper
func WithinTransaction(fn DBFunc) (err error) {
	tx := DB.Begin()
	defer tx.Commit()
	err = fn(tx)

	return err
}

// Create creates row based on model
func Create(i interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		if !DB.NewRecord(i) {
			return fmt.Errorf("cannot create row. not a new record")
		}
		if err = tx.Create(i).Error; err != nil {
			tx.Rollback()
			return err
		}
		return err
	})
}

// Save updates row based on model
func Save(i interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		if DB.NewRecord(i) {
			return fmt.Errorf("cannot save row. it is a new record")
		}
		if err = tx.Save(i).Error; err != nil {
			tx.Rollback()
			return err
		}
		return err
	})
}

// Delete removes row in db based on model.
func Delete(i interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		if err = tx.Delete(i).Error; err != nil {
			tx.Rollback()
			return err
		}
		return err
	})
}

// FirstOrCreate gets first matched record, or create a new one
func FirstOrCreate(i interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		if err = tx.FirstOrCreate(i).Error; err != nil {
			tx.Rollback()
			return err
		}
		return err
	})
}

// FindbyID finds row by id.
func FindbyID(i interface{}, id int) (err error) {
	return WithinTransaction(func(tx *gorm.DB) error {
		if err = tx.Last(i, id).Error; err != nil {
			tx.Rollback()
			return err
		}
		return err
	})
}

// SingleFindFilter finds by filter
func SingleFindFilter(i interface{}, filter interface{}) (err error) {
	query := DB // clone db connection

	query = conditionQuery(query, filter)

	err = query.Last(i).Error

	return err
}

// FindFilter finds by filter. limit 0 to find all
func FindFilter(i interface{}, order []string, sort []string, limit int, offset int, filter interface{}) (interface{}, error) {
	query := DB // clone db connection

	query = conditionQuery(query, filter)
	query = orderSortQuery(query, order, sort)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(i).Error

	return i, err
}

// PagedFindFilter same with FindFilter but with pagination
func PagedFindFilter(i interface{}, page int, rows int, order []string, sort []string, filter interface{}, field []string, allfieldcondition ...string) (result PagedFindResult, err error) {

	if page <= 0 {
		page = 1
	}

	query := DB

	query = conditionQuery(query, filter)
	query = orderSortQuery(query, order, sort)

	temp := query
	var totalRows int

	temp.Find(i).Count(&totalRows)

	if len(field) > 0 {
		query = query.Select(field)
	}

	var (
		offset   int
		lastPage int
	)

	if rows > 0 {
		offset = (page * rows) - rows
		lastPage = int(math.Ceil(float64(totalRows) / float64(rows)))

		query = query.Limit(rows).Offset(offset)
	}

	query.Find(i)

	result = PagedFindResult{
		TotalData:   totalRows,
		Rows:        rows,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + rows,
		Data:        i,
	}

	return result, err
}

func GetAll(i interface{}, filter interface{}) (data interface{}, err error) {

	db := DB

	// filtering
	refFilter := reflect.ValueOf(filter).Elem()
	refType := refFilter.Type()
	for x := 0; x < refFilter.NumField(); x++ {
		field := refFilter.Field(x)
		if field.Interface() != "" {
			switch refType.Field(x).Tag.Get("condition") {
			default:
				db = db.Where(fmt.Sprintf("%s = ?", refType.Field(x).Tag.Get("json")), field.Interface())
			case "LIKE":
				db = db.Where(fmt.Sprintf("%s %s ?", refType.Field(x).Tag.Get("json"), refType.Field(x).Tag.Get("condition")), "%"+field.Interface().(string)+"%")
			}
		}
	}

	db.Find(i)
	return i, err
}

func conditionQuery(query *gorm.DB, filter interface{}) *gorm.DB {
	refFilter := reflect.ValueOf(filter).Elem()
	refType := refFilter.Type()
	for x := 0; x < refFilter.NumField(); x++ {
		field := refFilter.Field(x)
		// check if empty
		if !reflect.DeepEqual(field.Interface(), reflect.Zero(reflect.TypeOf(field.Interface())).Interface()) {
			con := strings.Split(refType.Field(x).Tag.Get("condition"), ",")
			tags := parseTag(refType.Field(x).Tag.Get("condition"))
			switch con[0] {
			default:
				format := fmt.Sprintf("%s IN (?)", refType.Field(x).Tag.Get("json"))
				if tags.Contains("optional") {
					query = query.Or(format, field.Interface())
				} else {
					query = query.Where(format, field.Interface())
				}
			case "LIKE":
				format := fmt.Sprintf("LOWER(%s) %s ?", refType.Field(x).Tag.Get("json"), con[0])
				field := "%" + strings.ToLower(field.Interface().(string)) + "%"
				if tags.Contains("optional") {
					query = query.Or(format, field)
				} else {
					query = query.Where(format, field)
				}
			case "BETWEEN":
				if values, ok := field.Interface().(CompareFilter); ok && values.Value1 != "" {
					format := fmt.Sprintf("%s %s ? %s ?", refType.Field(x).Tag.Get("json"), con[0], "AND")
					if tags.Contains("optional") {
						query = query.Or(format, values.Value1, values.Value2)
					} else {
						query = query.Where(format, values.Value1, values.Value2)
					}
				}
			case "OR":
				var e []string
				for _, v := range field.Interface().([]string) {
					e = append(e, refType.Field(x).Tag.Get("json")+" = '"+v+"'")
				}
				if tags.Contains("optional") {
					query = query.Or(strings.Join(e, " OR "))
				} else {
					query = query.Where(strings.Join(e, " OR "))
				}
			}
		}
	}

	return query
}

func orderSortQuery(query *gorm.DB, order []string, sort []string) *gorm.DB {
	for k, v := range order {
		q := v
		if len(sort) > k {
			value := sort[k]
			if strings.ToUpper(value) == "ASC" || strings.ToUpper(value) == "DESC" {
				q = v + " " + strings.ToUpper(value)
			}
		}
		query = query.Order(q)
	}

	return query
}
