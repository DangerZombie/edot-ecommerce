package main

import (
	"monorepo-ecommerce/micro-services/shop/db"
	"monorepo-ecommerce/micro-services/shop/handler"
	"monorepo-ecommerce/micro-services/shop/repository"
	"monorepo-ecommerce/micro-services/shop/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initate Database
	dbConn := db.InitDatabase("./../../data/ecommerce.db")
	db.RunMigrations(dbConn, "./migrations/init.sql")
	defer dbConn.Close()

	// Initiate Echo
	e := echo.New()

	// Use middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize repository, service, handler
	warehouseRepo := repository.NewWarehouseRepository("http://localhost:7005")
	userRepo := repository.NewShopRepository(dbConn)
	userService := service.NewShopService(userRepo, warehouseRepo)
	handler.RegisterShopRoutes(e, userService)

	// Start server
	e.Logger.Fatal(e.Start(":7004"))
}
