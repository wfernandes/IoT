package broker

import (
	"time"

	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

type MQTTBroker struct {
	clientID  string
	brokerURL string
	client    mqtt.ClientInt
}

func NewMQTTBroker(clientId string, brokerUrl string) *MQTTBroker {
	return &MQTTBroker{
		clientID:  clientId,
		brokerURL: brokerUrl,
	}
}

func (m *MQTTBroker) Connect() error {
	opts := mqtt.NewClientOptions().AddBroker(m.brokerURL).SetClientID(m.clientID)

	m.client = mqtt.NewClient(opts)
	if token := m.client.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MQTTBroker) Disconnect() {
	if m.client != nil {
		m.client.Disconnect(500)
	}
}

func (m *MQTTBroker) Publish(topic string, value []byte) {
	token := m.client.Publish(topic, 0, false, value)
	token.Wait()
}

func (m *MQTTBroker) Subscribe(event string, f func([]byte)) {
	m.client.Subscribe(event, 0, func(client *mqtt.Client, msg mqtt.Message) {
		f(msg.Payload())
	})
}

func (m *MQTTBroker) IsConnected() bool {
	return m.client.IsConnected()
}
