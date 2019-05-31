package credentials

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("BasicCredentialsMockup", func() {
	var provider = NewBasicCredentialMockup()
	CredentialsContexts(provider)
})
