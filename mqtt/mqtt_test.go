package mqtt_test

import (
	"github.com/wfernandes/homesec/notification_processor/mqtt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mqtt", func() {

	var (
		mc *mqtt.MQTTClient
	)

	Context("cannot connect", func() {
		It("returns error if cannont connect", func() {
			mc := mqtt.New("some_client_id", "test.mosquitto.org:1883")

			err := mc.Connect()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("connects", func() {
		BeforeEach(func() {
			mc = mqtt.New("some_client_id", "tcp://test.mosquitto.org:1883")
		})

		It("connects and disconnects", func() {
			err := mc.Connect()
			Expect(err).ToNot(HaveOccurred())

			Expect(mc.IsConnected()).To(BeTrue())
			mc.Disconnect()
			Expect(mc.IsConnected()).To(BeFalse())
		})

		It("pub subs", func() {
			outChan := make(chan []byte)
			err := mc.Connect()
			Expect(err).ToNot(HaveOccurred())
			// subscribe to the event
			f := func(s []byte) {
				outChan <- s
			}
			mc.Subscribe("/test/hello/", f)
			// publish to the event
			for i := 0; i < 5; i++ {
				mc.Publish("/test/hello/", []byte("world123"))
			}

			Eventually(outChan).Should(Receive(BeEquivalentTo("world123")))

			mc.Disconnect()
		})
	})

})
