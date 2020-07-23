package config

import (
<<<<<<< HEAD
	"database/sql"
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
=======
	"fmt"
	"log"
	"os"
	"strings"

	"web-backend-patal/validator"

	"github.com/fsnotify/fsnotify"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
>>>>>>> staging
)

var (
	App *Application
)

type (
	Application struct {
		Name   string  `json:"name"`
		Port   string  `json:"port"`
		Config Config  `json:"app_config"`
		DB     *sql.DB `json:"db"`
	}

	Config struct {
		Port        string `envconfig:"APPPORT"`
		JWT         string `envconfig:"JWT_SECRET"`
		DB_Host     string `envconfig:"DB_HOST"`
		DB_Username string `envconfig:"DB_USERNAME"`
		DB_Port     string `envconfig:"DB_PORT"`
		DB_Password string `envconfig:"DB_PASSWORD"`
		DB_Name     string `envconfig:"DB_NAME"`
		DB_SSL      string `envconfig:"DB_SSL"`
	}
)

// Initiate news instances
func init() {
	var err error
	App = &Application{}
	App.Name = "Palembang digital"
	if err = App.LoadConfigs(); err != nil {
		log.Printf("Load config error : %v", err)
	}
	if err = App.DBinit(); err != nil {
		log.Printf("DB init error : %v", err)
	}

}

func (x *Application) Close() (err error) {
	if err = x.DB.Close(); err != nil {
		return err
	}

	return nil
}

// Loads general configs
func (x *Application) LoadConfigs() error {

	err := envconfig.Process("patal", &x.Config)

	return err
}

// Loads DBinit configs
func (x *Application) DBinit() error {
	conf := x.Config

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", conf.DB_Username, conf.DB_Password, conf.DB_Host, conf.DB_Port, conf.DB_Name, conf.DB_SSL)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err // proper error handling instead of panic
	}
	// db.SetMaxOpenConns(dbconf["dbMaxConns"].(int))
	// db.SetMaxIdleConns(dbconf["dbMaxIdleConns"].(int))
	x.DB = db
	return nil
}
