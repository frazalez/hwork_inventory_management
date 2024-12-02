package model

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
)

type Producto struct {
	Producto_id       int64
	Producto_nombre   string
	Producto_codigo   int64
	Producto_margen   float64
	Producto_precio   float64
	Producto_activado bool
}

type ProductSmall struct {
	Producto_codigo int64
	Producto_nombre string
	Producto_precio float64
}

type TransactionProduct struct {
	Name     string
	Code     int64
	Price    int64
	Quantity int64
	Type     int64
}

func AllProducts(db *sql.DB) ([]Producto, error) {
	var producto []Producto
	rows, err := db.Query("SELECT * FROM producto")
	if err != nil {
		return nil, fmt.Errorf("AllProducts: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var prod Producto
		if err := rows.Scan(&prod.Producto_id, &prod.Producto_nombre,
			&prod.Producto_codigo, &prod.Producto_margen,
			&prod.Producto_precio, &prod.Producto_activado); err != nil {

			return nil, fmt.Errorf("AllProducts: %v", err)
		}
		producto = append(producto, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AllProducts: %v", err)
	}

	return producto, nil
}
func AddNewProduct(c echo.Context, name string, code int64, margin float64, price float64) error {
	{
		query := DB.QueryRow(`SELECT producto_id FROM producto WHERE producto_codigo = ?`, code)
		var row int
		if err := query.Scan(&row); err == nil {
			c.Response().Header().Add("HX-Trigger", "duplicateproduct")
			return fmt.Errorf("AddNewProduct DuplicateProduct %v already exists", row)
		}
	}

	{
		_, err := DB.Exec(`INSERT INTO producto
		(producto_nombre, producto_codigo, producto_margen,
		producto_precio, producto_activado) VALUES
		(?, ?, ?, ?, ?)`, name, code, margin, price, 1)
		if err != nil {
			c.Response().Header().Add("HX-Trigger", "insertionerror")
			return fmt.Errorf("AddNewProduct InsertInto %v", err)
		}
	}

	return nil
}
func DisableProduct(id int64) error {

	DB.QueryRow(`update producto set producto_activado = 0 where producto_id = ?;`, id)
	return nil
}

func EnableProduct(id int64) error {
	DB.QueryRow(`update producto set producto_activado = 1 where producto_id = ?;`, id)
	return nil
}

func ModifyProduct(ogName string, name string, code int64, margin float64, price float64) error {
	idQuery := DB.QueryRow(`SELECT producto_codigo FROM producto p WHERE p.producto_nombre = ?`, ogName)
	var id int
	if idQErr := idQuery.Scan(&id); idQErr != nil {
		return fmt.Errorf("ModifyProduct IdQuery Scan: %v", idQErr)
	}
	fmt.Println(ogName, name, code, margin, price, id)
	updateQuery := DB.QueryRow(`
    UPDATE producto p SET
    p.producto_nombre = ?,
    p.producto_codigo = ?,
    p.producto_margen = ?,
    p.producto_precio = ?
    WHERE p.producto_codigo = ?;`, name, code, margin, price, id)

	if UQErr := updateQuery.Scan(); UQErr != nil {
		fmt.Printf("ModifyProduct Update: %v \n", UQErr)
	}
	return nil
}

func GetSmallTable(code int64) ([]ProductSmall, error) {
	codeStr := fmt.Sprintf("%s%v%s", "%", code, "%")
	fmt.Println(codeStr)
	rows, err := DB.Query(`SELECT producto_codigo, producto_nombre, producto_precio FROM producto WHERE producto_codigo LIKE ?`, codeStr)
	var products []ProductSmall
	if err != nil {
		return nil, fmt.Errorf("GetSmallTable: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var prod ProductSmall
		if err := rows.Scan(&prod.Producto_codigo, &prod.Producto_nombre, &prod.Producto_precio); err != nil {
			return nil, fmt.Errorf("GetSmallTable Scan: %v", err)
		}
		fmt.Println(prod.Producto_codigo, prod.Producto_nombre, prod.Producto_precio)
		products = append(products, prod)
	}

	return products, nil
}
