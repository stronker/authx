/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package role

import (
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RoleContexts(provider Role) {

	ginkgo.Context("with a register", func() {
		role := entities.NewRoleData("o1", "r1", "n1", false, []string{"p1", "p2"})
		ginkgo.BeforeEach(func() {
			err := provider.Add(role)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("must exist", func() {
			exists, err := provider.Exist("o1", "r1")
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(*exists).To(gomega.BeTrue())

			r, err := provider.Get(role.OrganizationID, role.RoleID)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(r).NotTo(gomega.BeNil())

		})

		ginkgo.It("can be edited the name", func() {
			err := provider.Edit(role.OrganizationID, role.RoleID, entities.NewEditRoleData().WithName("nNew"))
			gomega.Expect(err).To(gomega.Succeed())
			r, err := provider.Get(role.OrganizationID, role.RoleID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(r).NotTo(gomega.BeNil())
			gomega.Expect(r.Name).To(gomega.Equal("nNew"))

		})

		ginkgo.It("can be edited the roleID", func() {
			err := provider.Edit(role.OrganizationID, role.RoleID, entities.NewEditRoleData().WithPrimitives([]string{"pNew"}))
			gomega.Expect(err).To(gomega.Succeed())
			r, err := provider.Get(role.OrganizationID, role.RoleID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(r).NotTo(gomega.BeNil())
			gomega.Expect(r.Primitives).To(gomega.Equal([]string{"pNew"}))

		})
		ginkgo.It("can be edited without changes", func() {
			err := provider.Edit(role.OrganizationID, role.RoleID, entities.NewEditRoleData())
			gomega.Expect(err).To(gomega.Succeed())
			r, err := provider.Get(role.OrganizationID, role.RoleID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(r).NotTo(gomega.BeNil())
			gomega.Expect(r.Name).To(gomega.Equal(role.Name))
			gomega.Expect(r.Primitives).To(gomega.Equal(role.Primitives))

		})

		ginkgo.It("can delete the token", func() {
			err := provider.Delete(role.OrganizationID, role.RoleID)
			gomega.Expect(err).To(gomega.Succeed())
			r, err := provider.Get(role.OrganizationID, role.RoleID)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(r).To(gomega.BeNil())
		})

		ginkgo.AfterEach(func() {
			err := provider.Truncate()
			gomega.Expect(err).To(gomega.BeNil())
		})
	})
	ginkgo.Context("empty data store", func() {

		ginkgo.It("doesn't exist", func() {
			c, err := provider.Exist("o1", "r1")
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(*c).To(gomega.BeFalse())
		})

		ginkgo.It("must add correctly", func() {
			c, err := provider.Get("o1", "r1")
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(c).To(gomega.BeNil())
		})

		ginkgo.It("edit doesn't work", func() {
			err := provider.Edit("o1", "r1", entities.NewEditRoleData().WithName("nNew"))
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("delete doesn't work", func() {
			err := provider.Delete("o1", "r1")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

	})
}
