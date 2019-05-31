/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package devinterceptor

import (
	"context"
	"github.com/google/uuid"
	"github.com/nalej/authx/internal/app/authx/handler"
	"github.com/nalej/authx/internal/app/authx/manager"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var _ = ginkgo.Describe("Management secret access", func() {

	// gRPC server
	var server *grpc.Server

	// grpc test listener
	var listener *bufconn.Listener
	// client
	var cachedClient SecretAccess
	var client grpc_authx_go.AuthxClient

	var targetDeviceGroup *grpc_authx_go.DeviceGroupCredentials

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		// Register the handler with a mockup manager
		mgr := manager.NewAuthxMockup()
		handler := handler.NewAuthx(mgr)
		grpc_authx_go.RegisterAuthxServer(server, handler)

		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_authx_go.NewAuthxClient(conn)

		cachedClient, err = NewMngtSecretAccess("", 100)
		gomega.Expect(err).To(gomega.Succeed())
		cachedClient.(*MngtSecretAccess).Client = client

	})

	// Create a device group with a device linked to it
	ginkgo.BeforeEach(func() {
		// Add the device group
		deviceGroup := &grpc_authx_go.AddDeviceGroupCredentialsRequest{
			OrganizationId: uuid.New().String(),
			DeviceGroupId:  uuid.New().String(),
			Enabled:        true,
		}
		group, err := client.AddDeviceGroupCredentials(context.Background(), deviceGroup)
		gomega.Expect(err).To(gomega.Succeed())
		targetDeviceGroup = group
	})

	ginkgo.It("should be able to retrieve the secret from an existing device group", func() {
		deviceGroupId := &grpc_device_go.DeviceGroupId{
			OrganizationId: targetDeviceGroup.OrganizationId,
			DeviceGroupId:  targetDeviceGroup.DeviceGroupId,
		}
		secret, err := cachedClient.RetrieveSecret(deviceGroupId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(secret).ShouldNot(gomega.BeNil())
	})
	ginkgo.It("should be able to retrieve twice the secret from an existing device group", func() {
		deviceGroupId := &grpc_device_go.DeviceGroupId{
			OrganizationId: targetDeviceGroup.OrganizationId,
			DeviceGroupId:  targetDeviceGroup.DeviceGroupId,
		}
		secret, err := cachedClient.RetrieveSecret(deviceGroupId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(secret).ShouldNot(gomega.BeEmpty())
		secret, err = cachedClient.RetrieveSecret(deviceGroupId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(secret).ShouldNot(gomega.BeEmpty())
	})
	ginkgo.It("should not be able to retrieve the secret from an unknown device group", func() {
		deviceGroupId := &grpc_device_go.DeviceGroupId{
			OrganizationId: "not-found",
			DeviceGroupId:  targetDeviceGroup.DeviceGroupId,
		}
		secret, err := cachedClient.RetrieveSecret(deviceGroupId)
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(secret).To(gomega.BeEmpty())
	})

})
