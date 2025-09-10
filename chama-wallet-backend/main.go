package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"chama-wallet-backend/config"
	"chama-wallet-backend/database"
	"chama-wallet-backend/routes"
)

var DB *gorm.DB

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Initialize Stellar configuration
	config.InitStellarConfig()

	// Validate mainnet configuration if needed
	if err := config.ValidateMainnetConfig(); err != nil {
		log.Fatalf("❌ Configuration validation failed: %v", err)
	}

	// Connect to database and run migrations
	database.ConnectDB()
	database.RunMigrations()

	// Create Fiber app
	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Setup routes
	routes.Setup(app)
	routes.SetupSorobanRoutes(app)
	routes.GroupRoutes(app)
	routes.AuthRoutes(app)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Print startup information
	fmt.Printf("🚀 Chama Wallet API starting on port %s\n", port)
	fmt.Printf("🌐 Network: %s\n", config.Config.Network)
	fmt.Printf("🔗 Horizon: %s\n", config.Config.HorizonURL)
	fmt.Printf("📡 Soroban RPC: %s\n", config.Config.SorobanRPCURL)
	if config.Config.ContractID != "" {
		fmt.Printf("📋 Contract ID: %s\n", config.Config.ContractID)
	}
	
	if config.Config.IsMainnet {
		fmt.Println("⚠️  MAINNET MODE: Real funds will be used!")
	} else {
		fmt.Println("🧪 TESTNET MODE: Using test funds")
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
