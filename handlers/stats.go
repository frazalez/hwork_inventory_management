package handlers

import (
	"fmt"
	"inventoryManagement/model"
	"inventoryManagement/views"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func getDateAndCode(c echo.Context) (time.Time, int64, error) {

	date, dErr := time.Parse(time.DateOnly, c.FormValue("date"))
	if dErr != nil {
		fmt.Printf("CalculateSales %v", dErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return time.Time{}, -1, c.NoContent(http.StatusBadRequest)
	}
	code, err := strconv.ParseInt(c.FormValue("code"), 10, 64)
	if err != nil {
		fmt.Printf("CalculateSales %v", dErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return time.Time{}, -1, c.NoContent(http.StatusBadRequest)
	}

	return date, code, nil
}

func TableSearchHandler(c echo.Context) error {
	code, codeErr := strconv.ParseInt(c.FormValue("code"), 10, 64)
	if codeErr != nil {
		fmt.Printf("TableSearch ParseIntCode: %v", codeErr)
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

func StatsMainScreenHandler(c echo.Context) error {
	return views.ViewStats().Render(getParams(c))
}

func ViewMostSoldProductHandler(c echo.Context) error {
	date, dErr := time.Parse(time.DateOnly, c.FormValue("date"))
	if dErr != nil {
		fmt.Printf("ViewMostSold %v", dErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	name, qtty, err := model.MostSoldProductInTimespan(date)
	if err != nil {
		fmt.Printf("ViewMostSold %v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.MostSoldTable(name, qtty).Render(getParams(c))
}

func CalculateSalesOfProductInTimespan(c echo.Context) error {
	date, code, dcErr := getDateAndCode(c)
	if dcErr != nil {
		fmt.Printf("calculateSalesofProduct %v", dcErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	mostSoldProductData, err := model.ProductSalesInTimespan(code, date)
	if err != nil {
		fmt.Printf("mostSoldProductData %v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.SalesTimespan(mostSoldProductData).Render(getParams(c))
}

func CalculateProductDeltaInTimespan(c echo.Context) error {
	date, code, dcErr := getDateAndCode(c)
	if dcErr != nil {
		fmt.Printf("%v", dcErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	productSaleDelta, err := model.ProductDeltaInTimespan(code, date)
	if err != nil {

		fmt.Printf("%v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.DeltaSalesTimespan(productSaleDelta).Render(getParams(c))
}

func CalculateMostProfitableProductInTimespan(c echo.Context) error {
	date, _, dcErr := getDateAndCode(c)
	if dcErr != nil {
		fmt.Printf("%v", dcErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	mostProfitableProduct, err := model.MostProfitableProductInTimespan(date)
	if err != nil {

		fmt.Printf("%v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.ProductHighestIncome(mostProfitableProduct).Render(getParams(c))
}

func CalculateProfitsOfProductInTimespan(c echo.Context) error {
	date, code, dcErr := getDateAndCode(c)
	if dcErr != nil {
		fmt.Printf("%v", dcErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	mostProfitable, err := model.ProductProfitsInTimespan(code, date)
	if err != nil {
		fmt.Printf("%v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.IncomeToDate(mostProfitable).Render(getParams(c))
}

func CalculateProfitDeltaInTimespan(c echo.Context) error {
	date, _, dcErr := getDateAndCode(c)
	if dcErr != nil {
		fmt.Printf("%v", dcErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	profitDelta, err := model.ProfitDeltaInTimespan(date)
	if err != nil {

		fmt.Printf("%v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.IncomeVar(profitDelta).Render(getParams(c))
}

func CountSalesInTimespan(c echo.Context) error {
	date, _, dcErr := getDateAndCode(c)
	if dcErr != nil {
		fmt.Printf("%v", dcErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	sales, err := model.TotalSalesInTimespan(date)
	if err != nil {

		fmt.Printf("%v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.SalesToDate(sales).Render(getParams(c))
}

func CalculateSaleAmountDeltaTimespan(c echo.Context) error {
	date, code, dcErr := getDateAndCode(c)
	if dcErr != nil {
		fmt.Printf("%v", dcErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	productSaleDelta, err := model.ProductDeltaInTimespan(code, date)
	if err != nil {

		fmt.Printf("%v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.DeltaSalesTimespan(productSaleDelta).Render(getParams(c))
}

func CalculateMoneySpentProduct(c echo.Context) error {
	date, code, dcErr := getDateAndCode(c)
	if dcErr != nil {
		fmt.Printf("%v", dcErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	moneySpentProduct, err := model.MoneySpentProductTimespan(code, date)
	if err != nil {

		fmt.Printf("%v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.PExpense(moneySpentProduct).Render(getParams(c))
}

func CalculateDeltaMoneySpentProduct(c echo.Context) error {
	date, code, dcErr := getDateAndCode(c)
	if dcErr != nil {
		fmt.Printf("%v", dcErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	moneySpentDelta, err := model.MoneySpentProductDeltaTimespan(code, date)
	if err != nil {

		fmt.Printf("%v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	return views.PExpenseVar(moneySpentDelta).Render(getParams(c))
}

func CalculateMoneySpentTimespan(c echo.Context) error {
	date, _, dcErr := getDateAndCode(c)
	if dcErr != nil {
		fmt.Printf("%v", dcErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	moneySpent, err := model.MoneySpentTimespan(date)
	if err != nil {
		fmt.Printf("%v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)

	}
	return views.Expense(moneySpent).Render(getParams(c))
}

func CalculateMoneySpentDeltaTimespan(c echo.Context) error {
	date, _, dcErr := getDateAndCode(c)
	if dcErr != nil {
		fmt.Printf("%v", dcErr)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}
	moneySpentDelta, err := model.MoneySpentDeltaTimespan(date)
	if err != nil {

		fmt.Printf("%v", err)
		c.Response().Header().Add("HX-Trigger", "cancel")
		return c.NoContent(http.StatusBadRequest)
	}

	return views.ExpenseVar(moneySpentDelta).Render(getParams(c))
}
