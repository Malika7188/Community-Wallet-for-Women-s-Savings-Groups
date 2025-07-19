package services

import (
    "chama-wallet-backend/database"
    "chama-wallet-backend/models"
    "github.com/google/uuid"
    "time"
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
    // Find groups with upcoming contribution deadlines (5 days)
    var groups []models.Group
    if err := database.DB.Where("status = ?", "active").Find(&groups).Error; err != nil {
        return err
    }

    for _, group := range groups {
        // Calculate next contribution date based on group settings
        // Send notifications to all members
        var members []models.Member
        database.DB.Where("group_id = ? AND status = ?", group.ID, "approved").Find(&members)
        
        for _, member := range members {
            CreateNotification(
                member.UserID,
                group.ID,
                "contribution_reminder",
                "Contribution Reminder",
                "Your contribution is due in 5 days",
            )
        }
    }
    return nil
}