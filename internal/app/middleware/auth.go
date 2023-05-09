package middleware

import (
	"github.com/gofiber/fiber/v2"
	userservice "newanysock/internal/app/handler/user"
	"newanysock/internal/entity/user"
	"newanysock/pkg"
)

func JwtCheck(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	var userClaim user.UserClaim

	err := userservice.ParseJWTToken(cookie, &userClaim)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	c.Locals("userClaim", userClaim)
	return c.Next()
}

func IsAdmin(c *fiber.Ctx) error {
	data := c.Locals("userClaim")
	if data == nil {
		return c.JSON(fiber.Map{"error": "we cant detect your infomation"})
	}
	if pkg.CompareType(data, user.UserClaim{}) == false {
		return c.JSON(fiber.Map{"error": "we cant detect your infomation"})
	}
	newdata := data.(user.UserClaim)

	if newdata.Role != "admin" {
		c.Status(fiber.StatusForbidden)
		return c.JSON(fiber.Map{
			"message": "your action is disallowed",
		})
	} else if newdata.Role == "admin" {
		return c.Next()
	}
	c.Status(fiber.StatusForbidden)
	return c.JSON(fiber.Map{
		"message": "your action is disallowed",
	})
}
