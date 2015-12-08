package subscribe_test

import (
	"fmt"

	"github.com/wfernandes/homesec/notification_processor/subscribe"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Subscriber", func() {

	var (
		outputChan chan string
	)
	Context("without errors", func() {

		BeforeEach(func() {
			outputChan = make(chan string, 100)
		})

		It("reads initial list of sensor keys", func() {
			mockMqtt := &MockMQTTClient{
				ConnectErr: nil,
			}
			s := subscribe.New(mockMqtt, outputChan)
			s.Start()
			Eventually(func() []string { return s.Subscriptions().Sensors }).Should(HaveLen(2))
			Expect(s.Subscriptions().Sensors).To(ContainElement("/wff/v1/sp/touchsensor"))
			Expect(s.Subscriptions().Sensors).To(ContainElement("/wff/v1/sp/soundsensor"))
		})

		It("disconnects mqtt client when stop is called", func() {
			mockMqtt := &MockMQTTClient{
				ConnectErr: nil,
			}
			s := subscribe.New(mockMqtt, outputChan)
			s.Stop()

			Expect(mockMqtt.DisconnectCalled).To(BeTrue())
		})

		It("subscribes to sensor keys from sensor list", func() {

			mockMQTT := &MockMQTTClient{
				ConnectErr: nil,
			}
			s := subscribe.New(mockMQTT, outputChan)
			s.Start()
			Eventually(outputChan).Should(Receive(Equal("/wff/v1/sp/touchsensor")))
			Eventually(outputChan).Should(Receive(Equal("/wff/v1/sp/soundsensor")))

		})
	})

	Context("with errors", func() {
		It("panics", func() {
			mockMqtt := &MockMQTTClient{
				ConnectErr: fmt.Errorf("some error"),
			}
			s := subscribe.New(mockMqtt, nil)
			Expect(s.Start).To(Panic())
		})
	})
})

type MockMQTTClient struct {
	ConnectErr       error
	DisconnectCalled bool
}

func (m *MockMQTTClient) Connect() error {
	if m.ConnectErr != nil {
		return m.ConnectErr
	}
	return nil
}

func (m *MockMQTTClient) Disconnect() {
	m.DisconnectCalled = true
}

func (m *MockMQTTClient) IsConnected() bool {
	if m.ConnectErr != nil {
		return false
	}
	return true
}

func (m *MockMQTTClient) Subscribe(event string, f func([]byte)) {
	if event == subscribe.SENSORS_LIST_KEY {
		data := []byte(`{"sensors":["/wff/v1/sp/touchsensor", "/wff/v1/sp/soundsensor"]}`)
		f(data)
	} else {
		// for now pass in the event name as the data
		f([]byte(event))
	}

}

func (m *MockMQTTClient) Publish(topic string, data []byte) {
	// do nothing
}
