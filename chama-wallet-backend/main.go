package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/stellar/go/keypair"
	"gorm.io/gorm"

	"chama-wallet-backend/database"
	"chama-wallet-backend/routes"
	"chama-wallet-backend/services"
)

var DB *gorm.DB

func main() {
	database.ConnectDB()
	app := fiber.New()
	routes.Setup(app)
	routes.SetupSorobanRoutes(app)
	

	// routes.GroupRoutes(app)

	fmt.Println("✅ Setting up group routes...")
	routes.GroupRoutes(app)

	fmt.Println("🚀 Server starting on localhost:3000")
	log.Fatal(app.Listen("localhost:3000"))

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
