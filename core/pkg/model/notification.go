package model

type StateChangeNotificationType string

const (
	NotificationDelete StateChangeNotificationType = "delete"
	NotificationCreate StateChangeNotificationType = "write"
	NotificationUpdate StateChangeNotificationType = "update"
)

type StateChangeNotification struct {
	Type    StateChangeNotificationType `json:"type"`
	Source  string                      `json:"source"`
	FlagKey string                      `json:"flagKey"`
}
