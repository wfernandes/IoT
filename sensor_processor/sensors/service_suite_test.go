package sensors_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSensors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sensors Suite")
}
