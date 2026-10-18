package middleware

import (
	"github.com/MaulanaBarzaqi/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func RoleMiddleware(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		role := claims["role"].(string)

		if role != requiredRole {
			return utils.Forbidden(c, "Forbidden")
		}
		return c.Next()
	}
}
func RoleMiddlewareWithException(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		publicId := c.Params("id")
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		role := claims["role"].(string)
		pubIDToken := claims["pub_id"].(string)

		if (role != requiredRole) && (publicId != pubIDToken) {
			return utils.Forbidden(c, "Forbidden")
		}
		return c.Next()
	}
}