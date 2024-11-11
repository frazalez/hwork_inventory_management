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
		fmt.Printf("CreateProductHandler AddNewProduct: %v", err)
		c.Response().Header().Add("HX-Trigger", "createProductError")
		return c.NoContent(http.StatusBadRequest)
	}

	c.Response().Header().Add("HX-Trigger", "refreshTable")
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

func DisableProductHandler(c echo.Context) error {
	product := c.Request().FormValue("productId")
	productId, err := strconv.ParseInt(product, 10, 64)
	if err != nil {
		fmt.Printf("DisableProduct ParseIntError: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	if err := model.DisableProduct(productId); err != nil {
		fmt.Printf("DisableProduct DatabaseError: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	c.Response().Header().Add("HX-Trigger", "refreshTable")
	return c.NoContent(http.StatusOK)
}

func EnableProductHandler(c echo.Context) error {
	product := c.Request().FormValue("productId")
	productId, err := strconv.ParseInt(product, 10, 64)
	if err != nil {
		fmt.Printf("EnableProduct ParseIntError: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	if err := model.EnableProduct(productId); err != nil {
		fmt.Printf("EnableProduct DatabaseError: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	c.Response().Header().Add("HX-Trigger", "refreshTable")
	return c.NoContent(http.StatusOK)
}

func ModifyProductFormHandler(c echo.Context) error {
	view := views.ModifyProduct(
		c.FormValue("name"),
		c.FormValue("code"),
		c.FormValue("margin"),
		c.FormValue("price"))
	if err := view.Render(getParams(c)); err != nil {
		fmt.Printf("ModifyProductFormHandler: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}
	print("RAAA \n")
	return nil
}

func ModifyProductHandler(c echo.Context) error {
	ogName := c.FormValue("ogName")
	name := c.FormValue("name")
	code, codeErr := strconv.ParseInt(c.Request().FormValue("code"), 10, 64)
	if codeErr != nil {
		fmt.Printf("ModifyProduct ParseIntCode: %v", codeErr)
		return c.NoContent(http.StatusBadRequest)
	}

	margin, marginErr := strconv.ParseFloat(c.FormValue("margin"), 64)
	if codeErr != nil {
		fmt.Printf("ModifyProduct ParseFloatMargin: %v", marginErr)
		return c.NoContent(http.StatusBadRequest)
	}

	price, priceErr := strconv.ParseFloat(c.FormValue("code"), 64)
	if priceErr != nil {
		fmt.Printf("ModifyProduct ParseFloatPrice: %v", priceErr)
		return c.NoContent(http.StatusBadRequest)
	}

	if err := model.ModifyProduct(ogName, name, code, margin, price); err != nil {
		fmt.Printf("ModifyProduct DBInsert: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	c.Response().Header().Add("HX-Trigger", "refreshTable")
	return nil

}

func SmallTableSearchHandler(c echo.Context) error {
	fmt.Println("test")
	code, codeErr := strconv.ParseInt(c.FormValue("barcode"), 10, 64)
	if codeErr != nil {
		fmt.Printf("SmallTableSearch ParseIntCode: %v", codeErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}

	data, dataErr := model.GetSmallTable(code)
	if dataErr != nil {
		fmt.Printf("SmallTableSearch DatabaseError: %v", dataErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}

	return views.ProductsTableSmall(data).Render(getParams(c))
}

// falta agregar sumatoria a la venta
// falta agregar boton de confirmacion a venta
// falta agregar mensaje que indica venta en proceso
func CreateSaleHandler(c echo.Context) error {
	barcode, bcErr := strconv.ParseInt(c.FormValue("barcode"), 10, 64)
	if bcErr != nil {
		fmt.Printf("CreateSale ParseIntCode: %v", bcErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}

	quantity, qErr := strconv.ParseInt(c.FormValue("quantity"), 10, 64)
	if qErr != nil {
		fmt.Printf("CreateSale ParseIntQuantity: %v", qErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}

	saleType, sErr := strconv.ParseInt(c.FormValue("type"), 10, 64)
	if sErr != nil {
		fmt.Printf("CreateSale ParseIntType: %v", qErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}

	err := model.AddToSale(barcode, quantity, saleType)
	if err != nil {
		fmt.Printf("CreateSale DatabaseError: %v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}

	newTable, ntErr := model.GetSaleTransactionTable()
	if ntErr != nil {
		fmt.Printf("GetSaleTransactionTable DatabaseError: %v", ntErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.SalesTransacTable(newTable).Render(getParams(c))
}

func TablesSalesHandler(c echo.Context) error {
	usrCookie, usrCookieError := c.Cookie("usrtype")
	loginCookie, loginCookieError := c.Cookie("login")

	if usrCookieError != nil || loginCookieError != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tables, err := model.AllSales(model.DB)
	if err != nil {
		log.Printf("TablesSalesHandler: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	for i := range tables {
		fmt.Printf("%v %v %v %v", tables[i].Nombre, tables[i].Codigo, tables[i].PrecioVenta, tables[i].Cantidad)
	}
	if loginCookie.Value == "yes" && usrCookie.Value == "admin" || usrCookie.Value == "user" || usrCookie.Value == "manager" {
		return views.SalesTableOnly(tables).Render(getParams(c))
	} else {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

}

func StartSaleHandler(c echo.Context) error {
	table := []model.Producto_salida_join{}
	if err := model.StartSale(); err != nil {
		fmt.Printf("StartSale DatabaseError: %v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.SalesTransacTable(table).Render(getParams(c))
}

func CompleteSaleHandler(c echo.Context) error {
	usrcookie, cookieErr := c.Cookie("usrname")
	if cookieErr != nil {
		fmt.Printf("CompleteSale CookieError: %v", cookieErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	username := usrcookie.Value
	if err := model.CompleteSale(username); err != nil {
		fmt.Printf("CompleteSale DatabaseError: %v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}

	allSales, err := model.AllSales(model.DB)
	if err != nil {
		fmt.Printf("CompleteSale GetAllSales: %v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.SalesTableOnly(allSales).Render(getParams(c))
}
