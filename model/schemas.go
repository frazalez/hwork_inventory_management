package model

import "time"

type Usuario struct {
	Usuario_id         int64
	Usuario_nombre     string
	Usuario_psswd      string
	Usuario_privilegio int64
	Usuario_activo     bool
}

type UsrAuth struct {
	Usuario    string
	Passwd     string
	Activo     int
	Privilegio string
}

type Distribuidor struct {
	Dist_id     int64
	Dist_nombre string
}

type Salida struct {
	Salida_id      int64
	Salida_fecha   time.Time
	Salida_usuario int64
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

type Privilegios struct {
	Privilegio_id     int64
	Privilegio_nombre string
}
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
