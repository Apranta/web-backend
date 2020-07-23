package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/thedevsaddam/govalidator"
)

type Validator struct {
	DB *gorm.DB `json:"db"`
}

func (a *Validator) CustomValidatorRules() {
	// unique value on each column. format : []string{"unique:[table],[column],[exclude_column],[excluded_value]"}
	govalidator.AddCustomRule("unique", func(field string, rule string, message string, value interface{}) error {
		var (
			queryRow *gorm.DB
			total    int
		)

		query := `SELECT COUNT(*) as total FROM %s WHERE %s = ?`
		params := strings.Split(strings.TrimPrefix(rule, fmt.Sprintf("%s:", "unique")), ",")

		if len(params) == 2 {
			query = fmt.Sprintf(query, params[0], params[1])
			queryRow = a.DB.Raw(query, value)
		} else if len(params) == 4 {
			query += ` AND %s != ?`
			query = fmt.Sprintf(query, params[0], params[1], params[2])
			queryRow = a.DB.Raw(query, value, params[3])
		} else {
			return fmt.Errorf("Arguments not enough")
		}

		queryRow.Row().Scan(&total)

		if total > 0 {
			if message != "" {
				return errors.New(message)
			}
			if field == "telepon" {
				field = "Phone Number"
			}
			return fmt.Errorf("The %s has already been taken", field)
		}

		return nil
	})

	// validator for pagination
	govalidator.AddCustomRule("asc_desc", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		if strings.ToUpper(val) != "ASC" && strings.ToUpper(val) != "DESC" {
			return fmt.Errorf("The %s field must be contain word: asc or desc", field)
		}
		return nil
	})

	// validator for indonesia phone number
	govalidator.AddCustomRule("id_phonenumber", func(field string, rule string, message string, value interface{}) error {
		reg := regexp.MustCompile(`\+?([ -]?\d+)+|\(\d+\)([ -]\d+)`)
		if value == nil {
			return fmt.Errorf("no value")
		}
		val := value.(string)
		if !reg.MatchString(val) {
			return fmt.Errorf("The %s field is not a valid indonesia phone number", field)
		}
		return nil
	})
}
