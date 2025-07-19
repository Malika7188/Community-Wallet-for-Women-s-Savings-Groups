package services

import (
	"github.com/google/uuid"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/utils"
)

var groups = make(map[string]models.Group)

func CreateGroup(name string, description string) (models.Group, error) {
	wallet, err := utils.GenerateStellarWallet()
	if err != nil {
		return models.Group{}, err
	}

	group := models.Group{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Wallet:      wallet.PublicKey,
	}

	if err := database.DB.Create(&group).Error; err != nil {
		return models.Group{}, err
	}
	return group, nil
}

func GetGroupByID(id string) (models.Group, error) {
	var group models.Group
	err := database.DB.Preload("Members").First(&group, "id = ?", id).Error
	if err != nil {
		return models.Group{}, err
	}
	return group, nil
}

func AddMemberToGroup(groupID, walletAddress string) (models.Group, error) {
	var group models.Group
	if err := database.DB.Preload("Members").First(&group, "id = ?", groupID).Error; err != nil {
		return group, err
	}

	// Check if member already exists
	for _, member := range group.Members {
		if member.Wallet == walletAddress {
			return group, nil // Member already exists
		}
	}

	member := models.Member{
		ID:      uuid.NewString(),
		Wallet:  walletAddress,
		GroupID: groupID,
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
