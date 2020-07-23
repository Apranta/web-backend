package handlers

import (
	"fmt"
	"net/http"
	"time"

	"web-backend-patal/config"

<<<<<<< HEAD
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
=======
	"github.com/labstack/echo/v4"
>>>>>>> staging
)

// Info main type
type Info struct {
	Time string `json:"time"`
	DB   bool   `json:"database"`
}

var (
	err  error
	info Info
)

// ServiceInfo check service info
func ServiceInfo(c echo.Context) error {
	defer c.Request().Body.Close()

	info.Time = fmt.Sprintf("%v", time.Now().Format("2006-01-02T15:04:05"))
	info.DB = true

<<<<<<< HEAD
	if err = healthcheckDB(); err != nil {
=======
	if err = pingDB(); err != nil {
>>>>>>> staging
		info.DB = false
	}

	return c.JSON(http.StatusOK, info)
}

<<<<<<< HEAD
func healthcheckDB() (err error) {
	dbconf := config.App.Config.GetStringMap("database")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", dbconf["username"].(string), dbconf["password"].(string), dbconf["host"].(string), dbconf["port"].(string), dbconf["table"].(string))

	db, err := gorm.Open("mysql", connectionString)
	defer db.Close()
	return err
=======
func pingDB() (err error) {
	return config.App.DB.DB().Ping()
>>>>>>> staging
}
