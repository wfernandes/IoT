package main

import (
	"flag"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
	"github.com/wfernandes/iot/broker"
	"github.com/wfernandes/iot/logging"
	"github.com/wfernandes/iot/sensor_processor/config"
	"github.com/wfernandes/iot/sensor_processor/sensors"
)

var configFilePath = flag.String("config", "config/sensor.json", "Path to the Sensor Processor json config file")

func main() {
	flag.Parse()
	config, err := config.Configuration(*configFilePath)
	if err != nil {
		panic(err)
	}
	logging.SetLogLevel(config.LogLevel)

	broker := broker.NewMQTTBroker("edison processor", config.BrokerUrl)
	err = broker.Connect()
	if err != nil {
		panic(err)
	}
	logging.Log.Info("Successfully connected to broker")

	gbot := gobot.NewGobot()
	adapter := edison.NewEdisonAdaptor("edison")
	service := sensors.Initialize(gbot, adapter, broker)

	for pin, stype := range config.Sensors {
		switch stype {
		case "touch":
			service.NewTouchSensor(pin)
		case "sound":
			service.NewSoundSensor(pin)
		}
	}

	logging.Log.Info("Starting gobot bot...")
	gbot.Start()
}
