package handlers

import (
	"context"
	"fmt"
	"inventoryManagement/model"
	"inventoryManagement/views"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func getParams(c echo.Context) (context.Context, http.ResponseWriter) {
	return c.Request().Context(), c.Response().Writer
}

func LoginHandler(c echo.Context) error {
	return views.Login().Render(getParams(c))
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

func TablesPutHandler(c echo.Context) error {
	usrCookie, usrCookieError := c.Cookie("usrtype")
	loginCookie, loginCookieError := c.Cookie("login")

	if usrCookieError != nil || loginCookieError != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tables, err := model.AllProducts(model.DB)
	if err != nil {
		return fmt.Errorf("TablesPutHandler: %v", err)
	}

	if loginCookie.Value == "yes" && usrCookie.Value == "admin" || usrCookie.Value == "user" || usrCookie.Value == "manager" {
		switch usrCookie.Value {
		case "admin":
			return views.ProductsTableAdmin(tables).Render(getParams(c))
		case "manager":
			return views.ProductsTableManager(tables).Render(getParams(c))
		default:
			return views.ProductsTable(tables).Render(getParams(c))
		}
	} else {
		return c.Redirect(http.StatusSeeOther, "/login")
	}
}

func CreateProductHandler(c echo.Context) error {
	nombre := c.FormValue("nombre")
	codigo, err := strconv.ParseInt(c.FormValue("codigo"), 10, 32)
	if err != nil || codigo < 0 {
		c.Response().Header().Add("HX-Trigger", "invalidCodeNumber")
		fmt.Printf("CreateProductHandler InvalidCode: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}
	margen, err := strconv.ParseFloat(c.FormValue("margen"), 64)
	if err != nil {
		c.Response().Header().Add("HX-Trigger", "invalidMarginNumber")
		fmt.Printf("CreateProductHandler InvalidMarginNumber: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}
	precio, err := strconv.ParseFloat(c.FormValue("precio"), 64)
	if err != nil || precio <= 0 {
		c.Response().Header().Add("HX-Trigger", "invalidPriceNumber")
		fmt.Printf("CreateProductHandler InvalidPriceNumber: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	if err := model.AddNewProduct(c, nombre, codigo, margen, precio); err != nil {
		c.Response().Header().Add("HX-Trigger", "createProductError")
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}

func VentasPutHandler(c echo.Context) error {
	usrCookie, usrCookieError := c.Cookie("usrtype")
	loginCookie, loginCookieError := c.Cookie("login")

	if usrCookieError != nil || loginCookieError != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tables, err := model.AllSales(model.DB)
	if err != nil {
		log.Printf("VentasPutHandler: %v", err)
		return nil
	}

	if loginCookie.Value == "yes" && usrCookie.Value == "admin" || usrCookie.Value == "user" || usrCookie.Value == "manager" {
		return views.SellTable(tables).Render(getParams(c))
	} else {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

}
func ComprasPutHandler(c echo.Context) error {
	usrCookie, usrCookieError := c.Cookie("usrtype")
	loginCookie, loginCookieError := c.Cookie("login")

	if usrCookieError != nil || loginCookieError != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tables, err := model.AllPurchases(model.DB)
	if err != nil {
		log.Printf("ComprasPutHandler: %v", err)
		return nil
	}

	if loginCookie.Value == "yes" && usrCookie.Value == "admin" || usrCookie.Value == "user" || usrCookie.Value == "manager" {
		return views.PurchaseTable(tables).Render(getParams(c))
	} else {
		return c.Redirect(http.StatusSeeOther, "/login")
	}
}
func AdminPutHandler(c echo.Context) error {
	//TODO
	return echo.ErrNotImplemented
}
func UsuariosPutHandler(c echo.Context) error {
	return views.CrearUsuario().Render(getParams(c))
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

func CreateSessionHandler(c echo.Context) error {
	usr := c.FormValue("username")
	pwd := c.FormValue("password")

	if ok, err := model.Authenticate(usr, pwd, c); err != nil {
		c.Response().Header().Add("HX-Trigger", "loginFailed")
		return c.NoContent(http.StatusForbidden)
	} else if !ok {
		c.Response().Header().Add("HX-Trigger", "loginFailed")
		return c.NoContent(http.StatusForbidden)
	}
	c.Response().Header().Add("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
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
