/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package manager

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/nalej/authx/pkg/token"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"time"
)

var _ = ginkgo.Describe("BCryptPassword", func() {
	var manager = NewJWTTokenMockup()
	TokenContexts(manager)
})

func TokenContexts(manager Token) {
	ginkgo.Context("with a basic parameters", func() {
		claim := token.NewPersonalClaim("u1", "r1", []string{"p1", "p2"},"o1")
		expirationPeriod, _ := time.ParseDuration("10m")
		secret := "myLittleSecret112131"
		ginkgo.It("can generate a token", func() {
			gT, err := manager.Generate(claim, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT).NotTo(gomega.BeNil())

		})

		ginkgo.It("can generated two tokens for the same user", func() {
			gT, err := manager.Generate(claim, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT).NotTo(gomega.BeNil())

			gTNew, err := manager.Generate(claim, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gTNew).NotTo(gomega.BeNil())
			gomega.Expect(gTNew).NotTo(gomega.Equal(gT))
		})

		ginkgo.It("can refresh a token", func() {
			gT, err := manager.Generate(claim, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT).NotTo(gomega.BeNil())

			tk, jwtErr := jwt.ParseWithClaims(gT.Token, &token.Claim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			gomega.Expect(jwtErr).To(gomega.Succeed())

			cl, ok := tk.Claims.(*token.Claim)
			gomega.Expect(ok).To(gomega.BeTrue())
			gomega.Expect(cl).NotTo(gomega.BeNil())

			gTNew, err := manager.Refresh(claim, cl.Id, gT.RefreshToken, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gTNew).NotTo(gomega.BeNil())
			gomega.Expect(gTNew).NotTo(gomega.Equal(gT))
		})

		ginkgo.It("must be able to reject an expired refresh token", func() {

			d, _ := time.ParseDuration("-1s")

			gT, err := manager.Generate(claim, d, secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT).NotTo(gomega.BeNil())

			parser := jwt.Parser{SkipClaimsValidation: true}
			tk, jwtErr := parser.ParseWithClaims(gT.Token, &token.Claim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			gomega.Expect(jwtErr).To(gomega.Succeed())

			cl, ok := tk.Claims.(*token.Claim)
			gomega.Expect(ok).To(gomega.BeTrue())
			gomega.Expect(cl).NotTo(gomega.BeNil())

			gTNew, err := manager.Refresh(claim, cl.Id, gT.RefreshToken, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(gTNew).To(gomega.BeNil())

		})

		ginkgo.It("must be able to reject the refresh token is incorrect", func() {

			gT, err := manager.Generate(claim, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT).NotTo(gomega.BeNil())

			parser := jwt.Parser{SkipClaimsValidation: true}
			tk, jwtErr := parser.ParseWithClaims(gT.Token, &token.Claim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			gomega.Expect(jwtErr).To(gomega.Succeed())

			cl, ok := tk.Claims.(*token.Claim)
			gomega.Expect(ok).To(gomega.BeTrue())
			gomega.Expect(cl).NotTo(gomega.BeNil())

			gTNew, err := manager.Refresh(claim, cl.Id, gT.RefreshToken+"wrong", expirationPeriod, secret)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(gTNew).To(gomega.BeNil())

		})

		ginkgo.It("must be able to reject the tokenID is incorrect", func() {

			gT, err := manager.Generate(claim, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT).NotTo(gomega.BeNil())

			parser := jwt.Parser{SkipClaimsValidation: true}
			tk, jwtErr := parser.ParseWithClaims(gT.Token, &token.Claim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			gomega.Expect(jwtErr).To(gomega.Succeed())

			cl, ok := tk.Claims.(*token.Claim)
			gomega.Expect(ok).To(gomega.BeTrue())
			gomega.Expect(cl).NotTo(gomega.BeNil())

			gTNew, err := manager.Refresh(claim, cl.Id+"wrong", gT.RefreshToken, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(gTNew).To(gomega.BeNil())

		})

		ginkgo.It("can't use two times the same refresh token", func() {

			gT, err := manager.Generate(claim, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT).NotTo(gomega.BeNil())

			tk, jwtErr := jwt.ParseWithClaims(gT.Token, &token.Claim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			gomega.Expect(jwtErr).To(gomega.Succeed())

			cl, ok := tk.Claims.(*token.Claim)
			gomega.Expect(ok).To(gomega.BeTrue())
			gomega.Expect(cl).NotTo(gomega.BeNil())

			gTNew, err := manager.Refresh(claim, cl.Id, gT.RefreshToken, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gTNew).NotTo(gomega.BeNil())
			gomega.Expect(gTNew).NotTo(gomega.Equal(gT))

			gTWrong, err := manager.Refresh(claim, cl.Id, gT.RefreshToken, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(gTWrong).To(gomega.BeNil())

		})
		ginkgo.AfterEach(func() {
			err := manager.Clean()
			gomega.Expect(err).To(gomega.Succeed())
		})

	})

}
