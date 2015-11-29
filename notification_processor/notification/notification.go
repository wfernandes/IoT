package notification

import (
	"github.com/cloudfoundry/gosteno"
	"github.com/wfernandes/homesec/notification_processor/notifiers"
)

var CLOSE_MESSAGE = "Notification Service Shutdown"

type NotificationService struct {
	notifier  notifiers.Notifier
	inputChan chan string
	logger    *gosteno.Logger
}

func New(notifier notifiers.Notifier, inputChan chan string, logger *gosteno.Logger) *NotificationService {

	return &NotificationService{
		notifier:  notifier,
		inputChan: inputChan,
		logger:    logger,
	}
}

func (n *NotificationService) Start() {
	var err error
	n.logger.Info("Notification service started...")
	for body := range n.inputChan {
		n.logger.Info(body)
		err = n.notifier.Notify(body)
		if err != nil {
			n.logger.Warn(err.Error())
		}
	}
	err = n.notifier.Notify(CLOSE_MESSAGE)
	if err != nil {
		n.logger.Errorf("Error sending close notification: %s", err.Error())
	}
}

func (n *NotificationService) Stop() {
	close(n.inputChan)
}
