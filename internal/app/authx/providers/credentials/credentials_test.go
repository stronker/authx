/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package credentials

import (
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func CredentialsContexts(provider BasicCredentials) {

	ginkgo.Context("with a register", func() {
		credentials := entities.NewBasicCredentialsData("u1", []byte("p1"), "r1", "o1")
		ginkgo.BeforeEach(func() {
			err := provider.Add(credentials)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("must exist", func() {
			exists, err := provider.Exist(credentials.Username)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(*exists).To(gomega.BeTrue())

			c, err := provider.Get(credentials.Username)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(c).NotTo(gomega.BeNil())

		})

		ginkgo.It("can be edited the password", func() {
			err := provider.Edit(credentials.Username, entities.NewEditBasicCredentialsData().WithPassword([]byte("pNew")))
			gomega.Expect(err).To(gomega.Succeed())
			c, err := provider.Get(credentials.Username)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(c).NotTo(gomega.BeNil())
			gomega.Expect(c.Password).To(gomega.Equal([]byte("pNew")))

		})

		ginkgo.It("can be edited the roleID", func() {
			err := provider.Edit(credentials.Username, entities.NewEditBasicCredentialsData().WithRoleID("rNew"))
			gomega.Expect(err).To(gomega.Succeed())
			c, err := provider.Get(credentials.Username)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(c).NotTo(gomega.BeNil())
			gomega.Expect(c.RoleID).To(gomega.Equal("rNew"))

		})
		ginkgo.It("can be edited without changes", func() {
			err := provider.Edit(credentials.Username, entities.NewEditBasicCredentialsData())
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

		ginkgo.It("should not exist", func() {
			c, err := provider.Exist("u1")
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(*c).To(gomega.BeFalse())
		})

		ginkgo.It("should not work", func() {
			c, err := provider.Get("u1")
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(c).To(gomega.BeNil())
		})

		ginkgo.It("should  add correctly", func() {
			err := provider.Add(entities.NewBasicCredentialsData("u1", []byte("pwd"), "r1", "o1"))
			gomega.Expect(err).To(gomega.Succeed())
		})

		ginkgo.It("should  not work", func() {
			err := provider.Edit("u1", entities.NewEditBasicCredentialsData().WithRoleID("rNew"))
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should not work", func() {
			err := provider.Delete("u1")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

	})
	ginkgo.AfterEach(func() {
		err := provider.Truncate()
		gomega.Expect(err).To(gomega.BeNil())
	})
}
