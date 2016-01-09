package notification_test

import (
	"github.com/wfernandes/iot/notification_processor/notification"

	"time"

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
			ns = notification.New(mockNotifier, inputChan, time.Second)

			go ns.Start()

			Expect(mockNotifier.NotifyCallCount()).To(BeZero())
			inputChan <- sensorEvent
			Eventually(mockNotifier.NotifyCallCount).Should(Equal(1))
			Expect(mockNotifier.lastNotification).To(Equal("test message"))
		})

		It("doesn't notify again if new notification is received within notification period", func() {

			ns = notification.New(mockNotifier, inputChan, time.Second)
			go ns.Start()

			Expect(mockNotifier.NotifyCallCount()).To(BeZero())
			inputChan <- sensorEvent
			time.Sleep(100 * time.Millisecond)
			inputChan <- sensorEvent
			inputChan <- sensorEvent
			Expect(mockNotifier.NotifyCallCount()).To(Equal(1))

		})

		It("notifies again if new notification is received after notification period", func() {

			ns = notification.New(mockNotifier, inputChan, 100*time.Millisecond)
			go ns.Start()

			Expect(mockNotifier.NotifyCallCount()).To(BeZero())
			inputChan <- sensorEvent
			// This event should not be recorded since it sent instantaneously.
			inputChan <- sensorEvent
			// The next event should be recorded since its past the previous event.
			time.Sleep(200 * time.Millisecond)
			inputChan <- sensorEvent
			Eventually(mockNotifier.NotifyCallCount).Should(Equal(2))

		})

		PIt("should do something if notify fails", func() {})
	})

	Context("Stop", func() {

		It("sends a shutdown notification message", func() {
			ns = notification.New(mockNotifier, inputChan, time.Second)
			go ns.Start()
			Expect(mockNotifier.NotifyCallCount()).To(BeZero())
			ns.Stop()
			Eventually(mockNotifier.NotifyCallCount).Should(Equal(1))
			Expect(mockNotifier.lastNotification).To(Equal("Notification Service Shutdown"))
		})
	})
})
