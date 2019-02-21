package device_token

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("DeviceTokenMockup", func() {

	DeviceTokenContexts(NewDeviceTokenMockup())
})
