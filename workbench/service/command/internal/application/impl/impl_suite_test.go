package impl

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestImpl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Implementation Suite")
}
