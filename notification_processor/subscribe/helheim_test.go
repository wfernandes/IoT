package subscribe_test

type mockBroker struct {
	ConnectCalled chan bool
	ConnectOutput struct {
		ret0 chan error
	}
	SubscribeCalled chan bool
	SubscribeInput  struct {
		arg0 chan string
		arg1 chan func([]byte)
	}
	DisconnectCalled chan bool
}

func newMockBroker() *mockBroker {
	m := &mockBroker{}
	m.ConnectCalled = make(chan bool, 100)
	m.ConnectOutput.ret0 = make(chan error, 100)
	m.SubscribeCalled = make(chan bool, 100)
	m.SubscribeInput.arg0 = make(chan string, 100)
	m.SubscribeInput.arg1 = make(chan func([]byte), 100)
	m.DisconnectCalled = make(chan bool, 100)
	return m
}
func (m *mockBroker) Connect() error {
	m.ConnectCalled <- true
	return <-m.ConnectOutput.ret0
}
func (m *mockBroker) Subscribe(arg0 string, arg1 func([]byte)) {
	m.SubscribeCalled <- true
	m.SubscribeInput.arg0 <- arg0
	m.SubscribeInput.arg1 <- arg1
}
func (m *mockBroker) Disconnect() {
	m.DisconnectCalled <- true
}
