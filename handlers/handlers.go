package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


type Handler struct {
	DB        *gorm.DB
	Validator *validator.Validate
}

func NewHandler(db *gorm.DB, val *validator.Validate) *Handler {
	return &Handler{DB: db, Validator: val}
}

func CheckHealth(ctx *fiber.Ctx) error {
	fmt.Println("Health check hit!") // Debug log
		if (ctx.Method() != "GET") {
			return fiber.NewError(504, "Method not allowed")
		}
	
		return ctx.JSON(fiber.Map{
			"status": "OK",
	})
}

