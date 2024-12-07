package model

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Salida struct {
	Id      int64
	Fecha   string
	Usuario string
	Tipo string
}

type Salida_tipo struct {
	Salida_tipo_id     int64
	Salida_tipo_nombre string
}
type Producto_salida struct {
	Pro_sal_id       int64
	pro_sal_venta    int64
	pro_sal_producto int64
	pro_sal_cantidad int64
	pro_sal_precio   int64
	pro_sal_tipo     int64
}

type Producto_salida_join struct {
	Nombre        string
	Codigo        int64
	PrecioVenta   int64
	Cantidad      int64
	TipoSalida    string
	Fecha         string
	UsuarioSalida string
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
func AllSalesMain(db *sql.DB) ([]Salida, error) {
	var producto []Salida
	rows, err := db.Query(`SELECT s.salida_id, s.salida_fecha, st.salida_tipo_nombre, u.usuario_nombre from salida s
JOIN salida_tipo st on s.salida_tipo = st.salida_tipo_id
JOIN usuario u on s.salida_usuario = u.usuario_id;`)
	if err != nil {
		return nil, fmt.Errorf("AllSalesMain Query error in line 67: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var prod Salida
		if err := rows.Scan(&prod.Id, &prod.Fecha, &prod.Tipo, &prod.Usuario); err != nil {
			return nil, fmt.Errorf("AllSales: %v", err)
		}
		producto = append(producto, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AllSales: %v", err)
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

func GetSaleTransactionTable() ([]TransactionProduct, error) {
	var table []TransactionProduct
	rows, err := DB.Query(`select
	name, code, price, quantity, type
	from transaccion_salida_producto
	JOIN salida_tipo s
	on s.salida_tipo_id = type;`)
	if err != nil {
		return nil, fmt.Errorf("getSaleTransactionTable: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var prod TransactionProduct
		if err := rows.Scan(&prod.Name, &prod.Code, &prod.Price, &prod.Quantity, &prod.Type); err != nil {
			return nil, fmt.Errorf("getSaleTransactionTable Scan: %v", err)
		}
		table = append(table, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AllSales: %v", err)
	}

	return table, nil
}

func CompleteSale(username string) error {
	transacRows, trError := DB.Query(`SELECT tsp.name, tsp.code, tsp.price , tsp.quantity, tsp.type
FROM transaccion_salida_producto tsp;`)

	if trError != nil {
		return fmt.Errorf("CompleteSale SelectTransactionErr: %v", trError)
	}
	var transacTable []TransactionProduct

	defer transacRows.Close()
	for transacRows.Next() {
		var transacRowValues TransactionProduct
		if err := transacRows.Scan(&transacRowValues.Name, &transacRowValues.Code, &transacRowValues.Price, &transacRowValues.Quantity, &transacRowValues.Type); err != nil {
			return fmt.Errorf("CompleteSale transacRowsScan %v", err)
		}
		transacTable = append(transacTable, transacRowValues)
	}

	user, uError := getIDFromUsername(username)
	if uError != nil {
		return fmt.Errorf("CompleteSale GetIDFromUsername: %v", uError)
	}

	log.Print(transacTable[0].Type)

	insertSaleRows, isErr := DB.Exec(`INSERT INTO salida (salida_fecha, salida_tipo, salida_usuario)
		VALUES (?, ?, ?)`, time.Now().Format(time.DateTime), transacTable[0].Type, user)
	if isErr != nil {
		return fmt.Errorf("CompleteSale insertSalida %v", isErr)
	}
	log.Printf("inserted rows: %v", insertSaleRows)

	lastInsertRow := DB.QueryRow(`SELECT LAST_INSERT_ID() FROM salida`)
	var lastInsertID int

	if err := lastInsertRow.Scan(&lastInsertID); err != nil {
		return fmt.Errorf(`CompleteTable findLastInsertID`, err)
	}

	for i := range transacTable {
		PIDrow := DB.QueryRow(`SELECT producto_id FROM producto WHERE producto_codigo = ?`, transacTable[i].Code)
		var PID int
		if err := PIDrow.Scan(&PID); err != nil {
			return fmt.Errorf("CompleteSale InsertTransaction FindPID: %v", err)
		}
		rows, err := DB.Exec(`INSERT INTO producto_salida (pro_sal_salida, pro_sal_producto, pro_sal_cantidad, pro_sal_precio)
				VALUES (?,?,?,?)`, lastInsertID, PID, transacTable[i].Quantity, transacTable[i].Price)
		if err != nil {
			return fmt.Errorf("CompleteSale InsertTransaction InsertIntoProductoSalida: %v", err)
		}
		log.Printf("rows affected: %v", rows)
	}

	return nil
}

func CalculateTotalSale() (int, error) {
	row := DB.QueryRow(`select SUM(sq.p*sq.q)
	from (SELECT price as p, quantity as q from transaccion_salida_producto tsp) sq
	;`)

	var total int
	if err := row.Scan(&total); err != nil {
		return -1, fmt.Errorf("calculateTotalSale: %v", err)
	}
	return total, nil
}
