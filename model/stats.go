package model

import (
	"fmt"
	"time"
)

type MoneySpentTime struct {
	Sales     int
	StartDate string
}
type MoneySpentDelta struct {
	Sales      int
	Difference float64
	Date       string
}

type MoneySpentProduct struct {
	Name     string
	Expenses int
}
type MoneySpentProductDelta struct {
	Name       string
	Expenses   int
	Difference float64
	Date       string
}
type SaleDelta struct {
	Sales      int
	Difference float64
	Date       string
}
type ProfitDelta struct {
	Profits            int
	DifferencePrevious float64
	Date               string
}
type MostSoldProductData struct {
	Name     string
	Quantity int
	Date     string
}

type MostProfitableProduct struct {
	Name    string
	Profits int
}

type ProductSaleDelta struct {
	Name            string
	CurrentSales    int
	MonthlyVariance float64
	Date            string
}

func ProductSalesInTimespan(productCode int64, startDate time.Time) ([]MostSoldProductData, error) {
	dateString := fmt.Sprintf("%v-%v-%v", startDate.Year(), startDate.Month(), startDate.Day())
	row, err := DB.Query(`SELECT p.producto_nombre, SUM(ps.pro_sal_cantidad), DATE(s.salida_fecha)
FROM producto_salida ps 
JOIN producto p ON ps.pro_sal_producto = p.producto_id
JOIN salida s ON ps.pro_sal_salida = s.salida_id 
WHERE p.producto_codigo = ? AND DATE(s.salida_fecha) > ? 
GROUP BY DATE(s.salida_fecha)`, productCode, dateString)

	if err != nil {
		return nil, fmt.Errorf("MostSoldProduct: Query: %v", err)
	}

	var data []MostSoldProductData

	defer row.Close()
	for row.Next() {
		var rowData MostSoldProductData
		if err := row.Scan(&rowData.Name, &rowData.Quantity, &rowData.Date); err != nil {
			return nil, fmt.Errorf("MostSoldProduct: Scan: %v", err)
		}
		data = append(data, rowData)
	}

	return data, nil
}

func MostSoldProductInTimespan(startDate time.Time) (string, int, error) {
	dateString := fmt.Sprintf("%v-%v-%v", startDate.Year(), startDate.Month(), startDate.Day())
	row := DB.QueryRow(`SELECT pn, MAX(sq.total)
FROM (SELECT p.producto_id, p.producto_nombre as pn, SUM(ps.pro_sal_cantidad) as total, DATE(s.salida_fecha)
FROM producto_salida ps 
JOIN producto p ON ps.pro_sal_producto = p.producto_id
JOIN salida s ON ps.pro_sal_salida = s.salida_id
WHERE DATE(s.salida_fecha) > ?
GROUP BY p.producto_id) as sq;`, dateString)

	var name string
	var quantity int

	if scanErr := row.Scan(&name, &quantity); scanErr != nil {
		return "", -1, fmt.Errorf("MostSoldProductInTimespan: Scan: %v", scanErr)
	}

	return name, quantity, nil
}

func ProductDeltaInTimespan(productCode int64, startDate time.Time) ([]ProductSaleDelta, error) {
	dateString := fmt.Sprintf("%v-%v-%v", startDate.Year(), startDate.Month(), startDate.Day())
	row, rErr := DB.Query(`SELECT pnom, ma, (ma*100)/mp, dt
FROM( SELECT pnom, sm ma, lag(sm) over(order by dt) mp, dt
FROM ( SELECT p.producto_nombre pnom, SUM(ps.pro_sal_cantidad) sm, DATE(s.salida_fecha) dt
FROM producto_salida ps 
JOIN producto p ON ps.pro_sal_producto = p.producto_id
JOIN salida s ON ps.pro_sal_salida = s.salida_id 
WHERE p.producto_codigo = ? AND DATE(s.salida_fecha) > ?
GROUP BY DATE(s.salida_fecha)) sq) sq2;`, productCode, dateString)

	if rErr != nil {
		return nil, fmt.Errorf("ProductDeltaInTimespan: Query: %v", rErr)
	}

	var data []ProductSaleDelta

	defer row.Close()
	for row.Next() {
		var rowData ProductSaleDelta
		if err := row.Scan(&rowData.Name, &rowData.CurrentSales, &rowData.MonthlyVariance, &rowData.Date); err != nil {
			return nil, fmt.Errorf("MostSoldProduct: Scan: %v", err)
		}
		data = append(data, rowData)
	}
	return data, nil
}

