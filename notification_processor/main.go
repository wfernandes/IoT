package main

import (
	"flag"

	"time"

	"github.com/wfernandes/iot/broker"
	"github.com/wfernandes/iot/event"
	"github.com/wfernandes/iot/logging"
	"github.com/wfernandes/iot/notification_processor/config"
	"github.com/wfernandes/iot/notification_processor/notification"
	"github.com/wfernandes/iot/notification_processor/notifiers"
	"github.com/wfernandes/iot/notification_processor/subscribe"
)

var configFilePath = flag.String("config", "config/homesec.json", "Path to the HomeSec json config file")

func main() {

	flag.Parse()
	config, err := config.Configuration(*configFilePath)
	if err != nil {
		panic(err)
	}
	logging.SetLogLevel(config.LogLevel)

	alertChan := make(chan *event.Event, 100)
	mqttBroker := broker.NewMQTTBroker("wff_notification", config.BrokerUrl)
	subscriber := subscribe.New(mqttBroker, alertChan)
	// Subscribe to all available sensor keys
	go subscriber.Start()

	// Initialize twilio notification service
	notifier := notifiers.NewTwilio(config.TwilioAccountSid, config.TwilioAuthToken, config.TwilioFromPhone, config.To)

	service := notification.New(notifier, alertChan, time.Duration(config.NotificationIntervalMinutes)*time.Minute)
	service.Start()
}
