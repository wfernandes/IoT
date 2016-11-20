package sensors_test

import (
	"github.com/wfernandes/iot/sensor_processor/sensors"

	"log"

	"encoding/json"

	"github.com/hybridgroup/gobot"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wfernandes/iot/event"
	"github.com/wfernandes/iot/sensor_processor/testutils"
)

var _ = Describe("Sensors", func() {

	var (
		service *sensors.SensorService
		broker  *mockBroker
		gbot    *gobot.Gobot
	)

	BeforeEach(func() {
		// Remove line below to get more debug info
		log.SetOutput(&testutils.NullReadWriteCloser{})

		gbot = gobot.NewGobot()
		adapter := newMockAdapter("mock adapter")
		broker = newMockBroker()
		service = sensors.Initialize(gbot, adapter, broker)
	})

	It("publishes sensor name to sensor list", func() {
		service.NewTouchSensor("4")
		go gbot.Start()
		broker.IsConnectedOutput.ret0 <- true
		Eventually(broker.PublishCalled).Should(Receive(BeTrue()))
		Eventually(broker.PublishInput.arg0).Should(Receive(Equal(sensors.SENSORS_LIST_KEY)))
		Eventually(broker.PublishInput.arg1).Should(Receive(MatchJSON(`{"sensors":["/wff/v1/sp1/touchsensor4"]}`)))
	})

	Context("Touch", func() {

		It("adds touch robot to gobot", func() {
			service.NewTouchSensor("2")
			Eventually(gbot.Robots().Len()).Should(Equal(1))
			Expect(gbot.Robot("touchsensor2")).ToNot(BeNil())
		})

		It("publishes alert on touch push event", func() {
			service.NewTouchSensor("3")
			evt := event.Event{
				Name: "touchsensor3",
				Data: "touched",
			}
			eventTouched, _ := json.Marshal(evt)
			evt = event.Event{
				Name: "touchsensor3",
				Data: "released",
			}
			eventReleased, _ := json.Marshal(evt)

			go gbot.Start()
			broker.IsConnectedOutput.ret0 <- true
			Eventually(broker.IsConnectedCalled).Should(Receive(BeTrue()))
			Eventually(broker.PublishCalled).Should(Receive(BeTrue()))
			Eventually(broker.PublishInput.arg0).Should(Receive(Equal("/wff/v1/sp1/touchsensor3")))
			Eventually(broker.PublishInput.arg1).Should(Receive(Equal(eventTouched)))
			// Doing this again, since the channel gets flushed
			broker.IsConnectedOutput.ret0 <- true
			Eventually(broker.IsConnectedCalled).Should(Receive(BeTrue()))
			Eventually(broker.PublishCalled).Should(Receive(BeTrue()))
			Eventually(broker.PublishInput.arg0).Should(Receive(Equal("/wff/v1/sp1/touchsensor3")))
			Eventually(broker.PublishInput.arg1).Should(Receive(Equal(eventReleased)))
		})
	})

	Context("Light", func() {

		It("adds light grove sensor robot to gobot", func() {
			service.NewLightSensor("2")
			Eventually(gbot.Robots().Len()).Should(Equal(1))
			Expect(gbot.Robot("lightsensor2")).ToNot(BeNil())
		})

		It("publishes alert on light event", func() {
			service.NewLightSensor("2")
			evt := event.Event{
				Name: "lightsensor2",
				Data: "light detected",
			}
			evtLight, _ := json.Marshal(evt)

			go gbot.Start()
			broker.IsConnectedOutput.ret0 <- true
			Eventually(broker.IsConnectedCalled).Should(Receive(BeTrue()))
			Eventually(broker.PublishCalled).Should(Receive(BeTrue()))
			Eventually(broker.PublishInput.arg0).Should(Receive(Equal("/wff/v1/sp1/lightsensor2")))
			Eventually(broker.PublishInput.arg1).Should(Receive(Equal(evtLight)))

		})
	})

	Context("Sound", func() {
		It("adds sound robot to gobot", func() {
			service.NewSoundSensor("4")
			go gbot.Start()
			broker.IsConnectedOutput.ret0 <- true
			Eventually(broker.PublishCalled).Should(Receive(BeTrue()))
			Eventually(broker.PublishInput.arg0).Should(Receive(Equal(sensors.SENSORS_LIST_KEY)))
			Eventually(broker.PublishInput.arg1).Should(Receive(MatchJSON(`{"sensors":["/wff/v1/sp1/soundsensor4"]}`)))
		})

		It("publishes alert on acoustic event ", func() {
			service.NewSoundSensor("5")
			evt := event.Event{
				Name: "soundsensor5",
				Data: "noise detected",
			}
			eventSound, _ := json.Marshal(evt)
			go gbot.Start()
			broker.IsConnectedOutput.ret0 <- true
			Eventually(broker.IsConnectedCalled).Should(Receive(BeTrue()))
			Eventually(broker.PublishCalled).Should(Receive(BeTrue()))
			Eventually(broker.PublishInput.arg0).Should(Receive(Equal("/wff/v1/sp1/soundsensor5")))
			Eventually(broker.PublishInput.arg1).Should(Receive(Equal(eventSound)))
		})
	})

})
