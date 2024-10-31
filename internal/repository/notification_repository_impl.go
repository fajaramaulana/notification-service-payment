package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/fajaramaulana/notification-service-payment/internal/config"
	"github.com/fajaramaulana/notification-service-payment/internal/config/kafkaconfig"
	"github.com/fajaramaulana/notification-service-payment/internal/model"
	"github.com/sirupsen/logrus"
)

type NotificationRepositoryImpl struct {
	Producer kafkaconfig.KafkaProducer
	config   config.Config
	db       *sql.DB
}

func NewNotificationRepository(producer kafkaconfig.KafkaProducer, config config.Config, db *sql.DB) NotificationRepository {
	return &NotificationRepositoryImpl{
		Producer: producer,
		config:   config,
		db:       db,
	}
}

func (n *NotificationRepositoryImpl) InsertNotificationKafka(data model.NotificationToKafka) error {
	logrus.Info("Inserting notification to kafka")

	// Prepare time and location data
	t := time.Now()
	timeLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		logrus.Error("Failed to load location")
		return err
	}
	logrus.Infof("Processing data at %s", t.In(timeLocation).String())

	mesData := model.InsertToKafka{
		Status:  true,
		Message: "Success",
		Time:    t.In(timeLocation).String(),
		Data:    data.DetailNotificationToKafka,
	}

	message, err := json.Marshal(mesData)
	if err != nil {
		logrus.Error("Failed to marshal data")
		return err
	}

	// Send message via Kafka

	err = n.SendMessageToKafka(n.Producer, message)
	if err != nil {
		logrus.Errorf("Failed to send message: %v", err)
	}

	logrus.Info("Data inserted successfully")
	return nil
}

func (n *NotificationRepositoryImpl) SendMessageToKafka(producer kafkaconfig.KafkaProducer, message []byte) error {
	topic := n.config.Get("KAFKA_TOPIC_MAIN")
	err := producer.SendMessage(topic, message)
	if err != nil {
		logrus.Errorf("Failed to send message: %v", err)
	} else {
		logrus.Infof("Message sent to topic %s successfully", topic)
	}
	return err
}

func (n *NotificationRepositoryImpl) GetDetailUserByUserID(userID int) (*model.DetailUser, error) {
	var user model.DetailUser
	query := "SELECT id, username, email, first_name, last_name, phone_number FROM users WHERE id = ?"
	row := n.db.QueryRow(query, userID)
	env := n.config.Get("ENV")
	if env != "production" {
		logrus.WithFields(logrus.Fields{
			"query": query,
			"id":    userID,
		}).Info("Querying database")
	}

	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.FName, &user.LName, &user.Phone); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err // Other error
	}
	return &user, nil
}
