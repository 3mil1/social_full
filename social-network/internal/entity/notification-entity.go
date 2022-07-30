package entity

type Notification struct {
	NotificationObjID int
	ObjectID          int
	ActorID           string
	ReceiverID        string
	NotificationType  int
}
