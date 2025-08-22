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
        &models.PayoutSchedule{},
        &models.Notification{},
        &models.RoundContribution{},
        &models.RoundStatus{},
    )
    
    if err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }
    
    // Add indexes for better performance
    DB.Exec("CREATE INDEX IF NOT EXISTS idx_round_status_group_round ON round_statuses(group_id, round)")
    DB.Exec("CREATE INDEX IF NOT EXISTS idx_payout_requests_auto_generated ON payout_requests(auto_generated, status)")
    DB.Exec("CREATE INDEX IF NOT EXISTS idx_round_contributions_group_round ON round_contributions(group_id, round, status)")
    
    log.Println("✅ Database migrations completed successfully")
}
