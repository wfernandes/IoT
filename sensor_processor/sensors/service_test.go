package sensors_test

import (
	"github.com/wfernandes/homesec/sensor_processor/sensors"

	"log"

	"github.com/hybridgroup/gobot"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wfernandes/homesec/sensor_processor/testutils"
)

var _ = Describe("Touch", func() {

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

	It("adds touch robot to gobot", func() {
		service.NewTouchSensor("2")
		Eventually(gbot.Robots().Len()).Should(Equal(1))
		Expect(gbot.Robot("touchsensor2")).ToNot(BeNil())
	})

	It("sends publishes alert on touch push event", func() {
		service.NewTouchSensor("3")
		go gbot.Start()
		broker.IsConnectedOutput.ret0 <- true
		Eventually(broker.IsConnectedCalled).Should(Receive(BeTrue()))
		Eventually(broker.PublishCalled).Should(Receive(BeTrue()))
		Eventually(broker.PublishInput.arg0).Should(Receive(Equal("/wff/v1/sp1/touchsensor3")))
		Eventually(broker.PublishInput.arg1).Should(Receive(BeEquivalentTo("touched")))
		// Doing this again, since the channel gets flushed
		broker.IsConnectedOutput.ret0 <- true
		Eventually(broker.IsConnectedCalled).Should(Receive(BeTrue()))
		Eventually(broker.PublishCalled).Should(Receive(BeTrue()))
		Eventually(broker.PublishInput.arg0).Should(Receive(Equal("/wff/v1/sp1/touchsensor3")))
		Eventually(broker.PublishInput.arg1).Should(Receive(BeEquivalentTo("released")))
	})

})
