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

package token

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/stronker/authx/internal/app/authx/entities"
	"time"
)

func TokenContexts(provider Token) {
	
	ginkgo.Context("with a register", func() {
		token := entities.NewTokenData("u1", "t1", []byte("r1"), time.Now().Unix())
		ginkgo.BeforeEach(func() {
			err := provider.Add(token)
			gomega.Expect(err).To(gomega.BeNil())
		})
		
		ginkgo.It("must exist", func() {
			exist, err := provider.Exist(token.Username, token.TokenID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(*exist).To(gomega.BeTrue())
			
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
		ginkgo.It("should be able to update the token", func() {
			token.ExpirationDate = time.Now().Add(time.Second * 2).Unix()
			token.RefreshToken = []byte("r2")
			err := provider.Update(token)
			gomega.Expect(err).To(gomega.Succeed())
		})
		
	})
	ginkgo.Context("empty data store", func() {
		
		ginkgo.It("doesn't exist", func() {
			c, err := provider.Exist("u1", "t1")
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(*c).To(gomega.BeFalse())
		})
		
		ginkgo.It("must fail with get", func() {
			c, err := provider.Get("u1", "t1")
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(c).To(gomega.BeNil())
		})
		
		ginkgo.It("must add correctly", func() {
			err := provider.Add(entities.NewTokenData("u1", "t1", []byte("rt1"), 11111))
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
