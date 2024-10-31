package config

import (
	"errors"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

func ConnectProducer(brokersUrl []string) (sarama.SyncProducer, error) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Retry logic for connecting to Kafka
func RetryKafkaConnection(brokers []string, maxRetries int, retryInterval time.Duration) (sarama.SyncProducer, error) {
	var producer sarama.SyncProducer
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {

		logrus.Infof("Attempting to connect to Kafka. Attempt: %d", attempt)
		producer, err = ConnectProducer(brokers)
		if err == nil {
			return producer, nil // Successfully connected
		}

		// Log the error and wait before retrying
		logrus.Errorf("Failed to connect to Kafka: %v. Retrying in %s...", err, retryInterval)
		time.Sleep(retryInterval)
		// Exponential backoff: Increase retry interval for each attempt
		retryInterval *= 2
	}

	return nil, err // Return error after max retries
}

func RetryKafkaConnectionMock(brokers []string, maxRetries int, retryInterval time.Duration, dialFunc func([]string) (sarama.SyncProducer, error)) (sarama.SyncProducer, error) {
	var producer sarama.SyncProducer
	var err error

	// Try to connect with retry logic
	for i := 0; i < maxRetries; i++ {
		producer, err = dialFunc(brokers)
		if err == nil {
			return producer, nil
		}
		time.Sleep(retryInterval)
	}

	return nil, errors.New("failed to connect to Kafka after multiple retries")
}
