package subscribe

import (
	"encoding/json"

	"sync"
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
	err := s.broker.Connect()
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
			s.broker.Subscribe(sensorKey, s.sensorHandler)
		}
		s.lock.Unlock()
	}
	s.broker.Subscribe(SENSORS_LIST_KEY, getSensorList)
}

// TODO: Do we need this function. Don't see its use if we can get the info from outputChan
func (s *Subscribe) Subscriptions() Sensors {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.subscriptions
}

func (s *Subscribe) Stop() {
	s.broker.Disconnect()
}

func (s *Subscribe) sensorHandler(dat []byte) {
	s.outputChan <- string(dat)
}
