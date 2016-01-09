package notification_test

import (
	"github.com/wfernandes/iot/notification_processor/notification"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wfernandes/iot/event"
)

var _ = Describe("Notification", func() {

	var (
		inputChan    chan *event.Event
		mockNotifier *MockNotifier
		ns           *notification.NotificationService
		sensorEvent  *event.Event
	)

	BeforeEach(func() {
		mockNotifier = &MockNotifier{}
		inputChan = make(chan *event.Event)
		sensorEvent = &event.Event{
			Name: "touchsensor1",
			Data: "test message",
		}
	})

	Context("Start", func() {
		AfterEach(func() {
			ns.Stop()
		})

		It("reads from inputChan and notifies", func() {
			ns = notification.New(mockNotifier, inputChan)

			go ns.Start()

			Expect(mockNotifier.NotifyCallCount()).To(BeZero())
			inputChan <- sensorEvent
			Eventually(mockNotifier.NotifyCallCount).Should(Equal(1))
			Expect(mockNotifier.lastNotification).To(Equal("test message"))
		})

		PIt("should do something if notify fails", func() {})
	})

	Context("Stop", func() {

		It("sends a shutdown notification message", func() {
			ns = notification.New(mockNotifier, inputChan)
			go ns.Start()
			Expect(mockNotifier.NotifyCallCount()).To(BeZero())
			ns.Stop()
			Eventually(mockNotifier.NotifyCallCount).Should(Equal(1))
			Expect(mockNotifier.lastNotification).To(Equal("Notification Service Shutdown"))
		})
	})
})
