package handlers

import (
	"fmt"
	"inventoryManagement/model"
	"inventoryManagement/views"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func SalesPutHandler(c echo.Context) error {
	usrCookie, usrCookieError := c.Cookie("usrtype")
	loginCookie, loginCookieError := c.Cookie("login")

	if usrCookieError != nil || loginCookieError != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tables, err := model.AllSales(model.DB)
	if err != nil {
		log.Printf("SalesPutHandler: %v", err)
		return nil
	}

	if loginCookie.Value == "yes" && usrCookie.Value == "admin" || usrCookie.Value == "user" || usrCookie.Value == "manager" {
		return views.SellTable(tables).Render(getParams(c))
	} else {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

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
	c.Response().Header().Add("HX-Trigger", "saleTotal")
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
	table := []model.TransactionProduct{}
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

func calculateTotalSaleHandler(c echo.Context) error {
	total, err := model.CalculateTotalSale()
	if err != nil {
		fmt.Printf("calculateTotalSaleHandler: %v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	totalString := fmt.Sprintf("Total: %v", strconv.FormatInt(int64(total), 10))
	return c.String(http.StatusOK, totalString)
}
