package token

import (
	"fmt"
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"math/rand"
	"time"
)

var _ = ginkgo.Describe("TokenMockup", func() {

	var provider = NewTokenMockup()

	TokenContexts(provider)

	ginkgo.It("with a register", func() {
		for i := 0; i < 10; i++ {
			var p = rand.Intn(10)
			token := entities.NewTokenData(fmt.Sprintf("u%d", i), fmt.Sprintf("t%d", i), []byte("r1"), time.Now().Add(time.Duration(p)*time.Second).Unix())
			err := provider.Add(token)
			gomega.Expect(err).To(gomega.Succeed())
		}

		err := provider.DeleteExpiredTokens()
		gomega.Expect(err).To(gomega.Succeed())

	})
})