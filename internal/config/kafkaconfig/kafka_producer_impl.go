package kafkaconfig

import (
	"github.com/IBM/sarama"
)

// SaramaProducer is a wrapper around sarama.SyncProducer that implements the KafkaProducer interface.
type SaramaProducer struct {
	producer sarama.SyncProducer
}

// NewSaramaProducer creates a new SaramaProducer.
func NewSaramaProducer(producer sarama.SyncProducer) *SaramaProducer {
	return &SaramaProducer{producer: producer}
}

// SendMessage sends a message to a specific Kafka topic.
func (p *SaramaProducer) SendMessage(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	_, _, err := p.producer.SendMessage(msg)
	return err
}
