package handlers

import (
	"fmt"
	"inventoryManagement/model"
	"inventoryManagement/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UsersPutHandler(c echo.Context) error {
	return views.CrearUsuario().Render(getParams(c))
}

func CreateUserHandler(c echo.Context) error {
	username := c.Request().FormValue("username")
	password := c.Request().FormValue("password")
	privilege := c.Request().FormValue("privilege")

	if err := model.AddNewUser(c, username, password, privilege); err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}
