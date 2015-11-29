package listeners_test

import (
	"net"

	"github.com/wfernandes/homesec/notification_processor/listeners"

	. "github.com/onsi/ginkgo"
	gconf "github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"

	"strconv"

	"time"

	"github.com/cloudfoundry/gosteno"
)

var _ = Describe("Listener", func() {

	var (
		listener   *listeners.UDPListener
		address    string
		outputChan chan string
		logger     *gosteno.Logger
	)

	BeforeEach(func() {
		port := 8080 + gconf.GinkgoConfig.ParallelNode
		address = net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
		outputChan = make(chan string)
		logger = gosteno.NewLogger("Listener logger")
		listener = listeners.New(address, outputChan, logger)
	})

	Context("start", func() {
		It("reads from connection and sends data to outputChan", func() {

			go listener.Start()

			conn, err := net.Dial("udp", address)
			Expect(err).ToNot(HaveOccurred())

			doneChan := make(chan struct{})
			go func() {
				t := time.NewTicker(100 * time.Millisecond)
				for {
					select {
					case <-t.C:
						_, err = conn.Write([]byte("sensor alert"))
						Expect(err).ToNot(HaveOccurred())
					case <-doneChan:
						return
					}
				}
			}()

			var received string
			Eventually(outputChan, 2).Should(Receive(&received))
			Expect(received).To(Equal("sensor alert"))
			close(doneChan)
		})

		It("closes the dataChan if there is a network error", func() {
			address = "not a real address"
			listener = listeners.New(address, outputChan, logger)

			Expect(outputChan).ToNot(BeClosed())
			go listener.Start()
			Eventually(outputChan).Should(BeClosed())
		})
	})

	Context("Stop", func() {
		It("closes dataChan when stopped", func() {
			Expect(outputChan).ToNot(BeClosed())
			listener.Stop()
			Expect(outputChan).To(BeClosed())
		})
	})

})
