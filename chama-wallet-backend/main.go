package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/stellar/go/keypair"
	"gorm.io/gorm"

	"chama-wallet-backend/database"
	"chama-wallet-backend/routes"
	"chama-wallet-backend/services"
)

var DB *gorm.DB

func main() {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	database.ConnectDB()
	database.RunMigrations()
	app := fiber.New()

	// Add CORS middleware
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))
	routes.Setup(app)
	routes.SetupSorobanRoutes(app)
	routes.GroupRoutes(app)
	routes.AuthRoutes(app)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server starting on port %s\n", port)
	log.Fatal(app.Listen("0.0.0.0:" + port))

	kp1, err := keypair.Random()
	if err != nil {
		log.Fatalf("Failed to generate keypair1: %v", err)
	}
	kp2, err := keypair.Random()
	if err != nil {
		log.Fatalf("Failed to generate keypair2: %v", err)
	}

	addr1 := kp1.Address()
	addr2 := kp2.Address()

	fmt.Println("Account 1:", kp1.Address())
	fmt.Println("Seed 1:", kp1.Seed())
	fmt.Println("Account 2:", kp2.Address())
	fmt.Println("Seed 2:", kp2.Seed())

	fmt.Println("From:", addr1)
	fmt.Println("To:", addr2)

	// 🚀 Fund both accounts using Friendbot
	if err := services.FundTestAccount(addr1); err != nil {
		log.Fatalf("Funding account 1 failed: %v", err)
	}
	if err := services.FundTestAccount(addr2); err != nil {
		log.Fatalf("Funding account 2 failed: %v", err)
	}

	// 🪙 Send XLM from kp1 to kp2
	if _, err := services.SendXLM(kp1.Seed(), addr2, "10"); err != nil {
		log.Fatalf("Transaction failed: %v", err)
	} else {
		fmt.Println("✅ Transaction sent successfully.")
	}
}
