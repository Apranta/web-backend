package handlers

import (
	"net/http"
	"ugc-vote/config"
	"ugc-vote/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func GetEvent(c echo.Context) error {

	rows, err := config.App.DB.Query("select id, event, password from user")
	if err != nil {
		return err
	}
	defer rows.Close()

	var result []models.User

	for rows.Next() {
		var each = models.User{}
		var err = rows.Scan(&each.ID, &each.Username, &each.Password)

		if err != nil {
			return err
		}

		result = append(result, each)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
