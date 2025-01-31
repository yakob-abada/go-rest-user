package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	_ "github.com/yakob-abada/go-rest-user/docs" // Import Swagger docs
)

// @title User Management API
// @version 1.0
// @description A Golang API for user management with PostgreSQL, RabbitMQ, and bcrypt.
// @host localhost:8080
// @BasePath /

func main() {
	e := echo.New()

	// Serve Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start API
	log.Println("Starting server on port 8080...")
	e.Logger.Fatal(e.Start(":8080"))
}
