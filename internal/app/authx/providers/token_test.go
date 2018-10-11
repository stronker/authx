/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package providers

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"time"
)

var _ = ginkgo.Describe("TokenMockup", func() {
	var provider = NewTokenMockup()
	TokenContexts(provider)
})

func TokenContexts(provider Token) {

	ginkgo.Context("with a register", func() {
		token := NewTokenData("u1", "t1", []byte("r1"), time.Now().Unix())
		ginkgo.BeforeEach(func() {
			err := provider.Add(token)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("must exist", func() {
			t, err := provider.Get(token.Username, token.TokenID)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(t).NotTo(gomega.BeNil())

		})

		ginkgo.It("can delete the token", func() {
			err := provider.Delete(token.Username, token.TokenID)
			gomega.Expect(err).To(gomega.Succeed())
			t, err := provider.Get(token.Username, token.TokenID)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(t).To(gomega.BeNil())

		})

	})
	ginkgo.Context("empty data store", func() {
		ginkgo.It("must fail with get", func() {
			c, err := provider.Get("u1", "t1")
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(c).To(gomega.BeNil())
		})

		ginkgo.It("must add correctly", func() {
			err := provider.Add(NewTokenData("u1", "t1", []byte("rt1"), 11111))
			gomega.Expect(err).To(gomega.Succeed())
		})

		ginkgo.It("delete doesn't work", func() {
			err := provider.Delete("u1", "t1")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

	})
	ginkgo.AfterEach(func() {
		err := provider.Truncate()
		gomega.Expect(err).To(gomega.BeNil())
	})
}
