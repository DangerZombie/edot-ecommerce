package main

import (
	"monorepo-ecommerce/micro-services/user/db"
	"monorepo-ecommerce/micro-services/user/handler"
	"monorepo-ecommerce/micro-services/user/repository"
	"monorepo-ecommerce/micro-services/user/service"

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
	userRepo := repository.NewUserRepository(dbConn)
	userService := service.NewUserService(userRepo)
	handler.RegisterUserRoutes(e, userService)

	// Start server
	e.Logger.Fatal(e.Start(":7001"))
}