func MostProfitableProductInTimespan(startDate time.Time) (MostProfitableProduct, error) {
	dateString := fmt.Sprintf("%v-%v-%v", startDate.Year(), startDate.Month(), startDate.Day())

	row := DB.QueryRow(`SELECT pn, MAX(pfts)
FROM(
SELECT p.producto_nombre pn, SUM(ps.pro_sal_cantidad*ps.pro_sal_precio) pfts
FROM producto_salida ps 
JOIN producto p ON ps.pro_sal_producto = p.producto_id
JOIN salida s ON ps.pro_sal_salida = s.salida_id 
WHERE DATE(s.salida_fecha) > ?
GROUP BY p.producto_id
) sq;`, dateString)

	var product MostProfitableProduct

	if scanErr := row.Scan(&product.Name, &product.Profits); scanErr != nil {
		return MostProfitableProduct{}, fmt.Errorf("MostSoldProductInTimespan: Scan: %v", scanErr)
	}

	return product, nil
}

func ProductProfitsInTimespan(productCode int64, startDate time.Time) (MostProfitableProduct, error) {
	dateString := fmt.Sprintf("%v-%v-%v", startDate.Year(), startDate.Month(), startDate.Day())

	row := DB.QueryRow(`SELECT p.producto_nombre pn, SUM(ps.pro_sal_cantidad*ps.pro_sal_precio) pfts
FROM producto_salida ps 
JOIN producto p ON ps.pro_sal_producto = p.producto_id
JOIN salida s ON ps.pro_sal_salida = s.salida_id 
WHERE p.producto_codigo = ? AND DATE(s.salida_fecha) > ?`, productCode, dateString)

	var product MostProfitableProduct

	if scanErr := row.Scan(&product.Name, &product.Profits); scanErr != nil {
		return MostProfitableProduct{}, fmt.Errorf("MostSoldProductInTimespan: Scan: %v", scanErr)
	}

	return product, nil
}

func ProfitDeltaInTimespan(startDate time.Time) ([]ProfitDelta, error) {
	dateString := fmt.Sprintf("%v-%v-%v", startDate.Year(), startDate.Month(), startDate.Day())
	row, rErr := DB.Query(`SELECT ma, (ma*100)/mp, dt
FROM( SELECT sm ma, lag(sm) over(order by dt) mp, dt
FROM ( SELECT SUM(ps.pro_sal_cantidad*ps.pro_sal_precio) sm, DATE(s.salida_fecha) dt
FROM producto_salida ps 
JOIN producto p ON ps.pro_sal_producto = p.producto_id
JOIN salida s ON ps.pro_sal_salida = s.salida_id 
WHERE DATE(s.salida_fecha) > ?
GROUP BY DATE(s.salida_fecha)) sq) sq2;`, dateString)

	if rErr != nil {
		return nil, fmt.Errorf("ProfitDeltaInTimespan: Query: %v", rErr)
	}

	var data []ProfitDelta

	defer row.Close()
	for row.Next() {
		var rowData ProfitDelta
		if err := row.Scan(&rowData.Profits, &rowData.DifferencePrevious, &rowData.Date); err != nil {
			return nil, fmt.Errorf("MostSoldProduct: Scan: %v", err)
		}
		data = append(data, rowData)
	}
	return data, nil
}

func TotalSalesInTimespan(startDate time.Time) (int, error) {
	dateString := fmt.Sprintf("%v-%v-%v", startDate.Year(), startDate.Month(), startDate.Day())
	row := DB.QueryRow(`SELECT SUM(ps.pro_sal_cantidad)
FROM producto_salida ps 
JOIN salida s ON ps.pro_sal_salida = s.salida_id 
WHERE DATE(s.salida_fecha) > ?`, dateString)

	var data int

	if err := row.Scan(&data); err != nil {
		return -1, fmt.Errorf("MostSoldProduct: Scan: %v", err)
	}

	return data, nil
}

func TotalSaleDeltaInTimespan(startDate time.Time) ([]SaleDelta, error) {
	dateString := fmt.Sprintf("%v-%v-%v", startDate.Year(), startDate.Month(), startDate.Day())
	row, rErr := DB.Query(`SELECT ma, (ma*100)/mp, dt
FROM(SELECT sm ma, lag(sm) over(order by dt) mp, dt
FROM (SELECT SUM(ps.pro_sal_cantidad) sm, DATE(s.salida_fecha) dt
FROM producto_salida ps 
JOIN salida s ON ps.pro_sal_salida = s.salida_id 
WHERE DATE(s.salida_fecha) > ?
GROUP BY DATE(s.salida_fecha)) sq) sq2;`, dateString)

	if rErr != nil {
		return nil, fmt.Errorf("TotalSalesDeltaInTimespan: Query: %v", rErr)
	}

	var data []SaleDelta

	defer row.Close()
	for row.Next() {
		var rowData SaleDelta
		if err := row.Scan(&rowData.Sales, &rowData.Difference, &rowData.Date); err != nil {
			return nil, fmt.Errorf("TotalSaleDelta: Scan: %v", err)
		}
		data = append(data, rowData)
	}
	return data, nil
}

