package sensors

import (
	"encoding/json"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/wfernandes/iot/event"
	"github.com/wfernandes/iot/logging"
)

const SENSOR_KEY = "/wff/v1/sp1/"
const SENSORS_LIST_KEY = "/wff/v1/sp1/sensors/"
const SOUND_THRESHOLD = 129

type AdaptorReader interface {
	gobot.Adaptor
	AnalogRead(string) (val int, err error)
	DigitalRead(string) (val int, err error)
}

type SensorService struct {
	gobot   *gobot.Gobot
	adapter AdaptorReader
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

func Initialize(gobot *gobot.Gobot, adapter AdaptorReader, broker Broker) *SensorService {
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
	eventTouched := &event.Event{
		Name: name,
		Data: "touched",
	}
	eventReleased := &event.Event{
		Name: name,
		Data: "released",
	}

	s.publishSensorList(name)

	work := func() {
		touchSensor.On(touchSensor.Event(gpio.ButtonPush), func(data interface{}) {
			s.publish(eventTouched)
		})

		touchSensor.On(touchSensor.Event(gpio.ButtonRelease), func(data interface{}) {
			s.publish(eventReleased)
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

func (s *SensorService) NewLightSensor(pin string) {
	lightSensor := gpio.NewGroveLightSensorDriver(s.adapter, "light", pin)
	name := "lightsensor" + pin
	eventLight := &event.Event{
		Name: name,
		Data: "light detected",
	}
	s.publishSensorList(name)

	work := func() {
		lightSensor.On(lightSensor.Event(gpio.Data), func(data interface{}) {
			logging.Log.Debugf("light data %d", data)
			s.publish(eventLight)
		})
	}

	robot := gobot.NewRobot(name,
		[]gobot.Connection{s.adapter},
		[]gobot.Device{lightSensor},
		work,
	)
	s.gobot.AddRobot(robot)
	logging.Log.Infof("Added sensor %s", name)
}

func (s *SensorService) NewSoundSensor(pin string) {

	soundSensor := gpio.NewGroveSoundSensorDriver(s.adapter, "sound", pin)
	name := "soundsensor" + pin
	eventSound := &event.Event{
		Name: name,
		Data: "noise detected",
	}
	s.publishSensorList(name)

	work := func() {
		soundSensor.On(soundSensor.Event(gpio.Data), func(data interface{}) {
			logging.Log.Debugf("sound data %d", data)
			if data.(int) > SOUND_THRESHOLD {
				s.publish(eventSound)
			}
		})
	}

	robot := gobot.NewRobot(name,
		[]gobot.Connection{s.adapter},
		[]gobot.Device{soundSensor},
		work,
	)
	s.gobot.AddRobot(robot)
	logging.Log.Infof("Added sensor %s", name)
}

func (s *SensorService) publish(event *event.Event) {
	if s.broker.IsConnected() {
		eventBytes, _ := json.Marshal(event)
		logging.Log.Debugf("Publishing %s", SENSOR_KEY+event.Name)
		s.broker.Publish(SENSOR_KEY+event.Name, eventBytes)
	}
}

func (s *SensorService) buildSensorList(sensorName string) ([]byte, error) {
	s.sensors.List = append(s.sensors.List, SENSOR_KEY+sensorName)
	return json.Marshal(s.sensors)
}

func (s *SensorService) publishSensorList(sensorName string) {
	sensorList, err := s.buildSensorList(sensorName)
	if err != nil {
		logging.Log.Errorf("Error building sensor list: %s", err.Error())
	}
	s.broker.Publish(SENSORS_LIST_KEY, sensorList)
}
