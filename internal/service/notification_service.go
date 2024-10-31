package service

import (
	"context"

	pb "github.com/fajaramaulana/shared-proto-payment/proto/notification"
)

type NotificationService interface {
	SendNotification(ctx context.Context, req *pb.NotificationRequest) (*pb.NotificationResponse, error)
}
