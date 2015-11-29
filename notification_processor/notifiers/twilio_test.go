package notifiers_test

import (
	"github.com/wfernandes/homesec/notification_processor/notifiers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const TEST_SID = "AC7e86f1c52bebc2d081f66e4e380334d2"
const TEST_TOKEN = "c15098c9a04d95004eb671a89e65a3e7"

// https://www.twilio.com/docs/api/rest/test-credentials
const VALID_FROM = "+15005550006"

var _ = Describe("Twilio", func() {

	It("returns error when from number is not specified", func() {
		to := "+15005550006"
		service := notifiers.NewTwilio(TEST_SID, TEST_TOKEN, "", to)

		err := service.Notify("this is a test")
		Expect(err).To(HaveOccurred())
		// Code for missing from number
		Expect(err.Error()).To(ContainSubstring("21603"))
	})

	It("notifies", func() {
		to := VALID_FROM
		service := notifiers.NewTwilio(TEST_SID, TEST_TOKEN, VALID_FROM, to)

		err := service.Notify("this is a test")
		Expect(err).ToNot(HaveOccurred())
	})

	It("", func() {

	})
})
