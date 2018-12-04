package token

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("TokenMockup", func() {
	var provider = NewTokenMockup()
	TokenContexts(provider)
})