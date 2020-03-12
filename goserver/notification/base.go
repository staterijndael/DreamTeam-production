package notification

import "dt/models"

type notificationBase struct {
	seen          *bool
	Model         *models.Notification
	receivers     []uint
	dashReceivers []uint
}

func (base *notificationBase) GetModel() *models.Notification {
	return base.Model
}

func (base *notificationBase) Seen(state bool) {
	copyOfState := state
	base.seen = &copyOfState
}

func (base *notificationBase) Receivers() []uint {
	return base.receivers
}

func (base *notificationBase) DashReceivers() []uint {
	return base.dashReceivers
}
