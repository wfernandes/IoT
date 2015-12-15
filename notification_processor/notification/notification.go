package notification

import (
	"github.com/wfernandes/iot/logging"
	"github.com/wfernandes/iot/notification_processor/notifiers"
)

const CLOSE_MESSAGE = "Notification Service Shutdown"

type NotificationService struct {
	notifier  notifiers.Notifier
	inputChan chan string
}

func New(notifier notifiers.Notifier, inputChan chan string) *NotificationService {

	return &NotificationService{
		notifier:  notifier,
		inputChan: inputChan,
	}
}

func (n *NotificationService) Start() {
	var err error
	logging.Log.Info("Notification service started...")
	for body := range n.inputChan {
		err = n.notifier.Notify(body)
		if err != nil {
			logging.Log.Error("Error notifying", err)
		}
	}
	err = n.notifier.Notify(CLOSE_MESSAGE)
	if err != nil {
		logging.Log.Errorf("Error sending close notification", err)
	}
}

func (n *NotificationService) Stop() {
	close(n.inputChan)
}
