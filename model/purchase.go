package model

import (
	"database/sql"
	"fmt"
	"time"
)

type Entrada struct {
	Entrada_id           int64
	Entrada_fecha        time.Time
	Entrada_usuario      int64
	Entrada_distribuidor int64
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
