package writers_test

import (
	"github.com/wfernandes/homesec/sensor_processor/writers"

	"net"

	"log"

	"github.com/hybridgroup/gobot"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/localip"
	"github.com/wfernandes/homesec/sensor_processor/sensors"
	"github.com/wfernandes/homesec/sensor_processor/testutils"
)

var _ = Describe("Writers", func() {

	var (
		ip       string
		dataChan chan string
	)

	BeforeEach(func() {
		// Remove line below to get more debug info
		log.SetOutput(&testutils.NullReadWriteCloser{})
		gbot := gobot.NewGobot()
		adapter := testutils.NewMockAdapter("test adapter")
		dataChan = make(chan string)
		service := sensors.Initialize(gbot, adapter, dataChan)
		service.NewTouchSensor("2")
		go gbot.Start()
		ip, _ = localip.LocalIP()

	})

	It("reads sensor alerts from udp connection", func() {

		address := net.JoinHostPort(ip, "8080")
		writer := writers.New(address, dataChan)
		go writer.Start()

		conn, err := net.ListenPacket("udp", address)
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() string {
			readBuffer := make([]byte, 1000)
			readCount, _, err := conn.ReadFrom(readBuffer)
			Expect(err).ToNot(HaveOccurred())
			readData := make([]byte, readCount)
			copy(readData, readBuffer[:readCount])
			return string(readData)
		}).Should(Equal("TouchSensor-2 Touched"))
	})

	It("panics if unable to connect", func() {
		address := "bla:8888"
		writer := writers.New(address, dataChan)
		Expect(writer.Start).To(Panic())
	})

})
