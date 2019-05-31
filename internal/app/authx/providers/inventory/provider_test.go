/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package inventory

import (
	"github.com/google/uuid"
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
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
