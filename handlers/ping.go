package handlers

import (
	"fmt"
	"net/http"
	"time"

	"web-backend-patal/config"

	"github.com/labstack/echo/v4"
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

	if err = pingDB(); err != nil {
		info.DB = false
	}

	return c.JSON(http.StatusOK, info)
}

func pingDB() (err error) {
	return config.App.DB.DB().Ping()
}
