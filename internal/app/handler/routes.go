package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	adminHandler "newanysock/internal/app/handler/admin"
	userHandler "newanysock/internal/app/handler/user"
	"newanysock/internal/app/middleware"
)

func EmptyHandler(c *fiber.Ctx) error {
	return c.Next()
}

func testerror(c *fiber.Ctx) error {
	var a []string
	fmt.Println(a[1])
	return nil
}

func Setup(app *fiber.App) {
	apiGroup := app.Group("/Api", EmptyHandler)

	userGroup := apiGroup.Group("/User",
		middleware.JwtCheck,
		EmptyHandler)

	adminGroup := apiGroup.Group("/Admin",
		middleware.JwtCheck,
		middleware.IsAdmin,
		EmptyHandler)

	app.Get("/error", testerror)

	// FOR ADMIN ONLY
	adminGroup.Post("/AddUser", adminHandler.AddUser)
	adminGroup.Post("/DeleteUser", adminHandler.DeleteUser)
	adminGroup.Post("/UpdateUser", adminHandler.UpdateUser)
	adminGroup.Get("/ShowUsers", adminHandler.ShowUsers)
	adminGroup.Post("/AddLicense", adminHandler.AddLicense)
	adminGroup.Post("/RemoveLicense", adminHandler.RemoveLisence)
	adminGroup.Post("/UpdateLicense", adminHandler.UpdateLisence)

	// FOR USER
	userGroup.Get("/Profile", userHandler.Profile)
	userGroup.Post("/ChangePassword", userHandler.ChangePassword)
	userGroup.Post("/UpdateUser", userHandler.UpdateUser)
	userGroup.Post("/Logout", userHandler.Logout)
	app.Post("/Login", userHandler.Login)
	app.Post("/register", userHandler.Register)

}
