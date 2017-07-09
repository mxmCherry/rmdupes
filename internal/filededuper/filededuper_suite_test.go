package filededuper_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFilededuper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filededuper Suite")
}
