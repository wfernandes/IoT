package main

import (
	"flag"

	"github.com/cloudfoundry/gosteno"
	"github.com/wfernandes/homesec/broker"
	"github.com/wfernandes/homesec/notification_processor/config"
	"github.com/wfernandes/homesec/notification_processor/notification"
	"github.com/wfernandes/homesec/notification_processor/notifiers"
	"github.com/wfernandes/homesec/notification_processor/subscribe"
)

var configFilePath = flag.String("config", "config/homesec.json", "Path to the HomeSec json config file")

func main() {

	flag.Parse()
	config, err := config.Configuration(*configFilePath)
	if err != nil {
		panic(err)
	}
	alertChan := make(chan string, 100)

	mqttBroker := broker.NewMQTTBroker("wff_notification", config.BrokerUrl)
	subscriber := subscribe.New(mqttBroker, alertChan)
	// Subscribe to all available sensor keys
	go subscriber.Start()

	// Initialize twilio notification service
	notifier := notifiers.NewTwilio(config.TwilioAccountSid, config.TwilioAuthToken, config.TwilioFromPhone, config.To)

	logger := gosteno.NewLogger("Notification Service")
	service := notification.New(notifier, alertChan, logger)
	service.Start()
}
