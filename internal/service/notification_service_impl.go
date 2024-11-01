package service

import (
	"context"
	"strconv"

	"github.com/fajaramaulana/notification-service-payment/internal/config"
	"github.com/fajaramaulana/notification-service-payment/internal/model"
	"github.com/fajaramaulana/notification-service-payment/internal/repository"
	pb "github.com/fajaramaulana/shared-proto-payment/proto/notification"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationServiceImpl struct {
	repo   repository.NotificationRepository
	config config.Config
	pb.UnimplementedNotificationServiceServer
}

func NewNotificationService(repo repository.NotificationRepository, config config.Config) pb.NotificationServiceServer {
	return &NotificationServiceImpl{
		repo:   repo,
		config: config,
	}
}

func (n *NotificationServiceImpl) SendNotification(ctx context.Context, req *pb.NotificationRequest) (*pb.NotificationResponse, error) {
	logrus.Infof("Send Notification for : %s", req.GetUserId())

	userId, err := strconv.Atoi(req.GetUserId())
	if err != nil {
		logrus.Errorf("Failed to convert user id: %s", req.GetUserId())
		return &pb.NotificationResponse{
			Status:  "false",
			Message: "Failed to convert user id",
		}, status.Error(codes.InvalidArgument, "Failed to convert user id")
	}
	userDetail, err := n.repo.GetDetailUserByUserID(userId)

	if err != nil {
		logrus.Errorf("Failed to get user detail by user id: %d", userId)
		return &pb.NotificationResponse{
			Status:  "false",
			Message: "Failed to get user detail",
		}, status.Error(codes.NotFound, "Failed to get user detail")
	}

	if userDetail == nil {
		logrus.Errorf("User detail not found by user id: %d", userId)
		return &pb.NotificationResponse{
			Status:  "false",
			Message: "User detail not found",
		}, status.Error(codes.NotFound, "User not found")
	}

	// Insert to kafka
	notification := model.NotificationToKafka{
		DetailNotificationToKafka: model.DetailNotificationToKafka{
			ID:               userDetail.ID,
			Email:            userDetail.Email,
			Title:            req.GetTitle(),
			Message:          req.GetMessage(),
			NotificationType: req.GetType(),
			RequestTime:      req.GetTimestamp(),
		},
	}

	err = n.repo.InsertNotificationKafka(notification)

	if err != nil {
		logrus.Errorf("Failed to insert notification to kafka: %v", err)
		return &pb.NotificationResponse{
			Status:  "false",
			Message: "Failed to insert notification to kafka",
		}, status.Error(codes.Internal, "Failed to insert notification to kafka")
	}

	// Insert to db
	notificationDB := model.InsertToDB{
		UserId:           userDetail.ID,
		Title:            req.GetTitle(),
		Message:          req.GetMessage(),
		NotificationType: req.GetType(),
		Status:           "not yet",
	}

	err = n.repo.InsertNotificationDB(notificationDB)
	if err != nil {
		logrus.Errorf("Failed to insert notification to db: %v", err)
		return &pb.NotificationResponse{
			Status:  "false",
			Message: "Failed to insert notification to db",
		}, status.Error(codes.Internal, "Failed to insert notification to db")
	}

	return &pb.NotificationResponse{
		Status:  "true",
		Message: "Success",
	}, nil

}
