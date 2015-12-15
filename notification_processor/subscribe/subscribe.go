package subscribe

import (
	"encoding/json"

	"sync"

	"github.com/wfernandes/iot/logging"
)

const SENSORS_LIST_KEY = "/wff/v1/sp1/sensors/"

type Subscribe struct {
	broker        Broker
	subscriptions Sensors
	outputChan    chan string
	lock          sync.Mutex
}

type Broker interface {
	Connect() error
	Subscribe(string, func([]byte))
	Disconnect()
}

type Sensors struct {
	Sensors []string `json:"sensors"`
}

func New(broker Broker, outputChan chan string) *Subscribe {
	return &Subscribe{
		broker:     broker,
		outputChan: outputChan,
	}
}

func (s *Subscribe) Start() {
	logging.Log.Info("Starting subscriber...")
	err := s.broker.Connect()
	if err != nil {
		panic(err)
	}
	logging.Log.Info("Successfully connected to broker")

	getSensorList := func(dat []byte) {
		s.lock.Lock()
		err := json.Unmarshal(dat, &s.subscriptions)
		if err != nil {
			panic(err)
		}
		logging.Log.Infof("Subscribing %d sensors", len(s.subscriptions.Sensors))
		for _, sensorKey := range s.subscriptions.Sensors {
			s.broker.Subscribe(sensorKey, s.sensorHandler)
		}
		s.lock.Unlock()
	}
	logging.Log.Info("Subscribing sensors list")
	s.broker.Subscribe(SENSORS_LIST_KEY, getSensorList)
}

func (s *Subscribe) Stop() {
	logging.Log.Info("Stopping subscriber...")
	s.broker.Disconnect()
}

func (s *Subscribe) sensorHandler(dat []byte) {
	s.outputChan <- string(dat)
}
