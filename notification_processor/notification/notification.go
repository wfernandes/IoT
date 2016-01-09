package notification

import (
	"time"

	"github.com/wfernandes/iot/event"
	"github.com/wfernandes/iot/logging"
	"github.com/wfernandes/iot/notification_processor/notifiers"
)

const CLOSE_MESSAGE = "Notification Service Shutdown"

type NotificationService struct {
	notifier        notifiers.Notifier
	sensorsNotified map[string]time.Time
	inputChan       chan *event.Event
	duration        time.Duration
}

func New(notifier notifiers.Notifier, inputChan chan *event.Event, notificationPeriod time.Duration) *NotificationService {
	return &NotificationService{
		notifier:        notifier,
		inputChan:       inputChan,
		sensorsNotified: make(map[string]time.Time),
		duration:        notificationPeriod,
	}
}

func (n *NotificationService) Start() {
	var err error
	logging.Log.Info("Notification service started...")
	for event := range n.inputChan {
		logging.Log.Debugf("Received event: %s", event.Name)
		err = n.notify(event)
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

func (n *NotificationService) notify(evnt *event.Event) error {
	lastNotified, ok := n.sensorsNotified[evnt.Name]
	if !ok {
		n.sensorsNotified[evnt.Name] = time.Now()
		logging.Log.Debugf("Notifying event for %s", evnt.Name)
		return n.notifier.Notify(evnt.Data)
	}
	if lastNotified.Add(n.duration).After(time.Now()) {
		return nil
	} else {
		n.sensorsNotified[evnt.Name] = time.Now()
		logging.Log.Debugf("Notifying event for %s", evnt.Name)
		return n.notifier.Notify(evnt.Data)
	}

}
