package main

import (
	// "fmt"
	"fmt"
	"log"
	"os"
	"sesasi-backend-app/config"
	"sesasi-backend-app/handlers"
	"sesasi-backend-app/routes"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/joho/godotenv"
)



func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	//Inisialisasi Validator
	val := validator.New()

	// Connect DB
	db := config.SetupDatabase()

	// Inisialisasi Fiber
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, // 100 MB
	})

	// Middlewares
	app.Use(cors.New()) // default allow all


	// Buat handler instance (Depedency Injection)
	h := handlers.NewHandler(db, val)

	// Routes
	routes.SetupRoutes(app, h)

	// Print semua routes
	for _, route := range app.GetRoutes() {
		fmt.Printf("Method: %-6s Path: %s\n", route.Method, route.Path)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	log.Printf("🚀 Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))

}


