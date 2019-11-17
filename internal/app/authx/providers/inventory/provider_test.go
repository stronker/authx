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

package inventory

import (
	"github.com/google/uuid"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/stronker/authx/internal/app/authx/entities"
	"time"
)

func CreateTestECJoinToken() *entities.EICJoinToken {
	return entities.NewEICJoinToken(uuid.New().String(), time.Hour)
}

func RunTest(provider Provider) {
	
	ginkgo.BeforeEach(func() {
		provider.Clear()
	})
	
	ginkgo.Context("Edge controllers", func() {
		ginkgo.It("should be able to add a join token", func() {
			toAdd := CreateTestECJoinToken()
			err := provider.AddECJoinToken(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("should be able to retrieve a valid token", func() {
			toAdd := CreateTestECJoinToken()
			err := provider.AddECJoinToken(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			retrieved, err := provider.GetECJoinToken(toAdd.OrganizationID, toAdd.TokenID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
		})
	})
}
