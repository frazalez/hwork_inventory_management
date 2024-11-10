package model

import (
	"crypto/subtle"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

var DB *sql.DB

func ConnectToDatabase() *sql.DB {

	cfg := mysql.Config{
		User:   "root",
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

func AllSales(db *sql.DB) ([]Producto_salida_join, error) {
	var producto []Producto_salida_join
	rows, err := db.Query(`select 
    p.producto_nombre, p.producto_codigo, 
    ps.pro_sal_precio, ps.pro_sal_cantidad , 
    st.salida_tipo_nombre,  s.salida_fecha , u.usuario_nombre
	from 
    producto p 
	join 
    producto_salida ps on p.producto_id = ps.pro_sal_producto
	JOIN 
    salida s on ps.pro_sal_salida = s.salida_id
	JOIN
    salida_tipo st on s.salida_tipo = st.salida_tipo_id
	JOIN
    usuario u on s.salida_usuario = u.usuario_id;`)
	if err != nil {
		return nil, fmt.Errorf("AllSales: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var prod Producto_salida_join
		if err := rows.Scan(&prod.Nombre, &prod.Codigo, &prod.PrecioVenta, &prod.Cantidad, &prod.TipoSalida, &prod.Fecha, &prod.UsuarioSalida); err != nil {
			return nil, fmt.Errorf("AllSales: %v", err)
		}
		fmt.Println(prod.Cantidad, prod.Codigo, prod.Nombre)
		producto = append(producto, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AllSales: %v", err)
	}

	return producto, nil
}
func AllPurchases(db *sql.DB) ([]Producto_entrada_join, error) {
	var producto []Producto_entrada_join
	rows, err := db.Query(`select p.producto_nombre, p.producto_codigo, pe.pro_ent_precio , pe.pro_ent_cantidad , e.entrada_fecha , u.usuario_nombre
	from producto p join producto_entrada pe on p.producto_id = pe.pro_ent_pro_fk
	JOIN entrada e on pe.pro_ent_fk = e.entrada_id JOIN usuario u on e.entrada_usuario = u.usuario_id`)
	if err != nil {
		return nil, fmt.Errorf("allPurchases: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var prod Producto_entrada_join
		if err := rows.Scan(&prod.Nombre, &prod.Codigo, &prod.Precio, &prod.Cantidad, &prod.Fecha, &prod.Usuario); err != nil {
			return nil, fmt.Errorf("allPurchases: %v", err)
		}
		producto = append(producto, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("allPurchases: %v", err)
	}

	return producto, nil
}
func AllUsers(db *sql.DB) ([]Producto_entrada_join, error) {
	return nil, nil
}

func Authenticate(username, password string, c echo.Context) (bool, error) {

	var realData UsrAuth

	if strings.ContainsAny(username, "\"`'") || strings.ContainsAny(password, "\"`'") {
	}
	
	rows := DB.QueryRow(`SELECT u.usuario_nombre, u.usuario_psswd, u.usuario_activo, p.privilegio_nombre FROM usuario u
	JOIN privilegio p on p.privilegio_id = u.usuario_privilegio
WHERE u.usuario_nombre = ?`, username)

if err := rows.Scan(&realData.Usuario, &realData.Passwd, &realData.Activo, &realData.Privilegio); err != nil {
	return false, fmt.Errorf("auth/scan: %v", err)
}

fmt.Printf("%v %v %v %v", realData.Usuario, realData.Passwd, realData.Activo, realData.Privilegio)

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
	verifyQuery := DB.QueryRow(`SELECT usuario_nombre FROM usuario WHERE usuario_nombre = ?`, usr)
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

	_, err := DB.Exec(`INSERT INTO usuario (usuario_nombre, usuario_pass, usuario_activo, usuario_privilegio)
	VALUES (?,?,?,?)`, usr, pwd, 1, privInt)

	if err != nil {
		c.Response().Header().Add("HX-Trigger", "insertError")
		return fmt.Errorf("AddNewUser: %v", err)
	}

	return nil
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

func GetSmallTable (code int64) ([]ProductSmall, error) {
	codeStr := fmt.Sprintf("%s%v%s", "%",code,"%")
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

func StartSale() error {
	r, err := DB.Exec(
		`DROP TABLE IF EXISTS transaccion_salida_producto;
		`)

	r2, err2 := DB.Exec(`CREATE TABLE transaccion_salida_producto (
			id int(10) primary key auto_increment,
			name varchar(40),
			code int(20),
			price int(10),
			quantity int(10),
			type int(1),
			foreign key (type) references salida_tipo(salida_tipo_id)		
		);`)

	if err != nil {
		return fmt.Errorf("StartSale: %v", err)
	}
	if err2 != nil {
		return fmt.Errorf("StartSale: %v", err2)
	}
	fmt.Println(r)
	fmt.Println(r2)
	return nil
}

func AddToSale(barcode int64, quantity int64, saletype int64) error {
	var name string
	var price int64
	q, qErr := DB.Query(`SELECT producto_nombre, producto_precio 
	FROM producto 
	WHERE producto_codigo = ?`, barcode)
	if qErr != nil {
		return fmt.Errorf("addToSale Query: %v", qErr)
	}
	defer q.Close()
	for q.Next() {
		if err := q.Scan(&name, &price); err != nil {
			return fmt.Errorf("addToSale Scan: %v", err)
		}
	}

	r, rErr := DB.Exec(`INSERT INTO transaccion_salida_producto
	(name, code, price, quantity, type) VALUES
	(?, ?, ?, ?, ?)`, name, barcode, price, quantity, saletype)
	if rErr != nil {
		return fmt.Errorf("addToSale Insert: %v", rErr)
	}
	fmt.Println(r)
	return nil
}

func GetSaleTransactionTable() ([]Producto_salida_join, error) {
	var table []Producto_salida_join
	rows, err := DB.Query(`select 
	name, code, price, quantity, s.salida_tipo_nombre
	from transaccion_salida_producto
	JOIN salida_tipo s
	on s.salida_tipo_id = type;`)
	if err != nil {
		return nil, fmt.Errorf("getSaleTransactionTable: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var prod Producto_salida_join
		if err := rows.Scan(&prod.Nombre, &prod.Codigo, &prod.PrecioVenta, &prod.Cantidad, &prod.TipoSalida); err != nil {
			return nil, fmt.Errorf("getSaleTransactionTable Scan: %v", err)
		}
		prod.Fecha = time.Now().Format("yyyy-mm-dd")
		fmt.Println(prod.Cantidad, prod.Codigo, prod.Nombre)
		table = append(table, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AllSales: %v", err)
	}

	return table, nil
}

func CompleteSale() error {
	r, err := DB.Exec()
}