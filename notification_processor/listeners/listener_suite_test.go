package listeners_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestListener(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Listeners Suite")
}
