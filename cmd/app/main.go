package main

import (
	"github.com/gofiber/fiber/v2"
	"newanysock/internal/app/handler"
	"newanysock/internal/app/middleware"
	"newanysock/pkg/database/mysql"
)

func main() {

	mysql.Connect()

	app := fiber.New()
	app.Use(middleware.CORS())
	//app.Use(recover.New())

	handler.Setup(app)

	app.Listen(":8000")
}
