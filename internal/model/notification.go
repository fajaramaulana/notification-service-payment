package model

type Notification struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Title            string `json:"title"`
	Message          string `json:"message"`
	Status           string `json:"status"`
	NotificationType string `json:"notification_type"`
}

type NotificationToKafka struct {
	DetailNotificationToKafka DetailNotificationToKafka `json:"detail_notification_to_kafka"`
}

type DetailNotificationToKafka struct {
	ID               int    `json:"id"`
	Email            string `json:"email"`
	Title            string `json:"title"`
	Message          string `json:"message"`
	NotificationType string `json:"notification_type"`
	RequestTime      string `json:"request_time"`
	IpAddress        string `json:"ip_address"`
	UserAgent        string `json:"user_agent"`
	TypeId           int    `json:"type_id"`
}

type InsertToKafka struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Time    string `json:"time"`
	Data    DetailNotificationToKafka
}

type DetailUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FName    string `json:"f_name"`
	LName    string `json:"l_name"`
	Phone    string `json:"phone"`
}

type InsertToDB struct {
	UserId           int    `json:"user_id"`
	Title            string `json:"title"`
	Message          string `json:"message"`
	NotificationType string `json:"notification_type"`
	Status           string `json:"status"`
}
