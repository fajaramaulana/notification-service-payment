package repository

import (
	"github.com/fajaramaulana/notification-service-payment/internal/config/kafkaconfig"
	"github.com/fajaramaulana/notification-service-payment/internal/model"
)

type NotificationRepository interface {
	// GetNotificationByID gets a notification by its ID
	// GetNotificationByID(id int) (*Notification, error)

	// InsertNotificationDB inserts a new notification

	// InsertNotificationKafka insert a new notification to kafka topic

	GetDetailUserByUserID(userID int) (*model.DetailUser, error)
	SendMessageToKafka(producer kafkaconfig.KafkaProducer, message []byte) error
	InsertNotificationKafka(data model.NotificationToKafka) error
	InsertNotificationDB(notification model.InsertToDB) error
}
