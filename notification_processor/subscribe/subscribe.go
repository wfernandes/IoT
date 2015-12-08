package subscribe

import (
	"encoding/json"

	"sync"

	"github.com/wfernandes/homesec/notification_processor/mqtt"
)

const SENSORS_LIST_KEY = "/wff/v1/sp1/sensors/"

type Subscribe struct {
	mqttclient    mqtt.SimpleMQTT
	subscriptions Sensors
	outputChan    chan string
	lock          sync.Mutex
}

type Sensors struct {
	Sensors []string `json:"sensors"`
}

func New(mqttClient mqtt.SimpleMQTT, outputChan chan string) *Subscribe {
	return &Subscribe{
		mqttclient: mqttClient,
		outputChan: outputChan,
	}
}

func (s *Subscribe) Start() {
	err := s.mqttclient.Connect()
	if err != nil {
		panic(err)
	}

	getSensorList := func(dat []byte) {
		s.lock.Lock()
		err := json.Unmarshal(dat, &s.subscriptions)
		if err != nil {
			panic(err)
		}
		for _, sensorKey := range s.subscriptions.Sensors {
			s.mqttclient.Subscribe(sensorKey, s.sensorHandler)
		}
		s.lock.Unlock()
	}
	s.mqttclient.Subscribe(SENSORS_LIST_KEY, getSensorList)
}

// TODO: Do we need this function. Don't see its use if we can get the info from outputChan
func (s *Subscribe) Subscriptions() Sensors {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.subscriptions
}

func (s *Subscribe) Stop() {
	s.mqttclient.Disconnect()
}

func (s *Subscribe) sensorHandler(dat []byte) {
	s.outputChan <- string(dat)
}
