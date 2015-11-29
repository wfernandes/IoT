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
		service  *sensors.SensorService
		dataChan chan string
		gbot     *gobot.Gobot
	)

	BeforeEach(func() {
		// Remove line below to get more debug info
		log.SetOutput(&testutils.NullReadWriteCloser{})

		gbot = gobot.NewGobot()
		adapter := testutils.NewMockAdapter("mock adapter")
		dataChan = make(chan string)
		service = sensors.Initialize(gbot, adapter, dataChan)
	})

	It("adds touch robot to gobot", func() {
		service.NewTouchSensor("2")
		Eventually(gbot.Robots().Len()).Should(Equal(1))
		Expect(gbot.Robot("TouchSensor-2")).ToNot(BeNil())
	})

	It("sends alert on dataChan on touch push event", func() {
		service.NewTouchSensor("3")
		go gbot.Start()

		Eventually(dataChan).Should(Receive(Equal("TouchSensor-3 Touched")))
		Eventually(dataChan).Should(Receive(Equal("TouchSensor-3 Released")))
	})

})
