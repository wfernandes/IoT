package sensors

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

type SensorService struct {
	gobot    *gobot.Gobot
	adapter  gpio.DigitalReader
	dataChan chan string
}

// TODO: Add logging
func Initialize(gobot *gobot.Gobot, adapter gpio.DigitalReader, dataCh chan string) *SensorService {
	return &SensorService{
		gobot:    gobot,
		adapter:  adapter,
		dataChan: dataCh,
	}
}

func (s *SensorService) NewTouchSensor(pin string) {

	touchSensor := gpio.NewGroveTouchDriver(s.adapter, "touch", pin)
	name := "TouchSensor-" + pin

	work := func() {
		gobot.On(touchSensor.Event(gpio.Push), func(data interface{}) {
			s.dataChan <- name + " Touched"
		})

		gobot.On(touchSensor.Event(gpio.Release), func(data interface{}) {
			s.dataChan <- name + " Released"
		})
	}

	robot := gobot.NewRobot(name,
		[]gobot.Connection{s.adapter},
		[]gobot.Device{touchSensor},
		work,
	)
	s.gobot.AddRobot(robot)
}
