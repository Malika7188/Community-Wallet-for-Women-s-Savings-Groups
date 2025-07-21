package services

import (
	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func CreateNotification(userID, groupID, notificationType, title, message string) error {
	notification := models.Notification{
		ID:        uuid.NewString(),
		UserID:    userID,
		GroupID:   groupID,
		Type:      notificationType,
		Title:     title,
		Message:   message,
		CreatedAt: time.Now(),
	}
	return database.DB.Create(&notification).Error
}

func GetUserNotifications(userID string) ([]models.Notification, error) {
	var notifications []models.Notification
	err := database.DB.
		Where("user_id = ?", userID).
		Preload("Group").
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

func MarkNotificationAsRead(notificationID string) error {
	return database.DB.Model(&models.Notification{}).
		Where("id = ?", notificationID).
		Update("read", true).Error
}

func SendContributionReminders() error {
	// Find active groups with upcoming contribution deadlines (5 days)
	var groups []models.Group
	fiveDaysFromNow := time.Now().AddDate(0, 0, 5)

	if err := database.DB.Where("status = ? AND next_contribution_date <= ?", "active", fiveDaysFromNow).Find(&groups).Error; err != nil {
		return err
	}

	for _, group := range groups {
		var members []models.Member
		database.DB.Where("group_id = ? AND status = ?", group.ID, "approved").Find(&members)

		for _, member := range members {
			daysUntil := int(group.NextContributionDate.Sub(time.Now()).Hours() / 24)
			CreateNotification(
				member.UserID,
				group.ID,
				"contribution_reminder",
				"Contribution Reminder",
				fmt.Sprintf("Your contribution of %.2f XLM is due in %d days", group.ContributionAmount, daysUntil),
			)
		}
	}
	return nil
}

func UpdateNextContributionDate(groupID string) error {
	var group models.Group
	if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
		return err
	}

	nextDate := time.Now().AddDate(0, 0, group.ContributionPeriod)
	return database.DB.Model(&group).Update("next_contribution_date", nextDate).Error
}

func NotifyPayoutApproved(groupID, recipientID string, amount float64) error {
	var members []models.Member
	database.DB.Where("group_id = ? AND status = ?", groupID, "approved").Find(&members)

	var recipient models.User
	database.DB.First(&recipient, "id = ?", recipientID)

	for _, member := range members {
		CreateNotification(
			member.UserID,
			groupID,
			"payout_approved",
			"Payout Approved",
			fmt.Sprintf("A payout of %.2f XLM to %s has been approved and will be processed", amount, recipient.Name),
		)
	}
	return nil
}
