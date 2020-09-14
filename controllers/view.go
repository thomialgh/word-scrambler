package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LoginPage -
func LoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

// Index -
func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "question.html", nil)
}
