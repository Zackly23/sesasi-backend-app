package middlewares

import (
	"sesasi-backend-app/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RoleMiddleware(db *gorm.DB, allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ambil user dari context (AuthMiddleware yang set user_id)
		userRawId := c.Locals("user_id")
		if userRawId == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User ID tidak ditemukan di context",
			})
		}

		var user models.User
		if err := db.Preload("Role").First(&user, "id = ?", userRawId).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User tidak ditemukan",
			})
		}


		// cek apakah role user ada di daftar allowedRoles
		for _, role := range allowedRoles {
			if user.Role.Name == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: role tidak diizinkan",
		})
	}
}
