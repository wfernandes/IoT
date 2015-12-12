package main

import (
	"flag"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
	"github.com/wfernandes/homesec/broker"
	"github.com/wfernandes/homesec/sensor_processor/config"
	"github.com/wfernandes/homesec/sensor_processor/sensors"
)

var configFilePath = flag.String("config", "config/sensor.json", "Path to the Sensor Processor json config file")

func main() {
	flag.Parse()
	config, err := config.Configuration(*configFilePath)
	if err != nil {
		panic(err)
	}

	gbot := gobot.NewGobot()
	broker := broker.NewMQTTBroker("edison processor", config.BrokerUrl)
	adapter := edison.NewEdisonAdaptor("edison")
	//	adapter := testutils.NewMockAdapter("mockAdapter")
	service := sensors.Initialize(gbot, adapter, broker)

	for pin, stype := range config.Sensors {
		switch stype {
		case "touch":
			service.NewTouchSensor(pin)
		}
	}

	gbot.Start()
}
