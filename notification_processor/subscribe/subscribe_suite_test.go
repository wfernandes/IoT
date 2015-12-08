package subscribe_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSubscribe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Subscribe Suite")
}
