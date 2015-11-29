package notification_test

import (
	"github.com/wfernandes/homesec/notification_processor/notification"

	"github.com/cloudfoundry/gosteno"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Notification", func() {

	var (
		inputChan    chan string
		mockNotifier *MockNotifier
		logger       *gosteno.Logger
		ns           *notification.NotificationService
	)

	BeforeEach(func() {
		mockNotifier = &MockNotifier{}
		inputChan = make(chan string)
		logger = gosteno.NewLogger("Notification Test Logger")
	})

	Context("Start", func() {
		AfterEach(func() {
			ns.Stop()
		})

		It("reads from inputCnan and notifies", func() {
			ns = notification.New(mockNotifier, inputChan, logger)

			go ns.Start()

			Expect(mockNotifier.NotifyCallCount()).To(BeZero())
			inputChan <- "test message"
			Eventually(mockNotifier.NotifyCallCount).Should(Equal(1))
			Expect(mockNotifier.lastNotification).To(Equal("test message"))
		})

		PIt("should do something if notify fails", func() {})
	})

	Context("Stop", func() {

		It("sends a shutdown notification message", func() {
			ns = notification.New(mockNotifier, inputChan, logger)
			go ns.Start()
			Expect(mockNotifier.NotifyCallCount()).To(BeZero())
			ns.Stop()
			Eventually(mockNotifier.NotifyCallCount).Should(Equal(1))
			Expect(mockNotifier.lastNotification).To(Equal("Notification Service Shutdown"))
		})
	})
})