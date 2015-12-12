package sensors

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

const SENSOR_KEY = "/wff/v1/sp1/"

type SensorService struct {
	gobot   *gobot.Gobot
	adapter gpio.DigitalReader
	broker  Broker
}

type Broker interface {
	Connect() error
	Publish(string, []byte)
	Subscribe(string, func([]byte))
	IsConnected() bool
	Disconnect()
}

// TODO: Add logging
func Initialize(gobot *gobot.Gobot, adapter gpio.DigitalReader, broker Broker) *SensorService {
	return &SensorService{
		gobot:   gobot,
		adapter: adapter,
		broker:  broker,
	}
}

func (s *SensorService) NewTouchSensor(pin string) {

	touchSensor := gpio.NewGroveTouchDriver(s.adapter, "touch", pin)
	name := "touchsensor" + pin

	work := func() {
		gobot.On(touchSensor.Event(gpio.Push), func(data interface{}) {
			s.publish(name, "touched")
		})

		gobot.On(touchSensor.Event(gpio.Release), func(data interface{}) {
			s.publish(name, "released")
		})
	}

	robot := gobot.NewRobot(name,
		[]gobot.Connection{s.adapter},
		[]gobot.Device{touchSensor},
		work,
	)
	s.gobot.AddRobot(robot)
}

func (s *SensorService) publish(sensorName string, value string) {
	if s.broker.IsConnected() {
		s.broker.Publish(SENSOR_KEY+sensorName, []byte(value))
	}
}
