package kafkaconfig

type KafkaProducer interface {
	SendMessage(topic string, message []byte) error
}
