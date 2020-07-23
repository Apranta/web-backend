package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Login authenticates and signs the user into the system
// @Summary Log the user in
// @Description Authenticate the user
// @Tags accounts
// @ID accounts-login
// @Accept mpfd
// @Produce plain
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {string} string	"ok"
// @Router /api/accounts/login [post]
func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	log.Println(username, password)
	c.String(http.StatusOK, "ok")
	return nil
}
