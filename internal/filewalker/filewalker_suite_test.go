package filewalker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFilewalker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filewalker Suite")
}
