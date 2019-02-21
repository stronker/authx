package device_token

import (
	"github.com/google/uuid"
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"time"
)

func DeviceTokenContexts(provider Provider) {

	ginkgo.Context("adding device token...", func(){
		ginkgo.It("should be able to add", func(){

			deviceToken := entities.DeviceTokenData{
				DeviceId: uuid.New().String(),
				TokenID: uuid.New().String(),
				RefreshToken: uuid.New().String(),
				ExpirationDate:time.Now().Unix(),
				OrganizationId: uuid.New().String(),
				DeviceGroupId: uuid.New().String(),

			}
			err := provider.Add(&deviceToken)
			gomega.Expect(err).To(gomega.Succeed())
		})
	})
	ginkgo.Context("deleting device token...", func(){
		ginkgo.It("should be able to delete", func(){

			deviceToken := entities.DeviceTokenData{
				DeviceId: uuid.New().String(),
				TokenID: uuid.New().String(),
				RefreshToken: uuid.New().String(),
				ExpirationDate:time.Now().Unix(),
				OrganizationId: uuid.New().String(),
				DeviceGroupId: uuid.New().String(),

			}
			err := provider.Add(&deviceToken)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.Delete(deviceToken.DeviceId, deviceToken.TokenID)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("should not be able to delete", func(){

			err := provider.Delete(uuid.New().String(), uuid.New().String())
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})
	ginkgo.Context("getting device token...", func(){
		ginkgo.It("should be able to get a device token", func(){

			deviceToken := entities.DeviceTokenData{
				DeviceId: uuid.New().String(),
				TokenID: uuid.New().String(),
				RefreshToken: uuid.New().String(),
				ExpirationDate:time.Now().Unix(),
				OrganizationId: uuid.New().String(),
				DeviceGroupId: uuid.New().String(),

			}
			err := provider.Add(&deviceToken)
			gomega.Expect(err).To(gomega.Succeed())

			retrieved, err := provider.Get(deviceToken.DeviceId, deviceToken.TokenID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).NotTo(gomega.BeNil())
			gomega.Expect(retrieved.OrganizationId).Should(gomega.Equal(retrieved.OrganizationId))
			gomega.Expect(retrieved.DeviceGroupId).Should(gomega.Equal(retrieved.DeviceGroupId))
		})
		ginkgo.It("should not be able to get", func(){

			_, err := provider.Get(uuid.New().String(), uuid.New().String())
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})
	ginkgo.Context("finding device token...", func(){
		ginkgo.It("should be able to find a device token", func(){

			deviceToken := entities.DeviceTokenData{
				DeviceId: uuid.New().String(),
				TokenID: uuid.New().String(),
				RefreshToken: uuid.New().String(),
				ExpirationDate:time.Now().Unix(),
				OrganizationId: uuid.New().String(),
				DeviceGroupId: uuid.New().String(),

			}
			err := provider.Add(&deviceToken)
			gomega.Expect(err).To(gomega.Succeed())

			exists, err := provider.Exist(deviceToken.DeviceId, deviceToken.TokenID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(*exists).To(gomega.BeTrue())
		})
		ginkgo.It("should not be able to find a device token", func(){

			exists, err := provider.Exist(uuid.New().String(), uuid.New().String())
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(*exists).NotTo(gomega.BeTrue())
		})
	})
	ginkgo.Context("updating device token...", func(){
		ginkgo.It("should be able to update a device token", func(){

			deviceToken := entities.DeviceTokenData{
				DeviceId: uuid.New().String(),
				TokenID: uuid.New().String(),
				RefreshToken: uuid.New().String(),
				ExpirationDate:time.Now().Unix(),
				OrganizationId: uuid.New().String(),
				DeviceGroupId: uuid.New().String(),

			}
			err := provider.Add(&deviceToken)
			gomega.Expect(err).To(gomega.Succeed())

			deviceToken.RefreshToken = uuid.New().String()

			err = provider.Update(&deviceToken)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("should not be update to get", func(){

			deviceToken := entities.DeviceTokenData{
				DeviceId: uuid.New().String(),
				TokenID: uuid.New().String(),
				RefreshToken: uuid.New().String(),
				ExpirationDate:time.Now().Unix(),
				OrganizationId: uuid.New().String(),
				DeviceGroupId: uuid.New().String(),

			}

			err := provider.Update(&deviceToken)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})
}
