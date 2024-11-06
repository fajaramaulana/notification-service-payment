package main

import (
	"database/sql"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/fajaramaulana/notification-service-payment/internal/config"
	"github.com/fajaramaulana/notification-service-payment/internal/config/kafkaconfig"
	"github.com/fajaramaulana/notification-service-payment/internal/repository"
	"github.com/fajaramaulana/notification-service-payment/internal/service"
	pb "github.com/fajaramaulana/shared-proto-payment/proto/notification"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var ListenerFunc = func() (net.Listener, error) {
	return net.Listen("tcp", ":50052")
}

func init() {
	// Create the logs directory if it doesn't exist
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	// Set up Lumberjack logger for log rotation with JSON format
	fileLogger := &lumberjack.Logger{
		Filename:   "logs/notification-service-payment.log",
		MaxSize:    10,   // Maximum size in megabytes before rotating
		MaxBackups: 3,    // Maximum number of old log files to keep
		MaxAge:     28,   // Maximum number of days to retain a log file
		Compress:   true, // Compress old log files
	}

	// Create a MultiWriter to write to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, fileLogger)

	// Configure the global logrus logger to write to both outputs
	logrus.SetOutput(multiWriter)
	logrus.SetLevel(logrus.InfoLevel)

	// Set JSON format for log files and plain text for terminal
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func SetupServer(configuration config.Config, db *sql.DB, ListenerFunc func() (net.Listener, error), kafkaconfig kafkaconfig.KafkaProducer) (*grpc.Server, net.Listener, error) {
	repo := repository.NewNotificationRepository(kafkaconfig, configuration, db)
	service := service.NewNotificationService(repo, configuration)

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, service)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Set up listener
	listener, err := ListenerFunc()
	if err != nil {
		return nil, nil, err
	}

	return grpcServer, listener, nil

}

func main() {
	configuration := config.LoadConfiguration()

	producer, err := setupKafka(configuration)

	if err != nil {
		logrus.Fatalf("Failed to connect to Kafka: %v", err)
	}
	logrus.Infof("Configuration loaded")
	if err := run(configuration, config.ConnectDBMysql, ListenerFunc, producer); err != nil {
		logrus.Fatalf("Application failed: %v", err)
	}
}

func run(config config.Config, connectDB func(config.Config) (*sql.DB, error), createListener func() (net.Listener, error), kafkaconfig kafkaconfig.KafkaProducer) error {
	// Attempt to connect to the database
	db, err := connectDB(config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Ensure db is closed only if it’s initialized
	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	// Attempt to create the listener
	// Set up the gRPC server and listener
	grpcServer, listener, err := SetupServer(config, db, createListener, kafkaconfig)
	if err != nil {
		return fmt.Errorf("failed to set up server: %w", err)
	}
	defer listener.Close()

	// Start serving the gRPC server
	logrus.Info("Starting gRPC server...")
	logrus.Infof("gRPC server started successfully on port: %s", listener.Addr().String())
	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	// Continue with the rest of the function logic...
	return nil
}

func setupKafka(configuration config.Config) (kafkaconfig.KafkaProducer, error) {
	brokersUrl := []string{configuration.Get("KAFKA_URL")}
	maxRetries := configuration.Get("MAX_RETRIES")
	maxRetriesInt, err := strconv.Atoi(maxRetries)

	if err != nil {
		logrus.Errorf("Failed to convert max retries to int: %v", err)
		return nil, err
	}

	retryIntervalStr := configuration.Get("MAX_RETRIES_SECOND")
	retryInterval, err := strconv.Atoi(retryIntervalStr)
	if err != nil {
		logrus.Errorf("Failed to parse retry interval: %v", err)
		return nil, err
	}
	retryIntervalDuration := time.Duration(retryInterval) * time.Second

	producer, err := config.RetryKafkaConnection(brokersUrl, maxRetriesInt, retryIntervalDuration)
	if err != nil {
		logrus.Errorf("Failed to connect to Kafka: %v", err)
		return nil, err // Return the error instead of exiting
	}

	logrus.Infof("Connected to Kafka")
	return kafkaconfig.NewSaramaProducer(producer), nil
}