func MoneySpentProductTimespan(productCode int64, startDate time.Time) (MoneySpentProduct, error) {
	dateString := fmt.Sprintf("%v-%v-%v", startDate.Year(), startDate.Month(), startDate.Day())

	row := DB.QueryRow(`SELECT p.producto_nombre, SUM(pe.pro_ent_cantidad*pe.pro_ent_precio)
FROM producto_entrada pe 
JOIN producto p ON pe.pro_ent_pro_fk = p.producto_id
JOIN entrada e ON pe.pro_ent_ent_fk = e.entrada_id 
WHERE p.producto_codigo = ? AND DATE(e.entrada_fecha) > ?;`, productCode, dateString)

	var product MoneySpentProduct

	if scanErr := row.Scan(&product.Name, &product.Expenses); scanErr != nil {
		return MoneySpentProduct{}, fmt.Errorf("MoneySpentProductTimespan: Scan: %v", scanErr)
	}

	return product, nil
}

func MoneySpentProductDeltaTimespan(productCode int64, startDate time.Time) ([]MoneySpentProductDelta, error) {
	dateString := fmt.Sprintf("%v-%v-%v", startDate.Year(), startDate.Month(), startDate.Day())
	row, rErr := DB.Query(`SELECT pnom, ma, (ma*100)/mp, dt
FROM( SELECT pnom, eqt ma, lag(eqt) over(order by dt) mp, dt
FROM ( SELECT p.producto_nombre pnom, SUM(pe.pro_ent_cantidad*pe.pro_ent_precio) eqt, DATE(e.entrada_fecha) dt
FROM producto_entrada pe 
JOIN producto p ON pe.pro_ent_pro_fk = p.producto_id
JOIN entrada e ON  pe.pro_ent_ent_fk = e.entrada_id
WHERE p.producto_codigo = ? AND DATE(e.entrada_fecha) > ?
GROUP BY DATE(e.entrada_fecha)) sq) sq2;`, productCode, dateString)

	if rErr != nil {
		return nil, fmt.Errorf("MoneySpentProductDeltaTimespan: Query: %v", rErr)
	}

	var data []MoneySpentProductDelta

	defer row.Close()
	for row.Next() {
		var rowData MoneySpentProductDelta
		if err := row.Scan(&rowData.Name, &rowData.Expenses, &rowData.Difference, &rowData.Date); err != nil {
			return nil, fmt.Errorf("MoneySpentProductDeltaTimespan: Scan: %v", err)
		}
		data = append(data, rowData)
	}
	return data, nil
}

func MoneySpentTimespan(startDate time.Time) (MoneySpentTime, error) {
	dateString := fmt.Sprintf("%v-%v-%v", startDate.Year(), startDate.Month(), startDate.Day())

	row := DB.QueryRow(`SELECT SUM(pe.pro_ent_cantidad*pe.pro_ent_precio)
FROM producto_entrada pe 
JOIN entrada e ON pe.pro_ent_ent_fk = e.entrada_id 
WHERE DATE(e.entrada_fecha) > ?;`, dateString)

	var product MoneySpentTime

	if scanErr := row.Scan(&product.Sales); scanErr != nil {
		return MoneySpentTime{}, fmt.Errorf("MoneySpentProductTimespan: Scan: %v", scanErr)
	}

	return product, nil
}

func MoneySpentDeltaTimespan(startDate time.Time) ([]MoneySpentDelta, error) {
	dateString := fmt.Sprintf("%v-%v-%v", startDate.Year(), startDate.Month(), startDate.Day())
	row, rErr := DB.Query(`SELECT ma, (ma*100)/mp, dt
FROM( SELECT eqt ma, lag(eqt) over(order by dt) mp, dt
FROM ( SELECT SUM(pe.pro_ent_cantidad*pe.pro_ent_precio) eqt, DATE(e.entrada_fecha) dt
FROM producto_entrada pe 
JOIN entrada e ON  pe.pro_ent_ent_fk = e.entrada_id
WHERE DATE(e.entrada_fecha) > ?
GROUP BY DATE(e.entrada_fecha)) sq) sq2;`, dateString)

	if rErr != nil {
		return nil, fmt.Errorf("MoneySpentDeltaTimespan: Query: %v", rErr)
	}

	var data []MoneySpentDelta

	defer row.Close()
	for row.Next() {
		var rowData MoneySpentDelta
		if err := row.Scan(&rowData.Sales, &rowData.Difference, &rowData.Date); err != nil {
			return nil, fmt.Errorf("MoneySpentDeltaTimespan: Scan: %v", err)
		}
		data = append(data, rowData)
	}
	return data, nil
}
