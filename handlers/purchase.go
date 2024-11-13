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

func PurchasesPutHandler(c echo.Context) error {
	usrCookie, usrCookieError := c.Cookie("usrtype")
	loginCookie, loginCookieError := c.Cookie("login")

	if usrCookieError != nil || loginCookieError != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tables, err := model.AllPurchases(model.DB)
	if err != nil {
		log.Printf("PurchasesPutHandler: %v", err)
		return nil
	}

	if loginCookie.Value == "yes" && usrCookie.Value == "admin" || usrCookie.Value == "user" || usrCookie.Value == "manager" {
		return views.PurchaseTable(tables).Render(getParams(c))
	} else {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

}

func SmallPurchaseTableSearchHandler(c echo.Context) error {
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
func CreatePurchaseHandler(c echo.Context) error {
	barcode, bcErr := strconv.ParseInt(c.FormValue("barcode"), 10, 64)
	if bcErr != nil {
		fmt.Printf("CreatePurchase ParseIntCode: %v", bcErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}

	quantity, qErr := strconv.ParseInt(c.FormValue("quantity"), 10, 64)
	if qErr != nil {
		fmt.Printf("CreatePurchase ParseIntQuantity: %v", qErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}

	PurchaseType, sErr := strconv.ParseInt(c.FormValue("type"), 10, 64)
	if sErr != nil {
		fmt.Printf("CreatePurchase ParseIntType: %v", qErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}

	err := model.AddToPurchase(barcode, quantity, PurchaseType)
	if err != nil {
		fmt.Printf("CreatePurchase DatabaseError: %v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}

	newTable, ntErr := model.GetPurchaseTransactionTable()
	if ntErr != nil {
		fmt.Printf("GetPurchaseTransactionTable DatabaseError: %v", ntErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.PurchasesTransacTable(newTable).Render(getParams(c))
}

func TablesPurchasesHandler(c echo.Context) error {
	usrCookie, usrCookieError := c.Cookie("usrtype")
	loginCookie, loginCookieError := c.Cookie("login")

	if usrCookieError != nil || loginCookieError != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tables, err := model.AllPurchases(model.DB)
	if err != nil {
		log.Printf("TablesPurchasesHandler: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	for i := range tables {
		fmt.Printf("%v %v %v %v", tables[i].Nombre, tables[i].Codigo, tables[i].Precio, tables[i].Cantidad)
	}
	if loginCookie.Value == "yes" && usrCookie.Value == "admin" || usrCookie.Value == "user" || usrCookie.Value == "manager" {
		return views.PurchasesTableOnly(tables).Render(getParams(c))
	} else {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

}

func StartPurchaseHandler(c echo.Context) error {
	table := []model.TransactionPurchase{}
	if err := model.StartPurchase(); err != nil {
		fmt.Printf("StartPurchase DatabaseError: %v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.PurchasesTransacTable(table).Render(getParams(c))
}

func CompletePurchaseHandler(c echo.Context) error {
	usrcookie, cookieErr := c.Cookie("usrname")
	if cookieErr != nil {
		fmt.Printf("CompletePurchase CookieError: %v", cookieErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	username := usrcookie.Value
	if err := model.CompletePurchase(username); err != nil {
		fmt.Printf("CompletePurchase DatabaseError: %v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}

	allPurchases, err := model.AllPurchases(model.DB)
	if err != nil {
		fmt.Printf("CompletePurchase GetAllPurchases: %v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.PurchasesTableOnly(allPurchases).Render(getParams(c))
}
