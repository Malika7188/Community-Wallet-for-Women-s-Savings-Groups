package services

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stellar/go/keypair"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
)

// TestAutomatedPayoutService tests the automated payout functionality
func TestAutomatedPayoutService(t *testing.T) {
	// Setup test database connection
	database.ConnectDB()
	
	aps := &AutomatedPayoutService{}

	// Test 1: Create test group and members
	testGroup := createTestGroup(t)
	testMembers := createTestMembers(t, testGroup.ID, 3)

	// Test 2: Simulate all members contributing
	simulateRoundContributions(t, testGroup, testMembers, 1)

	// Test 3: Check automatic payout creation
	err := aps.CheckAndCreateAutomaticPayouts()
	if err != nil {
		t.Fatalf("Failed to create automatic payouts: %v", err)
	}

	// Test 4: Verify payout was created
	var payout models.PayoutRequest
	if err := database.DB.Where("group_id = ? AND round = ? AND auto_generated = ?", 
		testGroup.ID, 1, true).First(&payout).Error; err != nil {
		t.Fatalf("Automatic payout was not created: %v", err)
	}

	// Test 5: Verify payout amount is correct
	expectedAmount := testGroup.ContributionAmount * float64(len(testMembers))
	if payout.Amount != expectedAmount {
		t.Errorf("Expected payout amount %.2f, got %.2f", expectedAmount, payout.Amount)
	}

	// Test 6: Test designated recipient
	recipient, err := aps.getDesignatedRecipient(testGroup, 1)
	if err != nil {
		t.Fatalf("Failed to get designated recipient: %v", err)
	}

	if recipient.UserID != payout.RecipientID {
		t.Errorf("Payout recipient mismatch")
	}

	// Cleanup
	cleanupTestData(t, testGroup.ID)
}

func createTestGroup(t *testing.T) models.Group {
	// Create test users
	creator := models.User{
		ID:     uuid.NewString(),
		Name:   "Test Creator",
		Email:  "creator@test.com",
		Wallet: "GCREATORWALLET123456789012345678901234567890123456",
	}
	database.DB.Create(&creator)

	// Create payout order
	payoutOrder := []string{creator.ID}
	payoutOrderJSON, _ := json.Marshal(payoutOrder)

	group := models.Group{
		ID:                 uuid.NewString(),
		Name:               "Test Group",
		Description:        "Test group for automated payouts",
		Wallet:             "GGROUPWALLET1234567890123456789012345678901234567890",
		CreatorID:          creator.ID,
		Status:             "active",
		ContributionAmount: 100.0,
		ContributionPeriod: 30,
		PayoutOrder:        string(payoutOrderJSON),
		CurrentRound:       1,
		MaxMembers:         10,
		MinMembers:         3,
		IsApproved:         true,
	}

	if err := database.DB.Create(&group).Error; err != nil {
		t.Fatalf("Failed to create test group: %v", err)
	}

	return group
}

func createTestMembers(t *testing.T, groupID string, count int) []models.Member {
	var members []models.Member

	for i := 0; i < count; i++ {
		user := models.User{
			ID:     uuid.NewString(),
			Name:   fmt.Sprintf("Test User %d", i+1),
			Email:  fmt.Sprintf("user%d@test.com", i+1),
			Wallet: fmt.Sprintf("GUSER%d234567890123456789012345678901234567890123456", i+1),
		}
		database.DB.Create(&user)

		member := models.Member{
			ID:       uuid.NewString(),
			GroupID:  groupID,
			UserID:   user.ID,
			Wallet:   user.Wallet,
			Role:     "member",
			Status:   "approved",
			JoinedAt: time.Now(),
		}
		database.DB.Create(&member)
		members = append(members, member)
	}

	return members
}

func simulateRoundContributions(t *testing.T, group models.Group, members []models.Member, round int) {
	for _, member := range members {
		contribution := models.RoundContribution{
			ID:        uuid.NewString(),
			GroupID:   group.ID,
			MemberID:  member.ID,
			Round:     round,
			Amount:    group.ContributionAmount,
			Status:    "confirmed",
			TxHash:    fmt.Sprintf("test-tx-hash-%s", uuid.NewString()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		database.DB.Create(&contribution)
	}
}

func cleanupTestData(t *testing.T, groupID string) {
	database.DB.Where("group_id = ?", groupID).Delete(&models.RoundContribution{})
	database.DB.Where("group_id = ?", groupID).Delete(&models.RoundStatus{})
	database.DB.Where("group_id = ?", groupID).Delete(&models.PayoutRequest{})
	database.DB.Where("group_id = ?", groupID).Delete(&models.Member{})
	database.DB.Where("id = ?", groupID).Delete(&models.Group{})
}