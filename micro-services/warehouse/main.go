package main

import (
	"log"
	cj "monorepo-ecommerce/micro-services/warehouse/cron"
	"monorepo-ecommerce/micro-services/warehouse/db"
	"monorepo-ecommerce/micro-services/warehouse/handler"
	"monorepo-ecommerce/micro-services/warehouse/repository"
	"monorepo-ecommerce/micro-services/warehouse/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
)

func main() {
	// Initialize Database
	dbConn := db.InitDatabase("./../../data/ecommerce.db")
	db.RunMigrations(dbConn, "./migrations/init.sql")
	defer dbConn.Close()

	// Initialize Echo
	e := echo.New()

	// Use Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize repository, service, and handler
	warehouseRepo := repository.NewWarehouseRepository(dbConn)
	stockRepo := repository.NewStockRepository(dbConn)
	productRepo := repository.NewProductRepository("http://localhost:7002")
	warehouseService := service.NewWarehouseService(warehouseRepo, stockRepo, productRepo)
	handler.RegisterWarehouseRoutes(e, warehouseService)

	// Init cronjob
	autoSyncStock := cj.NewAutoSyncStockJob(productRepo, stockRepo, warehouseRepo)
	c := cron.New()
	c.AddFunc("@every 2m", func() {
		autoSyncStock.Run()
	})
	c.Start()

	// Start server in goroutine
	go func() {
		log.Println("Starting server on port 7005...")
		if err := e.Start(":7005"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Waiting program
	select {}
}
