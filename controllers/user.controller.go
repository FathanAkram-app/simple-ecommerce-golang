package controllers

import (
	"depmod/models"

	"github.com/labstack/echo"
)

// func FetchAllUser(c echo.Context) error {
// 	result := models.FetchAllUser()
// 	return c.JSON(http.StatusOK, result)
// }

func RegisterUser(c echo.Context) error {
	result := models.RegisterUser(
		c.FormValue("username"),
		c.FormValue("email"),
		c.FormValue("password"))

	return c.JSON(result.Status, result)
}

func LoginUser(c echo.Context) error {
	result := models.LoginUser(
		c.FormValue("email"),
		c.FormValue("password"))

	return c.JSON(result.Status, result)
}

func LogoutUser(c echo.Context) error {
	result := models.LogoutUser(
		c.Request().Header["Authorization"][0])

	return c.JSON(result.Status, result)
}
