package middlewares

import (
	"fmt"
	"os"
	"sesasi-backend-app/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func JWTMiddleware(db *gorm.DB) fiber.Handler {
	fmt.Println("JWT Middleware initialized")
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing or invalid authorization header",
				"error":   "Expected 'Authorization: Bearer <token>'",
			})
		}


		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET_KEY")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token is not valid",
				"error": err.Error(),
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token is not valid",
			})
		}

		// Cek apakah token sudah direvoke
		var tokenRecord models.PrivateAccessToken
		if err := db.Where("access_token = ? AND revoked = false", tokenStr).First(&tokenRecord).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token is no longer valid or has been revoked",
				"error": err.Error(),
			})
		}

		// Simpan data user ke context
		c.Locals("user_id", claims["user_id"])
		c.Locals("email", claims["email"])

		return c.Next()
	}
}
