package handlers

import (
	"fmt"
	"inventoryManagement/model"
	"inventoryManagement/views"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func ProductsTableHandler(c echo.Context) error {
	usrCookie, usrCookieError := c.Cookie("usrtype")
	loginCookie, loginCookieError := c.Cookie("login")

	if usrCookieError != nil || loginCookieError != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tables, err := model.AllProducts(model.DB)
	if err != nil {
		return fmt.Errorf("ProductsTableHandler: %v", err)
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
	if marginErr != nil {
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
