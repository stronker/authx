package manager

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/authx/internal/app/authx/providers/device"
	"github.com/nalej/authx/internal/app/authx/providers/device_token"
	"github.com/nalej/authx/pkg/token"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"time"
)

var _ = ginkgo.Describe("Device Token tests", func() {
	var devProvider = device.NewMockupDeviceCredentialsProvider()
	var devTokenManager = NewJWTDeviceToken(devProvider, device_token.NewDeviceTokenMockup() )


	expirationPeriod, _ := time.ParseDuration("10m")
	secret := "myLittleSecret12345"

	deviceClaim := token.NewDeviceClaim(uuid.New().String(), uuid.New().String(), uuid.New().String(), expirationPeriod)

	ginkgo.Context("Device Token tests", func() {
		var group entities.DeviceGroupCredentials
		ginkgo.BeforeSuite(func(){
			group = entities.DeviceGroupCredentials{
				OrganizationID: deviceClaim.OrganizationID,
				DeviceGroupID: deviceClaim.DeviceGroupID,
				DeviceGroupApiKey: uuid.New().String(),
				Enabled: true,
				DefaultDeviceConnectivity: true,
				Secret: uuid.New().String(),
			}
			err := devProvider.AddDeviceGroupCredentials(&group)
			gomega.Expect(err).To(gomega.Succeed())
		})


		ginkgo.It("Can generate a token", func() {
			gT, err := devTokenManager.Generate(deviceClaim, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT).NotTo(gomega.BeNil())

		})
		ginkgo.It("can add a device token twice", func() {
			gT, err := devTokenManager.Generate(deviceClaim, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT).NotTo(gomega.BeNil())

			gT2, err := devTokenManager.Generate(deviceClaim, expirationPeriod, secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT2).NotTo(gomega.BeNil())

		})
		ginkgo.It("can refresh a device token", func() {



			gT, err := devTokenManager.Generate(deviceClaim, expirationPeriod, group.Secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT).NotTo(gomega.BeNil())


			tk, jwtErr := jwt.ParseWithClaims(gT.Token, &token.DeviceClaim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(group.Secret), nil
			})
			gomega.Expect(jwtErr).To(gomega.Succeed())

			cl, ok := tk.Claims.(*token.DeviceClaim)
			gomega.Expect(ok).To(gomega.BeTrue())
			gomega.Expect(cl).NotTo(gomega.BeNil())


			gTNew, err := devTokenManager.Refresh(gT.Token, gT.RefreshToken, expirationPeriod, group.Secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gTNew).NotTo(gomega.BeNil())
			gomega.Expect(gTNew).NotTo(gomega.Equal(gT))
		})
		ginkgo.It("must be able to reject an expired refresh token", func() {

			d, _ := time.ParseDuration("-1s")

			gT, err := devTokenManager.Generate(deviceClaim, d, group.Secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT).NotTo(gomega.BeNil())

			parser := jwt.Parser{SkipClaimsValidation: true}
			tk, jwtErr := parser.ParseWithClaims(gT.Token, &token.Claim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(group.Secret), nil
			})
			gomega.Expect(jwtErr).To(gomega.Succeed())

			cl, ok := tk.Claims.(*token.Claim)
			gomega.Expect(ok).To(gomega.BeTrue())
			gomega.Expect(cl).NotTo(gomega.BeNil())

			gTNew, err := devTokenManager.Refresh(gT.Token,gT.RefreshToken, expirationPeriod, group.Secret)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(gTNew).To(gomega.BeNil())

		})
		ginkgo.It("must be able to reject the refresh token is incorrect", func() {

			gT, err := devTokenManager.Generate(deviceClaim, expirationPeriod, group.Secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT).NotTo(gomega.BeNil())

			parser := jwt.Parser{SkipClaimsValidation: true}
			tk, jwtErr := parser.ParseWithClaims(gT.Token, &token.DeviceClaim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(group.Secret), nil
			})
			gomega.Expect(jwtErr).To(gomega.Succeed())

			cl, ok := tk.Claims.(*token.DeviceClaim)
			gomega.Expect(ok).To(gomega.BeTrue())
			gomega.Expect(cl).NotTo(gomega.BeNil())

			gTNew, err := devTokenManager.Refresh(gT.Token, gT.RefreshToken+"wrong", expirationPeriod, group.Secret)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(gTNew).To(gomega.BeNil())

		})

		ginkgo.It("must be able to reject the token is incorrect", func() {

			gT, err := devTokenManager.Generate(deviceClaim, expirationPeriod, group.Secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT).NotTo(gomega.BeNil())

			parser := jwt.Parser{SkipClaimsValidation: true}
			tk, jwtErr := parser.ParseWithClaims(gT.Token, &token.DeviceClaim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(group.Secret), nil
			})
			gomega.Expect(jwtErr).To(gomega.Succeed())

			cl, ok := tk.Claims.(*token.DeviceClaim)
			gomega.Expect(ok).To(gomega.BeTrue())
			gomega.Expect(cl).NotTo(gomega.BeNil())

			gTNew, err := devTokenManager.Refresh(gT.Token+"wrong", gT.RefreshToken, expirationPeriod, group.Secret)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(gTNew).To(gomega.BeNil())

		})

		ginkgo.It("can't use two times the same refresh token", func() {

			gT, err := devTokenManager.Generate(deviceClaim, expirationPeriod, group.Secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gT).NotTo(gomega.BeNil())

			tk, jwtErr := jwt.ParseWithClaims(gT.Token, &token.DeviceClaim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(group.Secret), nil
			})
			gomega.Expect(jwtErr).To(gomega.Succeed())

			cl, ok := tk.Claims.(*token.DeviceClaim)
			gomega.Expect(ok).To(gomega.BeTrue())
			gomega.Expect(cl).NotTo(gomega.BeNil())

			gTNew, err := devTokenManager.Refresh(gT.Token, gT.RefreshToken, expirationPeriod, group.Secret)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(gTNew).NotTo(gomega.BeNil())
			gomega.Expect(gTNew).NotTo(gomega.Equal(gT))

			gTWrong, err := devTokenManager.Refresh(gT.Token, gT.RefreshToken, expirationPeriod, group.Secret)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(gTWrong).To(gomega.BeNil())

		})

	})
	ginkgo.AfterEach(func() {
		err := devTokenManager.Clean()
		gomega.Expect(err).To(gomega.Succeed())
	})
})

