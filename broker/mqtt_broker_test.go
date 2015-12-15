package broker_test

import (
	"github.com/wfernandes/iot/broker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MQTT Broker", func() {

	var mc *broker.MQTTBroker

	BeforeEach(func() {
		mc = broker.NewMQTTBroker("some_client_id", "tcp://test.mosquitto.org:1883")
	})

	It("isConnected returns false if not connected yet", func() {
		Expect(mc.IsConnected()).To(BeFalse())
	})

	It("does not panic if it tries to publish without connecting first", func() {
		Expect(func() { mc.Publish("some_key", []byte("some value")) }).ToNot(Panic())
	})

	It("does not panic if it tries to subscribe without connecting first", func() {
		f := func([]byte) { /*do nothing*/ }
		Expect(func() { mc.Subscribe("some_key", f) }).ToNot(Panic())
	})

	Context("cannot connect", func() {
		It("returns error if cannont connect", func() {
			mc := broker.NewMQTTBroker("some_client_id", "test.mosquitto.org:1883")

			err := mc.Connect()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("connects", func() {
		BeforeEach(func() {
			mc = broker.NewMQTTBroker("some_client_id", "tcp://test.mosquitto.org:1883")
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
