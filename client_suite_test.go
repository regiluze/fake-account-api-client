package accountclient_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestKpiManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "account api client Suite")
}
