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
		broker     *mockBroker
	)
	Context("without errors", func() {

		BeforeEach(func() {
			outputChan = make(chan string, 100)
			broker = newMockBroker()
		})

		It("reads initial list of sensor keys", func() {
			broker.ConnectOutput.ret0 <- nil
			s := subscribe.New(broker, outputChan)
			s.Start()
			Eventually(broker.ConnectCalled).Should(Receive(BeTrue()))
			Eventually(broker.SubscribeCalled).Should(Receive(BeTrue()))
			Eventually(broker.SubscribeInput.arg0).Should(Receive(Equal(subscribe.SENSORS_LIST_KEY)))
			data := []byte(`{"sensors":["/wff/v1/sp1/touchsensor", "/wff/v1/sp1/soundsensor"]}`)
			getSensorList := <-broker.SubscribeInput.arg1
			getSensorList(data)
			Eventually(broker.SubscribeCalled).Should(Receive(BeTrue()))
			Eventually(broker.SubscribeInput.arg0).Should(Receive(Equal("/wff/v1/sp1/touchsensor")))
			Eventually(broker.SubscribeInput.arg0).Should(Receive(Equal("/wff/v1/sp1/soundsensor")))
		})

		It("disconnects mqtt client when stop is called", func() {
			s := subscribe.New(broker, outputChan)
			s.Stop()

			Eventually(broker.DisconnectCalled).Should(Receive(BeTrue()))
		})

		It("subscribes to sensor keys from sensor list", func() {
			broker.ConnectOutput.ret0 <- nil
			s := subscribe.New(broker, outputChan)
			s.Start()
			data := []byte(`{"sensors":["/wff/v1/sp1/touchsensor", "/wff/v1/sp1/soundsensor"]}`)
			// get the function passed in and invoke it
			getSensorList := <-broker.SubscribeInput.arg1
			getSensorList(data)

			Eventually(broker.SubscribeCalled).Should(Receive(BeTrue()))
			touchSensorHandler := <-broker.SubscribeInput.arg1
			touchSensorHandler([]byte("some touch data"))
			Eventually(outputChan).Should(Receive(Equal("some touch data")))

			Eventually(broker.SubscribeCalled).Should(Receive(BeTrue()))
			soundSensorHandler := <-broker.SubscribeInput.arg1
			soundSensorHandler([]byte("some sound data"))
			Eventually(outputChan).Should(Receive(Equal("some sound data")))

		})
	})

	Context("with errors", func() {
		It("panics", func() {
			broker.ConnectOutput.ret0 <- fmt.Errorf("some error")
			s := subscribe.New(broker, nil)
			Expect(s.Start).To(Panic())
		})
	})
})
