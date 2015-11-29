package config_test

import (
	"github.com/wfernandes/homesec/notification_processor/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	Context("defaults", func() {
		It("reads from the default configuration", func() {
			config, err := config.DefaultConfig()
			Expect(err).ToNot(HaveOccurred())
			Expect(config).ToNot(BeNil())
		})
	})

	Context("parsing the json", func() {

		It("returns error for missing twilio account number", func() {

			var jsonStr []byte = []byte(`{
	    		"TwilioAuthToken": "some_auth_token"
			}`)

			_, err := config.FromBytes(jsonStr)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Twilio Account SID is required"))
		})
		It("returns error for missing twilio auth token", func() {

			var jsonStr []byte = []byte(`{
	    		"TwilioAccountSid": "some_sid"
			}`)

			_, err := config.FromBytes(jsonStr)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Twilio Auth Token is required"))
		})

		It("returns error for missing from phone", func() {
			var jsonStr []byte = []byte(`{
				"TwilioAccountSid": "some_sid",
				"TwilioAuthToken": "some_token"
			}`)
			_, err := config.FromBytes(jsonStr)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Twilio From Phone is required"))
		})

		It("returns error for incorrect json", func() {
			var jsonStr []byte = []byte(`{
				"AccountSid": "some_sid"
				"AuthToken": "some_auth_token"
	  		}`)

			_, err := config.FromBytes(jsonStr)
			Expect(err).To(HaveOccurred())
		})

		It("returns config for valid config", func() {
			var jsonStr []byte = []byte(`{
				"TwilioAccountSid": "some_sid",
				"TwilioAuthToken": "some_auth_token",
				"TwilioFromPhone": "some_phone_number",
				"To": "some_phone_number",
				"Port": 1234
		  	}`)

			config, err := config.FromBytes(jsonStr)
			Expect(err).ToNot(HaveOccurred())
			Expect(config).ToNot(BeNil())
			Expect(config.TwilioAccountSid).To(Equal("some_sid"))
			Expect(config.TwilioAuthToken).To(Equal("some_auth_token"))
			Expect(config.TwilioFromPhone).To(Equal("some_phone_number"))
			Expect(config.To).To(Equal("some_phone_number"))
			Expect(config.Port).To(BeEquivalentTo(1234))
		})
	})

	Context("reading the file", func() {

		It("returns error for invalid config file path", func() {
			_, err := config.Configuration("idonotexist.json")
			Expect(err).To(HaveOccurred())
		})

		It("returns config for valid config file", func() {
			config, err := config.Configuration("fixtures/config.json")
			Expect(err).ToNot(HaveOccurred())
			Expect(config.TwilioAccountSid).ToNot(BeEmpty())
			Expect(config.TwilioAuthToken).ToNot(BeEmpty())
		})
	})
})
