package device

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestDeviceCredentialsPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Device Credentials providers package suite")
}
