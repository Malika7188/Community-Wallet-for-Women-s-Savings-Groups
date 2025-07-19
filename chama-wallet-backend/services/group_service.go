package services

import (
	"github.com/google/uuid"
	"time"
	"errors"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/utils"
)

var groups = make(map[string]models.Group)

func CreateGroup(name, description, creatorID string) (models.Group, error) {
	wallet, err := utils.GenerateStellarWallet()
	if err != nil {
		return models.Group{}, err
	}

	group := models.Group{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Wallet:      wallet.PublicKey,
		CreatorID:   creatorID,
		Status:      "pending",
	}

	if err := database.DB.Create(&group).Error; err != nil {
		return models.Group{}, err
	}

	// Automatically add creator as admin
	creator := models.Member{
		ID:       uuid.NewString(),
		GroupID:  group.ID,
		UserID:   creatorID,
		Role:     "creator",
		Status:   "approved",
		JoinedAt: time.Now(),
	}
	database.DB.Create(&creator)

	return group, nil
}

func GetGroupByID(groupID string) (models.Group, error) {
	var group models.Group
	err := database.DB.Preload("Members.User").Preload("Creator").First(&group, "id = ?", groupID).Error
	return group, err
}

func AddMemberToGroup(groupID, userID, walletAddress string) (models.Group, error) {
	var group models.Group
	if err := database.DB.Preload("Members").First(&group, "id = ?", groupID).Error; err != nil {
		return group, err
	}

	// Check if member already exists
	for _, member := range group.Members {
		if member.UserID == userID {
			return group, nil // Member already exists
		}
	}

	member := models.Member{
		ID:       uuid.NewString(),
		GroupID:  groupID,
		UserID:   userID,
		Wallet:   walletAddress,
		Role:     "member",
		Status:   "approved",
		JoinedAt: time.Now(),
	}
	if err := database.DB.Create(&member).Error; err != nil {
		return group, err
	}

	group.Members = append(group.Members, member)
	return group, nil
}

// func Contribute(groupID, memberID string, amount float64) error {
// 	// 1. Call Soroban
// 	resp, err := ContributeOnChain(fmt.Sprintf("%.0f", amount))
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("On-chain response:", resp)

// 	// 2. Save to DB
// 	contribution := models.Contribution{
// 		ID:       uuid.New().String(),
// 		GroupID:  groupID,
// 		MemberID: memberID,
// 		Amount:   amount,
// 	}
// 	return database.DB.Create(&contribution).Error
// }

func GetGroupWithMembers(groupID string) (models.Group, error) {
	var group models.Group
	err := database.DB.Preload("Members").First(&group, "id = ?", groupID).Error
	return group, err
}
func GetAllGroups() ([]models.Group, error) {
	var groups []models.Group
	if err := database.DB.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func GetUserGroups(userID string) ([]models.Group, error) {
	var groups []models.Group
	err := database.DB.
		Joins("JOIN members ON groups.id = members.group_id").
		Where("members.user_id = ? AND members.status = ?", userID, "approved").
		Preload("Members.User").
		Preload("Creator").
		Find(&groups).Error
	return groups, err
}

func InviteUserToGroup(groupID, inviterID, email string) error {
	// Check if inviter is admin/creator
	var member models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND role IN ?", 
		groupID, inviterID, []string{"creator", "admin"}).First(&member).Error; err != nil {
		return errors.New("only admins can invite users")
	}

	// Check if user exists
	var user models.User
	userExists := database.DB.Where("email = ?", email).First(&user).Error == nil

	invitation := models.GroupInvitation{
		ID:        uuid.NewString(),
		GroupID:   groupID,
		InviterID: inviterID,
		Email:     email,
		Status:    "pending",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if userExists {
		invitation.UserID = user.ID
	}

	return database.DB.Create(&invitation).Error
}

func GetNonGroupMembers(groupID string) ([]models.User, error) {
	var users []models.User
	err := database.DB.
		Where("id NOT IN (SELECT user_id FROM members WHERE group_id = ? AND status = ?)", 
			groupID, "approved").
		Find(&users).Error
	return users, err
}

func ApproveGroupActivation(groupID, adminID string, settings models.GroupSettings) error {
	// Verify admin permissions
	var member models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND role IN ?", 
		groupID, adminID, []string{"creator", "admin"}).First(&member).Error; err != nil {
		return errors.New("insufficient permissions")
	}

	// Update group with contribution settings
	updates := map[string]interface{}{
		"status": "active",
		"contribution_amount": settings.ContributionAmount,
		"contribution_period": settings.ContributionPeriod,
		"payout_order": settings.PayoutOrder,
	}

	return database.DB.Model(&models.Group{}).Where("id = ?", groupID).Updates(updates).Error
}

func JoinGroupRequest(groupID, userID, walletAddress string) error {
	// Check if user is already a member
	var existingMember models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ?", groupID, userID).First(&existingMember).Error; err == nil {
		return errors.New("user is already a member of this group")
	}

	// Create pending membership
	member := models.Member{
		ID:       uuid.NewString(),
		GroupID:  groupID,
		UserID:   userID,
		Wallet:   walletAddress,
		Role:     "member",
		Status:   "pending",
		JoinedAt: time.Now(),
	}

	return database.DB.Create(&member).Error
}
