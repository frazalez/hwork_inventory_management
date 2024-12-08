package model

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Entrada struct {
	Id      int64
	Fecha   string
	Usuario string
}

type Producto_entrada struct {
	pro_ent_id       int64
	pro_ent_ent_fk   int64
	pro_ent_pro_fk   int64
	pro_ent_cantidad int64
	pro_ent_precio   int64
}

type Producto_entrada_join struct {
	Nombre   string
	Codigo   int64
	Precio   int64
	Cantidad int64
	Fecha    string
	Usuario  string
}
type TransactionPurchase struct {
	Name     string
	Code     int64
	Price    int64
	Quantity int64
}

func StartPurchase() error {
	r, err := DB.Exec(
		`DROP TABLE IF EXISTS transaccion_entrada_producto;`)
	fmt.Printf("Rows affected: %v", r)

	r2, err2 := DB.Exec(`CREATE TABLE transaccion_entrada_producto (
			id int(10) primary key auto_increment,
			name varchar(40),
			code int(20),
			price int(10),
			quantity int(10));`)
	fmt.Printf("Rows Affected: %v", r2)

	if err != nil {
		return fmt.Errorf("StartPurchase: %v", err)
	}
	if err2 != nil {
		return fmt.Errorf("StartPurchase: %v", err2)
	}
	return nil
}

func AllPurchasesMain(db *sql.DB) ([]Entrada, error) {
	var producto []Entrada
	rows, err := db.Query(`select
    e.entrada_id, e.entrada_fecha, u.usuario_nombre
    from
    entrada e
	JOIN
    usuario u on e.entrada_usuario = u.usuario_id;`)
	if err != nil {
		return nil, fmt.Errorf("AllPurchases Query error in purchase.go line 70: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var prod Entrada
		if err := rows.Scan(&prod.Id, &prod.Fecha, &prod.Usuario); err != nil {
			return nil, fmt.Errorf("AllPurchases scan error in purchase.go line 77: %v", err)
		}
		producto = append(producto, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AllPurchases: %v", err)
	}

	return producto, nil
}

func AllPurchases(db *sql.DB) ([]Producto_entrada_join, error) {
	var producto []Producto_entrada_join
	rows, err := db.Query(`select
    p.producto_nombre, p.producto_codigo,
    pe.pro_ent_precio, pe.pro_ent_cantidad ,
    e.entrada_fecha , u.usuario_nombre
	from
    producto p
	join
    producto_entrada pe on p.producto_id = pe.pro_ent_pro_fk
	JOIN
    entrada e on pe.pro_ent_ent_fk = e.entrada_id
	JOIN
    usuario u on e.entrada_usuario = u.usuario_id;`)
	if err != nil {
		return nil, fmt.Errorf("AllPurchases: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var prod Producto_entrada_join
		if err := rows.Scan(&prod.Nombre, &prod.Codigo, &prod.Precio, &prod.Cantidad, &prod.Fecha, &prod.Usuario); err != nil {
			return nil, fmt.Errorf("AllPurchases: %v", err)
		}
		producto = append(producto, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AllPurchases: %v", err)
	}

	return producto, nil
}

func AddToPurchase(barcode int64, quantity int64) error {
	var name string
	var price int64
	q, qErr := DB.Query(`SELECT producto_nombre, producto_precio
	FROM producto
	WHERE producto_codigo = ?`, barcode)
	if qErr != nil {
		return fmt.Errorf("addToPurchase Query: %v", qErr)
	}
	defer q.Close()
	for q.Next() {
		if err := q.Scan(&name, &price); err != nil {
			return fmt.Errorf("addToPurchase Scan: %v", err)
		}
	}

	r, rErr := DB.Exec(`INSERT INTO transaccion_entrada_producto
	(name, code, price, quantity) VALUES
	(?, ?, ?, ?)`, name, barcode, price, quantity)
	if rErr != nil {
		return fmt.Errorf("addToPurchase Insert: %v", rErr)
	}
	fmt.Println(r)
	return nil
}

func GetPurchaseTransactionTable() ([]TransactionPurchase, error) {
	var table []TransactionPurchase
	rows, err := DB.Query(`select
	name, code, price, quantity
	from transaccion_entrada_producto`)
	if err != nil {
		return nil, fmt.Errorf("getPurchaseTransactionTable: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var prod TransactionPurchase
		if err := rows.Scan(&prod.Name, &prod.Code, &prod.Price, &prod.Quantity); err != nil {
			return nil, fmt.Errorf("getPurchaseTransactionTable Scan: %v", err)
		}
		table = append(table, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AllPurchases: %v", err)
	}

	return table, nil
}

func CompletePurchase(username string) error {
	transacRows, trError := DB.Query(`SELECT tep.name, tep.code, tep.price , tep.quantity
FROM transaccion_entrada_producto tep;`)

	if trError != nil {
		return fmt.Errorf("CompletePurchase SelectTransactionErr: %v", trError)
	}
	var transacTable []TransactionPurchase

	defer transacRows.Close()
	for transacRows.Next() {
		var transacRowValues TransactionPurchase
		if err := transacRows.Scan(&transacRowValues.Name, &transacRowValues.Code, &transacRowValues.Price, &transacRowValues.Quantity); err != nil {
			return fmt.Errorf("CompletePurchase transacRowsScan %v", err)
		}
		transacTable = append(transacTable, transacRowValues)
	}

	user, uError := getIDFromUsername(username)
	if uError != nil {
		return fmt.Errorf("CompletePurchase GetIDFromUsername: %v", uError)
	}

	insertPurchaseRows, isErr := DB.Exec(`INSERT INTO entrada (entrada_fecha, entrada_usuario)
		VALUES (?, ?)`, time.Now().Format(time.DateTime), user)
	if isErr != nil {
		return fmt.Errorf("CompletePurchase insertentrada %v", isErr)
	}
	log.Printf("inserted rows: %v", insertPurchaseRows)

	lastInsertRow := DB.QueryRow(`SELECT LAST_INSERT_ID() FROM entrada`)
	var lastInsertID int

	if err := lastInsertRow.Scan(&lastInsertID); err != nil {
		return fmt.Errorf(`CompleteTable findLastInsertID: %v`, err)
	}

	for i := range transacTable {
		PIDrow := DB.QueryRow(`SELECT producto_id FROM producto WHERE producto_codigo = ?`, transacTable[i].Code)
		var PID int
		if err := PIDrow.Scan(&PID); err != nil {
			return fmt.Errorf("CompletePurchase InsertTransaction FindPID: %v", err)
		}
		rows, err := DB.Exec(`INSERT INTO producto_entrada (pro_ent_ent_fk, pro_ent_pro_fk, pro_ent_cantidad, pro_ent_precio)
				VALUES (?,?,?,?)`, lastInsertID, PID, transacTable[i].Quantity, transacTable[i].Price)
		if err != nil {
			return fmt.Errorf("CompletePurchase InsertTransaction InsertIntoProductoentrada: %v", err)
		}
		log.Printf("rows affected: %v", rows)
	}

	return nil
}

func CalculateTotalPurchase() (int, error) {
	row := DB.QueryRow(`select SUM(sq.p*sq.q)
	from (SELECT price as p, quantity as q from transaccion_entrada_producto tsp) sq
	;`)

	var total int
	if err := row.Scan(&total); err != nil {
		return -1, fmt.Errorf("calculateTotalPurchase: %v", err)
	}
	return total, nil
}
