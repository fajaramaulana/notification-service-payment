package kafkaconfig_test

import (
	"errors"
	"testing"

	"github.com/IBM/sarama"
	"github.com/fajaramaulana/notification-service-payment/internal/config/kafkaconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSyncProducer is a mock implementation of sarama.SyncProducer.
type MockSyncProducer struct {
	mock.Mock
}

func (m *MockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	args := m.Called(msg)
	return args.Get(0).(int32), args.Get(1).(int64), args.Error(2)
}

func (m *MockSyncProducer) SendMessages(messages []*sarama.ProducerMessage) error {
	args := m.Called(messages)
	return args.Error(0)
}

func (m *MockSyncProducer) Close() error {
	return nil
}

func (m *MockSyncProducer) IsTransactional() bool {
	return false
}

func (m *MockSyncProducer) AbortTxn() error {
	return nil
}

func (m *MockSyncProducer) BeginTxn() error {
	return nil
}

func (m *MockSyncProducer) CommitTxn() error {
	return nil
}

func (m *MockSyncProducer) AddMessageToTxn(*sarama.ConsumerMessage, string, *string) error {
	return nil
}

func (m *MockSyncProducer) AddOffsetsToTxn(map[string][]*sarama.PartitionOffsetMetadata, string) error {
	return nil
}

// Mocking the TxnStatus method to return a default status.
func (m *MockSyncProducer) TxnStatus() sarama.ProducerTxnStatusFlag {
	return 0 // Using a basic integer, as the actual value is not needed for testing
}

// TestNewSaramaProducer tests the creation of a new SaramaProducer.
func TestNewSaramaProducer(t *testing.T) {
	mockProducer := new(MockSyncProducer)
	saramaProducer := kafkaconfig.NewSaramaProducer(mockProducer)

	assert.NotNil(t, saramaProducer)
}

// TestSendMessage_Success tests the successful sending of a message.
func TestSendMessage_Success(t *testing.T) {
	mockProducer := new(MockSyncProducer)
	mockProducer.On("SendMessage", mock.Anything).Return(int32(0), int64(0), nil)

	saramaProducer := kafkaconfig.NewSaramaProducer(mockProducer)
	topic := "test-topic"
	message := []byte("test-message")

	err := saramaProducer.SendMessage(topic, message)

	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

// TestSendMessage_Error tests the behavior when an error occurs during message sending.
func TestSendMessage_Error(t *testing.T) {
	mockProducer := new(MockSyncProducer)
	mockProducer.On("SendMessage", mock.Anything).Return(int32(0), int64(0), errors.New("send error"))

	saramaProducer := kafkaconfig.NewSaramaProducer(mockProducer)
	topic := "test-topic"
	message := []byte("test-message")

	err := saramaProducer.SendMessage(topic, message)

	assert.Error(t, err)
	assert.Equal(t, "send error", err.Error())
	mockProducer.AssertExpectations(t)
}
