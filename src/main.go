package main

import (
	"inventoryManagement/handlers"
	"inventoryManagement/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	echoInstance := echo.New()
	echoInstance.Use(middleware.Static("/static/*"))
	echoInstance.GET("/static/*", echo.WrapHandler(http.FileServer(http.FS(echoInstance.Filesystem))))
	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Secure())

	handlers.DefineRouting(echoInstance)
	model.DB = model.ConnectToDatabase()

	echoInstance.Logger.Fatal(echoInstance.Start(":8080"))

}
