package database

import (
    "chama-wallet-backend/models"
    "log"
)

func RunMigrations() {
    log.Println("Running database migrations...")
    
    err := DB.AutoMigrate(
        &models.User{},
        &models.Group{},
        &models.Member{},
        &models.Contribution{},
        &models.GroupInvitation{},
        &models.AdminNomination{},
        &models.PayoutRequest{},
        &models.PayoutApproval{},
        &models.Notification{},
    )
    
    if err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }
    
    log.Println("✅ Database migrations completed successfully")
}