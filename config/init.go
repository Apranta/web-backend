package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"web-backend-patal/validator"

	"github.com/fsnotify/fsnotify"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var (
	App *Application
)

type (
	Application struct {
		Name    string      `json:"name"`
		Port    string      `json:"port"`
		Version string      `json:"version"`
		ENV     string      `json:"env"`
		Config  viper.Viper `json:"prog_config"`
		DB      *gorm.DB    `json:"db"`
	}
)

// Initiate news instances
func init() {
	var err error
	App = &Application{}
	App.Name = "Palembang digital"
	App.Version = os.Getenv("APPVER")
	if err = App.LoadConfigs(); err != nil {
		log.Printf("Load config error : %v", err)
	}
	if err = App.DBinit(); err != nil {
		log.Printf("DB init error : %v", err)
	}

	// apply custom validator
	App.Port = App.Config.GetString("port")
	App.ENV = App.Config.GetString("env")
	v := validator.Validator{DB: App.DB}
	v.CustomValidatorRules()
}

func (x *Application) Close() (err error) {
	if err = x.DB.Close(); err != nil {
		return err
	}

	return nil
}

// Loads general configs
func (x *Application) LoadConfigs() error {
	var conf *viper.Viper

	conf = viper.New()
	conf.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	conf.AutomaticEnv()
	conf.SetConfigName("config")
	conf.AddConfigPath(".")
	conf.SetConfigType("yaml")
	if err := conf.ReadInConfig(); err != nil {
		return err
	}
	conf.WatchConfig()
	conf.OnConfigChange(func(e fsnotify.Event) {
		log.Println("App Config file changed %s:", e.Name)
		x.LoadConfigs()
	})
	x.Config = viper.Viper(*conf)
	return nil
}

// Loads DBinit configs
func (x *Application) DBinit() error {
	dbconf := x.Config.GetStringMap(fmt.Sprintf("database"))
	Cons := DBConfig{
		Adapter:  PostgresAdapter,
		Host:     dbconf["host"].(string),
		Port:     dbconf["port"].(string),
		Username: dbconf["username"].(string),
		Password: dbconf["password"].(string),
		Database: dbconf["database"].(string),
		SSL:      dbconf["sslmode"].(string),
	}
	Start(Cons)
	x.DB = DB
	return nil
}
