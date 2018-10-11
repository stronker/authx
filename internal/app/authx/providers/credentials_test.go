/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package providers

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("BasicCredentialsMockup", func() {
	var provider = NewBasicCredentialMockup()
	CredentialsContexts(provider)
})

func CredentialsContexts(provider BasicCredentials) {

	ginkgo.Context("with a register", func() {
		credentials := NewBasicCredentialsData("u1", []byte("p1"), "r1", "o1")
		ginkgo.BeforeEach(func() {
			err := provider.Add(credentials)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("must exist", func() {
			c, err := provider.Get(credentials.Username)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(c).NotTo(gomega.BeNil())

		})

		ginkgo.It("can be edited the password", func() {
			err := provider.Edit(credentials.Username, NewEditBasicCredentialsData().WithPassword([]byte("pNew")))
			gomega.Expect(err).To(gomega.Succeed())
			c, err := provider.Get(credentials.Username)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(c).NotTo(gomega.BeNil())
			gomega.Expect(c.Password).To(gomega.Equal([]byte("pNew")))

		})

		ginkgo.It("can be edited the roleID", func() {
			err := provider.Edit(credentials.Username, NewEditBasicCredentialsData().WithRoleID("rNew"))
			gomega.Expect(err).To(gomega.Succeed())
			c, err := provider.Get(credentials.Username)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(c).NotTo(gomega.BeNil())
			gomega.Expect(c.RoleID).To(gomega.Equal("rNew"))

		})
		ginkgo.It("can be edited without changes", func() {
			err := provider.Edit(credentials.Username, NewEditBasicCredentialsData())
			gomega.Expect(err).To(gomega.Succeed())
			c, err := provider.Get(credentials.Username)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(c).NotTo(gomega.BeNil())
			gomega.Expect(c.Password).To(gomega.Equal(credentials.Password))

		})

		ginkgo.It("can delete the credentials", func() {
			err := provider.Delete(credentials.Username)
			gomega.Expect(err).To(gomega.Succeed())
			c, err := provider.Get(credentials.Username)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(c).To(gomega.BeNil())
		})


	})
	ginkgo.Context("empty data store", func() {
		ginkgo.It("get doesn't work", func() {
			c, err := provider.Get("u1")
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(c).To(gomega.BeNil())
		})

		ginkgo.It("must add correctly", func() {
			err := provider.Add(NewBasicCredentialsData("u1",[] byte("pwd"),"r1","o1"))
			gomega.Expect(err).To(gomega.Succeed())
		})

		ginkgo.It("edit doesn't work", func() {
			err := provider.Edit("u1", NewEditBasicCredentialsData().WithRoleID("rNew"))
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("delete doesn't work", func() {
			err := provider.Delete("u1")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

	})
	ginkgo.AfterEach(func() {
		err := provider.Truncate()
		gomega.Expect(err).To(gomega.BeNil())
	})
}
