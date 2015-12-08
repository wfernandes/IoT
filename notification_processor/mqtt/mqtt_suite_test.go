package mqtt_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMqtt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mqtt Suite")
}
