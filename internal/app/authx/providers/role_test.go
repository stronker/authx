/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package providers

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("RoleMockup", func() {
	var provider = NewRoleMockup()
	RoleContexts(provider)
})

func RoleContexts(provider Role) {

	ginkgo.Context("with a register", func() {
		role := NewRoleData("o1", "r1", "n1", [] string{"p1", "p2"})
		ginkgo.BeforeEach(func() {
			err := provider.Add(role)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("must exist", func() {
			r, err := provider.Get(role.OrganizationId, role.RoleId)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(r).NotTo(gomega.BeNil())

		})

		ginkgo.It("can be edited the name", func() {
			err := provider.Edit(role.OrganizationId, role.RoleId, NewEditRoleData().WithName("nNew"))
			gomega.Expect(err).To(gomega.Succeed())
			r, err := provider.Get(role.OrganizationId, role.RoleId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(r).NotTo(gomega.BeNil())
			gomega.Expect(r.Name).To(gomega.Equal("nNew"))

		})

		ginkgo.It("can be edited the roleID", func() {
			err := provider.Edit(role.OrganizationId, role.RoleId, NewEditRoleData().WithPrimitives([]string{"pNew"}))
			gomega.Expect(err).To(gomega.Succeed())
			r, err := provider.Get(role.OrganizationId, role.RoleId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(r).NotTo(gomega.BeNil())
			gomega.Expect(r.Primitives).To(gomega.Equal([]string{"pNew"}))

		})
		ginkgo.It("can be edited without changes", func() {
			err := provider.Edit(role.OrganizationId, role.RoleId, NewEditRoleData())
			gomega.Expect(err).To(gomega.Succeed())
			r, err := provider.Get(role.OrganizationId, role.RoleId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(r).NotTo(gomega.BeNil())
			gomega.Expect(r.Name).To(gomega.Equal(role.Name))
			gomega.Expect(r.Primitives).To(gomega.Equal(role.Primitives))

		})

		ginkgo.It("can delete the token", func() {
			err := provider.Delete(role.OrganizationId, role.RoleId)
			gomega.Expect(err).To(gomega.Succeed())
			r, err := provider.Get(role.OrganizationId, role.RoleId)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(r).To(gomega.BeNil())
		})

		ginkgo.AfterEach(func() {
			err := provider.Truncate()
			gomega.Expect(err).To(gomega.BeNil())
		})
	})
	ginkgo.Context("empty data store", func() {
		ginkgo.It("must add correctly", func() {
			c, err := provider.Get("o1", "r1")
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(c).To(gomega.BeNil())
		})

		ginkgo.It("edit doesn't work", func() {
			err := provider.Edit("o1", "r1", NewEditRoleData().WithName("nNew"))
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("delete doesn't work", func() {
			err := provider.Delete("o1", "r1")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

	})
}
