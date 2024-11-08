package model

import (
	"crypto/subtle"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

var DB *sql.DB

func ConnectToDatabase() *sql.DB {

	cfg := mysql.Config{
		User:   "sorcerer",
		Passwd: "admin",
		Net:    "tcp",
		Addr:   "localhost",
		DBName: "inventario",
	}
	cfg.AllowNativePasswords = true

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected to database")
	return db
}

func AllProducts(db *sql.DB) ([]Producto, error) {
	var productos []Producto
	rows, err := db.Query("SELECT * FROM productos")
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
		productos = append(productos, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AllProducts: %v", err)
	}

	return productos, nil
}

func AllSales(db *sql.DB) ([]Producto_salida_join, error) {
	var productos []Producto_salida_join
	rows, err := db.Query(`select p.producto_nombre, p.producto_codigo, ps.pro_sal_precio, ps.pro_sal_cantidad , st.salida_tipo_nombre,  s.salida_fecha , u.usuario_nombre
	from productos p join producto_salida ps on p.producto_id = ps.pro_sal_producto
	JOIN salida_tipo st on ps.pro_sal_tipo = st.salida_tipo_id JOIN salida s
	on ps.pro_sal_venta = s.salida_id JOIN usuarios u on s.salida_usuario = u.usuario_id`)
	if err != nil {
		return nil, fmt.Errorf("AllSales: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var prod Producto_salida_join
		if err := rows.Scan(&prod.Nombre, &prod.Codigo, &prod.PrecioVenta, &prod.Cantidad, &prod.TipoSalida, &prod.Fecha, &prod.UsuarioSalida); err != nil {
			return nil, fmt.Errorf("AllSales: %v", err)
		}
		productos = append(productos, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AllSales: %v", err)
	}

	return productos, nil
}
func AllPurchases(db *sql.DB) ([]Producto_entrada_join, error) {
	var productos []Producto_entrada_join
	rows, err := db.Query(`select p.producto_nombre, p.producto_codigo, pe.pro_ent_precio , pe.pro_ent_cantidad , e.entrada_fecha , u.usuario_nombre
	from productos p join producto_entrada pe on p.producto_id = pe.pro_ent_pro_fk
	JOIN entrada e on pe.pro_ent_fk = e.entrada_id JOIN usuarios u on e.entrada_usuario = u.usuario_id`)
	if err != nil {
		return nil, fmt.Errorf("allPurchases: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var prod Producto_entrada_join
		if err := rows.Scan(&prod.Nombre, &prod.Codigo, &prod.Precio, &prod.Cantidad, &prod.Fecha, &prod.Usuario); err != nil {
			return nil, fmt.Errorf("allPurchases: %v", err)
		}
		productos = append(productos, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("allPurchases: %v", err)
	}

	return productos, nil
}
func AllUsers(db *sql.DB) ([]Producto_entrada_join, error) {
	return nil, nil
}

func Authenticate(username, password string, c echo.Context) (bool, error) {

	var realData UsrAuth

	if strings.ContainsAny(username, "\"`'") || strings.ContainsAny(password, "\"`'") {
		return false, fmt.Errorf("auth/invalidchars")
	}

	rows := DB.QueryRow(`SELECT u.usuario_nombre, u.usuario_pass, u.usuario_activo, p.privilegio_nombre FROM usuarios u
JOIN privilegios p on p.privilegio_id = u.usuario_privilegio
WHERE u.usuario_nombre = ?`, username)

	if err := rows.Scan(&realData.Usuario, &realData.Passwd, &realData.Activo, &realData.Privilegio); err != nil {
		return false, fmt.Errorf("auth/scan: %v", err)
	}

	if realData.Activo == 0 {
		return false, fmt.Errorf("auth/disableduser")
	}

	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("auth/rowserr: %v", err)
	}

	if subtle.ConstantTimeCompare([]byte(username), []byte(realData.Usuario)) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte(realData.Passwd)) == 1 {

		c.SetCookie(&http.Cookie{
			Name:     "usrtype",
			Value:    realData.Privilegio, //admin, user, manager validos.
			Path:     "/",
			MaxAge:   0,
			HttpOnly: true,
		})

		c.SetCookie(&http.Cookie{
			Name:     "login",
			Value:    "yes",
			Path:     "/",
			MaxAge:   0,
			HttpOnly: true,
		})

		return true, nil

	}
	return false, nil
}

func AddNewUser(c echo.Context, usr string, pwd string, priv string) error {
	verifyQuery := DB.QueryRow(`SELECT usuario_nombre FROM usuarios WHERE usuario_nombre = ?`, usr)
	var dupeUser string
	if verifyError := verifyQuery.Scan(&dupeUser); verifyError == nil {
		c.Response().Header().Add("HX-Trigger", "duplicateError")
		return fmt.Errorf("AddNewUser verify: %s already exists", dupeUser)
	}

	query := DB.QueryRow(`SELECT privilegio_id FROM privilegios WHERE privilegio_nombre = ?`, priv)
	var privInt int
	if queryError := query.Scan(&privInt); queryError != nil {
		c.Response().Header().Add("HX-Trigger", "invalidPrivilegeError")
		return fmt.Errorf("AddNewUser query: %v", queryError)
	}

	_, err := DB.Exec(`INSERT INTO usuarios (usuario_nombre, usuario_pass, usuario_activo, usuario_privilegio)
	VALUES (?,?,?,?)`, usr, pwd, 1, privInt)

	if err != nil {
		c.Response().Header().Add("HX-Trigger", "insertError")
		return fmt.Errorf("AddNewUser: %v", err)
	}

	return nil
}

func AddNewProduct(c echo.Context, name string, code int64, margin float64, price float64) error {
	{
		query := DB.QueryRow(`SELECT producto_id FROM productos WHERE producto_codigo = ?`, code)
		var row int
		if err := query.Scan(&row); err == nil {
			c.Response().Header().Add("HX-Trigger", "duplicateproduct")
			return fmt.Errorf("AddNewProduct DuplicateProduct %v already exists", row)
		}
	}

	{
		_, err := DB.Exec(`INSERT INTO productos
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

	DB.QueryRow(`update productos set producto_activado = 0 where producto_id = ?;`, id)
	return nil
}

func EnableProduct(id int64) error {
	DB.QueryRow(`update productos set producto_activado = 1 where producto_id = ?;`, id)
	return nil
}

func ModifyProduct(ogName string, name string, code int64, margin float64, price float64) error {
	print(name)
	print("\n")
	idQuery := DB.QueryRow(`SELECT producto_codigo FROM productos p WHERE p.producto_nombre = ?`, ogName)
	var id int
	if idQErr := idQuery.Scan(&id); idQErr != nil {
		println("error?")
		return fmt.Errorf("ModifyProduct IdQuery Scan: %v", idQErr)
	}
	fmt.Println(ogName, name, code, margin, price, id)
	print("\n")
	updateQuery := DB.QueryRow(`
    UPDATE productos p SET
    p.producto_nombre = ?,
    p.producto_codigo = ?,
    p.producto_margen = ?,
    p.producto_precio = ?
    WHERE p.producto_codigo = ?;`, name, code, margin, price, id)

	if UQErr := updateQuery.Scan(); UQErr != nil {
		fmt.Printf("ModifyProduct Update: %v \n", UQErr)
	}
	print("???2\n")
	return nil
}
