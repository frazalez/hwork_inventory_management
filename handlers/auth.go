package handlers

import (
	"fmt"
	"inventoryManagement/model"
	"inventoryManagement/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context) error {
	return views.Login().Render(getParams(c))
}
func CreateSessionHandler(c echo.Context) error {
	usr := c.FormValue("username")
	pwd := c.FormValue("password")

	if ok, err := model.Authenticate(usr, pwd, c); err != nil {
		fmt.Println(err)
		c.Response().Header().Add("HX-Trigger", "loginFailed")
		return c.NoContent(http.StatusForbidden)
	} else if !ok {
		fmt.Println("Login Failed")
		c.Response().Header().Add("HX-Trigger", "loginFailed")
		return c.NoContent(http.StatusForbidden)
	}
	c.Response().Header().Add("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

func SalirPostHandler(c echo.Context) error {

	c.SetCookie(&http.Cookie{
		Name:   "login",
		Value:  "no",
		Path:   "/",
		MaxAge: 0,
	})
	c.Response().Header().Add("HX-Redirect", "/login")

	return c.NoContent(http.StatusNoContent)

}
