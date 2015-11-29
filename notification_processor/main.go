package main

import (
	"flag"
	"fmt"

	"github.com/cloudfoundry/gosteno"
	"github.com/pivotal-golang/localip"
	"github.com/wfernandes/homesec/notification_processor/config"
	"github.com/wfernandes/homesec/notification_processor/listeners"
	"github.com/wfernandes/homesec/notification_processor/notification"
	"github.com/wfernandes/homesec/notification_processor/notifiers"
)

var configFilePath = flag.String("config", "config/homesec.json", "Path to the HomeSec json config file")

func main() {

	flag.Parse()
	config, err := config.Configuration(*configFilePath)
	if err != nil {
		panic(err)
	}
	// Initialize twilio notification service
	notifier := notifiers.NewTwilio(config.TwilioAccountSid, config.TwilioAuthToken, config.TwilioFromPhone, config.To)

	alertChan := make(chan string)
	logger := gosteno.NewLogger("Notification Service")
	service := notification.New(notifier, alertChan, logger)
	go service.Start()

	// Send stuff to alertChan once we get an alert from sensor.
	ipAddress, err := localip.LocalIP()
	if err != nil {
		panic(err)
	}
	address := fmt.Sprintf("%s:%d", ipAddress, config.Port)
	logger.Infof("Listening for incoming alerts on %s", address)
	listener := listeners.New(address, alertChan, logger)
	listener.Start() // Blocking
}
