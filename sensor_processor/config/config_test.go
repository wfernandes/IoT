package config_test

import (
	"github.com/wfernandes/iot/sensor_processor/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wfernandes/iot/logging"
)

var _ = Describe("Config", func() {

	It("reads sensor configs from config file", func() {

		jsonStr := []byte(`{
			"Sensors": {
				"2": "touch",
				"3": "sound"
			},
			"NotifierUrl": "10.10.10.10:1234",
			"BrokerUrl": "tcp://something:port",
			"LogLevel": "DEBUG"
		}`)

		conf, err := config.FromBytes(jsonStr)
		Expect(err).ToNot(HaveOccurred())

		Expect(conf.NotifierUrl).To(Equal("10.10.10.10:1234"))
		Expect(conf.Sensors).To(HaveLen(2))
		Expect(conf.Sensors["2"]).To(Equal("touch"))
		Expect(conf.Sensors["3"]).To(Equal("sound"))
		Expect(conf.LogLevel).To(Equal(logging.DEBUG))
	})

	It("defaults LogLevel to INFO", func() {
		jsonStr := []byte(`{
			"NotifierUrl": "10.10.10.10:1234",
			"BrokerUrl": "tcp://something:port"
		}`)
		conf, err := config.FromBytes(jsonStr)
		Expect(err).ToNot(HaveOccurred())
		Expect(conf.LogLevel).To(Equal(logging.INFO))

	})

	It("stores the last sensor type if two sensors have the same pin", func() {
		jsonStr := []byte(`{
			"Sensors": {
				"2": "touch",
				"2": "sound"
			},
			"NotifierUrl": "10.10.10.10:1234",
			"BrokerUrl": "tcp://something:port"
		}`)

		conf, err := config.FromBytes(jsonStr)
		Expect(err).ToNot(HaveOccurred())
		Expect(conf.Sensors).To(HaveLen(1))
		Expect(conf.Sensors["2"]).To(Equal("sound"))
	})

	It("returns error for missing notifier url", func() {
		jsonStr := []byte(`{
			"Sensors": {
				"2": "touch",
				"3": "sound"
			}
		}`)
		_, err := config.FromBytes(jsonStr)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("Notifier Url required"))
	})

	It("returns error for missing broker url", func() {
		jsonStr := []byte(`{
			"Sensors": {
				"2": "touch",
				"2": "sound"
			},
			"NotifierUrl": "10.10.10.10:1234"
		}`)
		_, err := config.FromBytes(jsonStr)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("Broker Url required"))
	})

	It("returns error for invalid json", func() {
		// missing comma
		jsonStr := []byte(`{
			"Sensors": {
				"2": "touch"
				"3": "sound"
			},
			"NotifierUrl": "10.10.10.10:1234"
		}`)

		_, err := config.FromBytes(jsonStr)
		Expect(err).To(HaveOccurred())
	})

	It("returns error for invalid sensor type", func() {
		jsonStr := []byte(`{
			"Sensors": {
				"2": "touch",
				"3": "bla"
			},
			"NotifierUrl": "10.10.10.10:1234"
		}`)

		_, err := config.FromBytes(jsonStr)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("Invalid sensor type: bla"))
	})

	It("returns error for invalid pin", func() {
		jsonStr := []byte(`{
			"Sensors": {
				"2a": "touch",
				"3": "sound"
			},
			"NotifierUrl": "10.10.10.10:1234"
		}`)

		_, err := config.FromBytes(jsonStr)
		Expect(err).To(HaveOccurred())
	})
})
