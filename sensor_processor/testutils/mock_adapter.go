package testutils

import "sync"

type MockAdapter struct {
	name  string
	state int
	lock  sync.Mutex
}

func NewMockAdapter(name string) *MockAdapter {
	return &MockAdapter{
		name: name,
	}
}

// Changing state returned every time this is called so that
// events are triggered which corresponds to events subscribed
// in sensor's work function.
func (m *MockAdapter) DigitalRead(string) (int, error) {
	m.lock.Lock()
	if m.state == 0 {
		m.state = 1
	} else {
		m.state = 0

	}
	m.lock.Unlock()
	return m.state, nil
}

func (m *MockAdapter) Name() string {
	return m.name
}

func (m *MockAdapter) Connect() []error {
	var errs []error
	return errs
}

func (m *MockAdapter) Finalize() []error {
	var errs []error
	return errs
}