package handlers

import (
	"github.com/labstack/echo/v4"
)

func DefineRouting(e *echo.Echo) {

	e.GET("/", IndexEntranceHandler)
	e.POST("/", IndexEntranceHandler)
	e.GET("/login", LoginHandler)
	e.PUT("/views/tables", TablesPutHandler)
	e.PUT("/views/ventas", VentasPutHandler)
	e.PUT("/views/compras", ComprasPutHandler)
	e.PUT("/views/admin", AdminPutHandler)
	e.PUT("/views/crear_usuario", UsuariosPutHandler)
	e.POST("/views/salir", SalirPostHandler)
	e.POST("/create-session", CreateSessionHandler)
	e.POST("/create-user", CreateUserHandler)
	e.POST("/create-product", CreateProductHandler)
	e.PUT("/disableProduct", DisableProductHandler)
	e.PUT("/enableProduct", EnableProductHandler)
	e.PUT("/modifyProductForm", ModifyProductFormHandler)
	e.POST("/modifyProduct", ModifyProductHandler)
	e.POST("/small-table-search", SmallTableSearchHandler)
	e.GET("/get-sales-table", TablesSalesHandler)
	e.POST("/start-sale", StartSaleHandler)
	e.POST("/create-sale", CreateSaleHandler)
}
