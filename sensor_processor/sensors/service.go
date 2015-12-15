package sensors

import (
	"encoding/json"
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/wfernandes/iot/logging"
)

const SENSOR_KEY = "/wff/v1/sp1/"
const SENSORS_LIST_KEY = "/wff/v1/sp1/sensors/"

type SensorService struct {
	gobot   *gobot.Gobot
	adapter gpio.DigitalReader
	broker  Broker
	sensors Sensors
}

type Sensors struct {
	List []string `json:"sensors"`
}

type Broker interface {
	Connect() error
	Publish(string, []byte)
	Subscribe(string, func([]byte))
	IsConnected() bool
	Disconnect()
}

func Initialize(gobot *gobot.Gobot, adapter gpio.DigitalReader, broker Broker) *SensorService {
	return &SensorService{
		gobot:   gobot,
		adapter: adapter,
		broker:  broker,
		sensors: Sensors{},
	}
}

func (s *SensorService) NewTouchSensor(pin string) {

	touchSensor := gpio.NewGroveTouchDriver(s.adapter, "touch", pin)
	name := "touchsensor" + pin

	sensorList, err := s.buildSensorList(name)
	if err != nil {
		fmt.Printf("Error building sensor list: %s", err.Error())
	}
	s.broker.Publish(SENSORS_LIST_KEY, sensorList)

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
	logging.Log.Infof("Added sensor %s", name)
}

func (s *SensorService) publish(sensorName string, value string) {
	if s.broker.IsConnected() {
		logging.Log.Debugf("Publishing %s", SENSOR_KEY+sensorName)
		s.broker.Publish(SENSOR_KEY+sensorName, []byte(value))
	}
}

func (s *SensorService) buildSensorList(sensorName string) ([]byte, error) {
	s.sensors.List = append(s.sensors.List, SENSOR_KEY+sensorName)
	return json.Marshal(s.sensors)
}
