package main

import (
	"monorepo-ecommerce/micro-services/product/db"
	"monorepo-ecommerce/micro-services/product/handler"
	"monorepo-ecommerce/micro-services/product/repository"
	"monorepo-ecommerce/micro-services/product/service"

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
	productRepo := repository.NewProductRepository(dbConn)
	productService := service.NewProductService(productRepo)
	handler.RegisterProductRoutes(e, productService)

	// Start server
	e.Logger.Fatal(e.Start(":7002"))
}
