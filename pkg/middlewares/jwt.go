package middlewares

import (
	"Beside-Mom-BE/configs"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(config configs.JWT) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenString := ctx.Get("Authorization")
		if tokenString == "" || len(tokenString) < 8 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":      "Unauthorized",
				"status_code": fiber.StatusUnauthorized,
				"message":     "Missing or invalid token",
				"result":      nil,
			})
		}

		tokenString = tokenString[7:]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(config.Secret), nil
		})

		if err != nil || !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":      "Unauthorized",
				"status_code": fiber.StatusUnauthorized,
				"message":     "Invalid or expired token",
				"result":      nil,
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":      "Unauthorized",
				"status_code": fiber.StatusUnauthorized,
				"message":     "Invalid token claims",
				"result":      nil,
			})
		}

		ctx.Locals("user_id", claims["user_id"])
		ctx.Locals("role", claims["role"])
		return ctx.Next()
	}
}

func AdminMiddleware(ctx *fiber.Ctx) error {
	userRole, ok := ctx.Locals("role").(string)
	if !ok || userRole != "Admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusForbidden,
			"message":     "Forbidden: Admin access required",
			"result":      nil,
		})
	}

	return ctx.Next()
}
