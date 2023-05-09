package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORS() func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Accept-Language, Content-Length",
	})
}
