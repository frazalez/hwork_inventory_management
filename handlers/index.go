package handlers

import (
	"context"
	"inventoryManagement/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getParams(c echo.Context) (context.Context, http.ResponseWriter) {
	return c.Request().Context(), c.Response().Writer
}

func IndexEntranceHandler(c echo.Context) error {
	usrCookie, usrCookieError := c.Cookie("usrtype")
	loginCookie, loginCookieError := c.Cookie("login")

	if usrCookieError != nil || loginCookieError != nil {
		print(usrCookieError, loginCookieError)
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	if loginCookie.Value == "yes" && (usrCookie.Value == "admin" || usrCookie.Value == "user" || usrCookie.Value == "manager") {
		switch usrCookie.Value {
		case "admin":
			return views.Index("admin").Render(getParams(c))
		case "manager":
			return views.Index("manager").Render(getParams(c))
		case "user":
			return views.Index("user").Render(getParams(c))
		default:
			return c.Redirect(http.StatusSeeOther, "/login")
		}
	} else {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

}
