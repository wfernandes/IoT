package mqtt

import (
	"time"

	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

type SimpleMQTT interface {
	Connect() error
	Disconnect()
	Publish(string, []byte)
	Subscribe(string, func([]byte))
	IsConnected() bool
}

type MQTTClient struct {
	clientID  string
	brokerURL string
	client    mqtt.ClientInt
}

func New(clientId string, brokerUrl string) *MQTTClient {
	return &MQTTClient{
		clientID:  clientId,
		brokerURL: brokerUrl,
	}
}

func (m *MQTTClient) Connect() error {
	opts := mqtt.NewClientOptions().AddBroker(m.brokerURL).SetClientID(m.clientID)

	m.client = mqtt.NewClient(opts)
	if token := m.client.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MQTTClient) Disconnect() {
	if m.client != nil {
		m.client.Disconnect(500)
	}
}

func (m *MQTTClient) Publish(topic string, value []byte) {
	token := m.client.Publish(topic, 0, false, value)
	token.Wait()
}

func (m *MQTTClient) Subscribe(event string, f func([]byte)) {
	m.client.Subscribe(event, 0, func(client *mqtt.Client, msg mqtt.Message) {
		f(msg.Payload())
	})
}

func (m *MQTTClient) IsConnected() bool {
	return m.client.IsConnected()
}
