package role

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("RoleMockup", func() {
	var provider = NewRoleMockup()
	RoleContexts(provider)
})
