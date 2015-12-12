package sensors_test

type mockBroker struct {
	ConnectCalled chan bool
	ConnectOutput struct {
		ret0 chan error
	}
	PublishCalled chan bool
	PublishInput  struct {
		arg0 chan string
		arg1 chan []byte
	}
	SubscribeCalled chan bool
	SubscribeInput  struct {
		arg0 chan string
		arg1 chan func([]byte)
	}
	IsConnectedCalled chan bool
	IsConnectedOutput struct {
		ret0 chan bool
	}
	DisconnectCalled chan bool
}

func newMockBroker() *mockBroker {
	m := &mockBroker{}
	m.ConnectCalled = make(chan bool, 100)
	m.ConnectOutput.ret0 = make(chan error, 100)
	m.PublishCalled = make(chan bool, 100)
	m.PublishInput.arg0 = make(chan string, 100)
	m.PublishInput.arg1 = make(chan []byte, 100)
	m.SubscribeCalled = make(chan bool, 100)
	m.SubscribeInput.arg0 = make(chan string, 100)
	m.SubscribeInput.arg1 = make(chan func([]byte), 100)
	m.IsConnectedCalled = make(chan bool, 100)
	m.IsConnectedOutput.ret0 = make(chan bool, 100)
	m.DisconnectCalled = make(chan bool, 100)
	return m
}
func (m *mockBroker) Connect() error {
	m.ConnectCalled <- true
	return <-m.ConnectOutput.ret0
}
func (m *mockBroker) Publish(arg0 string, arg1 []byte) {
	m.PublishCalled <- true
	m.PublishInput.arg0 <- arg0
	m.PublishInput.arg1 <- arg1
}
func (m *mockBroker) Subscribe(arg0 string, arg1 func([]byte)) {
	m.SubscribeCalled <- true
	m.SubscribeInput.arg0 <- arg0
	m.SubscribeInput.arg1 <- arg1
}
func (m *mockBroker) IsConnected() bool {
	m.IsConnectedCalled <- true
	return <-m.IsConnectedOutput.ret0
}
func (m *mockBroker) Disconnect() {
	m.DisconnectCalled <- true
}
