/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package manager

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("BCryptPassword", func() {
	var manager = NewBCryptPassword()
	PasswordContexts(manager)
})

func PasswordContexts(manager Password) {
	ginkgo.Context("with a password", func() {
		pass := "123assSSda132131"
		ginkgo.It("can hash it", func() {
			hashed, err := manager.GenerateHashedPassword(pass)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(hashed).NotTo(gomega.BeNil())

		})

		ginkgo.It("must be different with another hashed password", func() {
			hashed, err := manager.GenerateHashedPassword(pass)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(hashed).NotTo(gomega.BeNil())

			otherHashed, err := manager.GenerateHashedPassword(pass)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(otherHashed).NotTo(gomega.BeNil())
			gomega.Expect(otherHashed).NotTo(gomega.Equal(hashed))
		})

		ginkgo.It("must be able to validate a correct password", func() {
			hashed, err := manager.GenerateHashedPassword(pass)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(hashed).NotTo(gomega.BeNil())

			err = manager.CompareHashAndPassword(hashed, pass)
			gomega.Expect(err).To(gomega.Succeed())
		})

		ginkgo.It("must be able to reject an incorrect password", func() {
			hashed, err := manager.GenerateHashedPassword(pass)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(hashed).NotTo(gomega.BeNil())

			err = manager.CompareHashAndPassword(hashed, pass+"wrong")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

	})

}
