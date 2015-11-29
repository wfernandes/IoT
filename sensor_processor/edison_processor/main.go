package main

import (
	"flag"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
	"github.com/wfernandes/homesec/sensor_processor/config"
	"github.com/wfernandes/homesec/sensor_processor/sensors"
	"github.com/wfernandes/homesec/sensor_processor/writers"
)

var configFilePath = flag.String("config", "config/sensor.json", "Path to the Sensor Processor json config file")

func main() {
	flag.Parse()
	config, err := config.Configuration(*configFilePath)
	if err != nil {
		panic(err)
	}

	gbot := gobot.NewGobot()
	dataChan := make(chan string)
	adapter := edison.NewEdisonAdaptor("edison")
	//	adapter := testutils.NewMockAdapter("mockAdapter")
	service := sensors.Initialize(gbot, adapter, dataChan)

	for pin, stype := range config.Sensors {
		switch stype {
		case "touch":
			service.NewTouchSensor(pin)
		}
	}

	writer := writers.New(config.NotifierUrl, dataChan)
	go writer.Start()

	gbot.Start()
	close(dataChan)
}
