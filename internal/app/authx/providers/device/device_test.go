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

package device

import (
	"github.com/google/uuid"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/stronker/authx/internal/app/authx/entities"
	"github.com/stronker/authx/internal/app/authx/utils"
)

func RunTest(provider Provider) {
	
	ginkgo.AfterEach(func() {
		provider.Truncate()
	})
	
	testHelper := utils.NewDeviceTestHepler()
	
	ginkgo.Context("device group credential tests", func() {
		ginkgo.It("Should be able to add device group", func() {
			
			toAdd := testHelper.CreateDeviceGroupCredentials()
			
			err := provider.AddDeviceGroupCredentials(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("Should not be able to add device group twice", func() {
			
			toAdd := testHelper.CreateDeviceGroupCredentials()
			
			err := provider.AddDeviceGroupCredentials(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			
			err = provider.AddDeviceGroupCredentials(toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("Should be able to update device group", func() {
			
			toAdd := testHelper.CreateDeviceGroupCredentials()
			
			err := provider.AddDeviceGroupCredentials(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			
			toAdd.DefaultDeviceConnectivity = false
			toAdd.Enabled = false
			
			err = provider.UpdateDeviceGroupCredentials(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			
			// check the update worked
			updated, err := provider.GetDeviceGroup(toAdd.OrganizationID, toAdd.DeviceGroupID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(updated.Enabled).Should(gomega.Equal(toAdd.Enabled))
			gomega.Expect(updated.DefaultDeviceConnectivity).Should(gomega.Equal(toAdd.DefaultDeviceConnectivity))
			
		})
		ginkgo.It("Should not be able to update non existing device group", func() {
			
			toUpdate := testHelper.CreateDeviceGroupCredentials()
			
			err := provider.UpdateDeviceGroupCredentials(toUpdate)
			gomega.Expect(err).NotTo(gomega.Succeed())
			
		})
		ginkgo.It("Should be able to find a device group", func() {
			toAdd := testHelper.CreateDeviceGroupCredentials()
			
			err := provider.AddDeviceGroupCredentials(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			
			exists, err := provider.ExistsDeviceGroup(toAdd.OrganizationID, toAdd.DeviceGroupID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())
			
		})
		ginkgo.It("Should not be able to find a non existing device group", func() {
			toAdd := testHelper.CreateDeviceGroupCredentials()
			
			exists, err := provider.ExistsDeviceGroup(toAdd.OrganizationID, toAdd.DeviceGroupID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())
			
		})
		ginkgo.It("Should be able to get a device group", func() {
			toAdd := testHelper.CreateDeviceGroupCredentials()
			
			err := provider.AddDeviceGroupCredentials(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			
			recovered, err := provider.GetDeviceGroup(toAdd.OrganizationID, toAdd.DeviceGroupID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(recovered).NotTo(gomega.BeNil())
			gomega.Expect(recovered.DeviceGroupApiKey).Should(gomega.Equal(toAdd.DeviceGroupApiKey))
			
		})
		ginkgo.It("Should not be able to get a non existing device group", func() {
			toAdd := testHelper.CreateDeviceGroupCredentials()
			
			_, err := provider.GetDeviceGroup(toAdd.OrganizationID, toAdd.DeviceGroupID)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("Should be able to get a device group by group_api_key", func() {
			toAdd := testHelper.CreateDeviceGroupCredentials()
			
			err := provider.AddDeviceGroupCredentials(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			
			recovered, err := provider.GetDeviceGroupByApiKey(toAdd.DeviceGroupApiKey)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(recovered).NotTo(gomega.BeNil())
			gomega.Expect(recovered.DeviceGroupApiKey).Should(gomega.Equal(toAdd.DeviceGroupApiKey))
			gomega.Expect(recovered.DeviceGroupID).Should(gomega.Equal(toAdd.DeviceGroupID))
			
		})
		ginkgo.It("Should be able to remove a device group", func() {
			toAdd := testHelper.CreateDeviceGroupCredentials()
			
			err := provider.AddDeviceGroupCredentials(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			
			err = provider.RemoveDeviceGroup(toAdd.OrganizationID, toAdd.DeviceGroupID)
			gomega.Expect(err).To(gomega.Succeed())
			
			// Should not be able to find it
			exists, err := provider.ExistsDeviceGroup(toAdd.OrganizationID, toAdd.DeviceGroupID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())
			
		})
		ginkgo.It("Should not be able to remove a non existing device group", func() {
			toAdd := testHelper.CreateDeviceGroupCredentials()
			
			err := provider.RemoveDeviceGroup(toAdd.OrganizationID, toAdd.DeviceGroupID)
			gomega.Expect(err).NotTo(gomega.Succeed())
			
		})
	})
	ginkgo.Context("device credential tests", func() {
		var targetDeviceGroup *entities.DeviceGroupCredentials
		ginkgo.BeforeEach(func() {
			targetDeviceGroup = testHelper.CreateDeviceGroupCredentials()
			err := provider.AddDeviceGroupCredentials(targetDeviceGroup)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("Should be able to add device credentials ", func() {
			toAdd := testHelper.CreateDeviceCredentials(*targetDeviceGroup)
			err := provider.AddDeviceCredentials(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("Should not be able to add device credentials of a non existing group", func() {
			toAdd := testHelper.CreateDeviceCredentials(*targetDeviceGroup)
			toAdd.DeviceGroupID = uuid.New().String()
			
			err := provider.AddDeviceCredentials(toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("Should be able to update device credentials ", func() {
			toAdd := testHelper.CreateDeviceCredentials(*targetDeviceGroup)
			err := provider.AddDeviceCredentials(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			
			toAdd.Enabled = false
			err = provider.UpdateDeviceCredentials(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			
			// checks th update works
			retrieved, err := provider.GetDevice(toAdd.OrganizationID, toAdd.DeviceGroupID, toAdd.DeviceID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved.Enabled).NotTo(gomega.BeTrue())
			
		})
		ginkgo.It("Should not be able to update device credentials ", func() {
			toAdd := testHelper.CreateDeviceCredentials(*targetDeviceGroup)
			
			err := provider.UpdateDeviceCredentials(toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())
			
		})
		ginkgo.It("Should be able to find device credentials ", func() {
			toAdd := testHelper.CreateDeviceCredentials(*targetDeviceGroup)
			err := provider.AddDeviceCredentials(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			
			exists, err := provider.ExistsDevice(toAdd.OrganizationID, toAdd.DeviceGroupID, toAdd.DeviceID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
		ginkgo.It("Should not be able to find device credentials ", func() {
			toAdd := testHelper.CreateDeviceCredentials(*targetDeviceGroup)
			exists, err := provider.ExistsDevice(toAdd.OrganizationID, toAdd.DeviceGroupID, toAdd.DeviceID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())
		})
		ginkgo.It("Should be able to return a device credentials ", func() {
			toAdd := testHelper.CreateDeviceCredentials(*targetDeviceGroup)
			err := provider.AddDeviceCredentials(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			
			retrieved, err := provider.GetDevice(toAdd.OrganizationID, toAdd.DeviceGroupID, toAdd.DeviceID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved.DeviceApiKey).Should(gomega.Equal(toAdd.DeviceApiKey))
		})
		ginkgo.It("Should not be able to return a device credentials ", func() {
			toAdd := testHelper.CreateDeviceCredentials(*targetDeviceGroup)
			
			_, err := provider.GetDevice(toAdd.OrganizationID, toAdd.DeviceGroupID, toAdd.DeviceID)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("Should be able to remove a device credentials ", func() {
			toAdd := testHelper.CreateDeviceCredentials(*targetDeviceGroup)
			err := provider.AddDeviceCredentials(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			
			err = provider.RemoveDevice(toAdd.OrganizationID, toAdd.DeviceGroupID, toAdd.DeviceID)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("Should not be able to remove a device credentials ", func() {
			toAdd := testHelper.CreateDeviceCredentials(*targetDeviceGroup)
			
			err := provider.RemoveDevice(toAdd.OrganizationID, toAdd.DeviceGroupID, toAdd.DeviceID)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		
	})
}
