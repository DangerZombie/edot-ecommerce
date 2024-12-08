package main

import (
	"log"
	cj "monorepo-ecommerce/micro-services/order/cron"
	"monorepo-ecommerce/micro-services/order/db"
	"monorepo-ecommerce/micro-services/order/handler"
	"monorepo-ecommerce/micro-services/order/repository"
	"monorepo-ecommerce/micro-services/order/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
)

func main() {
	// Init database
	dbConn := db.InitDatabase("./../../data/ecommerce.db")
	db.RunMigrations(dbConn, "./migrations/init.sql")
	defer dbConn.Close()

	// Init Product Repository
	productRepo := repository.NewProductRepository("http://localhost:7002") // URL Product Service

	// Init Shop Repository
	shopRepo := repository.NewShopRepository("http://localhost:7004")

	// Initiate Echo
	e := echo.New()

	// Use middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Init Order Repository, Service, Handler
	orderRepo := repository.NewOrderRepository(dbConn)
	orderService := service.NewOrderService(orderRepo, productRepo, shopRepo)
	handler.RegisterOrderRoutes(e, orderService)

	// Init cronjob
	autoCancelJob := cj.NewAutoCancelJob(orderRepo, productRepo)
	c := cron.New()
	c.AddFunc("@every 2m", func() {
		autoCancelJob.Run()
	})
	c.Start()

	// Start server in goroutine
	go func() {
		log.Println("Starting server on port 7003...")
		if err := e.Start(":7003"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Waiting program
	select {}
}
