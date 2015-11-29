package notification_test

import "sync"

type MockNotifier struct {
	notifyCallCount  int
	lastNotification string
	returnError      error
	lock             sync.Mutex
}

func (m *MockNotifier) Notify(body string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.notifyCallCount++
	m.lastNotification = body

	if m.returnError != nil {
		return m.returnError
	} else {
		return nil
	}
}

func (m *MockNotifier) NotifyCallCount() int {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.notifyCallCount
}
