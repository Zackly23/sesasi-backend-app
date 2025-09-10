package routes

import (
	"fmt"
	"sesasi-backend-app/handlers"
	"sesasi-backend-app/middlewares"

	// "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	// "gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, h *handlers.Handler) {
	fmt.Println("Setting up routes...")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is running ✅")
	})

	app.Get("/health", handlers.CheckHealth)

	api := app.Group("/api")
	v1 := api.Group("/v1")
	
	// Public routes (tanpa JWT)
	auth := v1.Group("/auth")
	
	auth.Post("/signup", func(c *fiber.Ctx) error {
		return handlers.SignUp(c, h)
	})
	
	auth.Post("/login", func(c *fiber.Ctx) error {
		return handlers.Login(c, h)
	})

	auth.Post("/logout", middlewares.JWTMiddleware(h.DB) , func(c *fiber.Ctx) error {
		return handlers.Logout(c, h)
	})

	auth.Get("/refresh", func(c *fiber.Ctx) error {
		return handlers.Refresh(c, h)
	})

	adminRoutes := v1.Group(
		"/admin",
		middlewares.JWTMiddleware(h.DB),
		middlewares.RoleMiddleware(h.DB, []string{"admin"}),
	)

	adminRoutes.Get(
		"/users",
		func(c *fiber.Ctx) error {
			return handlers.GetAllUsers(c, h)
		},
	)

	adminRoutes.Post(
		"/users/verifikator",
		func(c *fiber.Ctx) error {
			return handlers.CreateVerifikator(c, h)
		},
	)

	adminRoutes.Patch(
		"/users/:userId/verifikator",
		func(c *fiber.Ctx) error {
			return handlers.ChangeUserRole(c, h)
		},
	)

	adminRoutes.Patch(
		"/users/:userId/reset-password",
		func(c *fiber.Ctx) error {
			return handlers.ResetPassword(c, h)
		},
	)

	adminRoutes.Get(
		"/izins",
		func(c *fiber.Ctx) error {
			return handlers.GetAllPengajuanIzin(c, h)
		},
	)

	verifikatorRoutes := v1.Group(
		"/verifikator",
		middlewares.JWTMiddleware(h.DB),
		middlewares.RoleMiddleware(h.DB, []string{"verifikator"}),
	)

	verifikatorRoutes.Get(
		"/users",
		func(c *fiber.Ctx) error {
			return handlers.GetUserVerified(c, h)
		},
	)

	verifikatorRoutes.Patch(
		"/users/:userId/verify" ,
		func(c *fiber.Ctx) error {
			return handlers.VerifyUser(c, h)
		},
	)

	verifikatorRoutes.Get(
		"/izins",
		func(c *fiber.Ctx) error {
			return handlers.GetFilteredPengajuanIzin(c, h)
		},
	)

	verifikatorRoutes.Patch(
		"/izins/:izinId/status",
		func(c *fiber.Ctx) error {
			return handlers.UpdateStatusPengajuanIzin(c, h)
		},
	)

	// Group routes for user
	userRoutes := v1.Group(
		"/user",
		middlewares.JWTMiddleware(h.DB),
		middlewares.RoleMiddleware(h.DB, []string{"user"}),
	)

	userRoutes.Post(
		"/izins",
		func(c *fiber.Ctx) error {
			return handlers.CreatePengajuanIzin(c, h)
		},
	)

	userRoutes.Get(
		"/izins",
		func(c *fiber.Ctx) error {
			return handlers.GetUserPengajuanIzin(c, h)
		},
	)

	userRoutes.Get(
		"/izins/:izinId",
		func(c *fiber.Ctx) error {
			return handlers.GetDetailPengajuanIzin(c, h)
		},
	)

	userRoutes.Put(
		"/izins/:izinId",
		func(c *fiber.Ctx) error {
			return handlers.UpdateDetailPengajuanIzin(c, h)
		},
	)

	userRoutes.Delete(
		"/izins/:izinId",
		func(c *fiber.Ctx) error {
			return handlers.DeletePengajuanIzin(c, h)
		},
	)

	userRoutes.Patch(
		"/izins/:izinId/cancel",
		func(c *fiber.Ctx) error {
			return handlers.CancelPengajuanIzin(c, h)
		},
	)

	userRoutes.Get(
		"/izins/:izinId/status",
		func(c *fiber.Ctx) error {
			return handlers.GetStatusPengajuanIzin(c, h)
		},
	)

	userRoutes.Patch(
		"/password",
		func(c *fiber.Ctx) error {
			return handlers.UpdatePassword(c, h)
		},
	)
}

