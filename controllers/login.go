package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"word-scrambler/libs"
	"word-scrambler/models"
	"word-scrambler/pkg"

	"github.com/labstack/echo/v4"
)

// Login -
func Login(c echo.Context) error {
	var user struct {
		Username string `json:"username"`
	}
	if err := c.Bind(&user); err != nil {
		fmt.Println(err)
		return pkg.BadRequestResp(c, "invalid request format")
	}

	tx := pkg.DB.Begin()
	uins := models.User{
		Username: user.Username,
	}
	if err := tx.Create(&uins).Error; err != nil {
		tx.Rollback()
		return pkg.InternalErrResp(c, err)
	}
	tx.Commit()

	key := libs.GuuID(user.Username)

	cookie := &http.Cookie{
		Name:     "mysession",
		Value:    key,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(3 * time.Hour),
	}

	conn := pkg.GetConn()
	defer conn.Conn.Close()
	conn.Set(key, strconv.Itoa(int(uins.ID)), 3*time.Hour)
	c.SetCookie(cookie)
	return pkg.OKResponse(c)
}

// AuthHandler -
func AuthHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("mysession")
		if err != nil {
			return pkg.NotAuthorizeResp(c)
		}
		conn := pkg.GetConn()
		defer conn.Conn.Close()
		val, err := conn.Get(cookie.Value)
		if err != nil {
			return pkg.InternalErrResp(c, err)
		}

		if val == "" {
			return pkg.NotAuthorizeResp(c)
		}

		// u, err := models.GetUserID(pkg.DB, val)
		// if err != nil {
		// 	return pkg.InternalErrResp(c, err)
		// }

		userID, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return pkg.InternalErrResp(c, err)
		}

		c.Set("user_id", userID)

		return next(c)
	}
}

// Protected -
func Protected(c echo.Context) error {
	return pkg.OKResponse(c)
}
