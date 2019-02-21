/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package interceptor

import (
	"context"
	"github.com/google/uuid"
	"github.com/nalej/authx/internal/app/authx/handler"
	"github.com/nalej/authx/internal/app/authx/manager"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-ping-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

// PingHandler that receives device pings and will be launched with a device interceptor.
type PingHandler struct {
}

func (ph * PingHandler) Ping(context.Context, *grpc_ping_go.PingRequest) (*gprc_ping_go.PingResponse, error){

}

var _ = ginkgo.Describe("Device Interceptor", func() {

	// gRPC server
	var server *grpc.Server
	var targetServer *grpc.Server

	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_authx_go.AuthxClient
	var interceptorClient grpc_authx_go.AuthxClient

	var targetDeviceGroup * grpc_authx_go.DeviceGroupCredentials
	var targetDevice * grpc_authx_go.DeviceCredentials

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
		interceptorClient = grpc_authx_go.NewAuthxClient(conn)
	})

	// Create a device group with a device linked to it
	ginkgo.BeforeEach(func() {
		// Add the device group
		deviceGroup := &grpc_authx_go.AddDeviceGroupCredentialsRequest{
			OrganizationId: uuid.New().String(),
			DeviceGroupId: uuid.New().String(),
			Enabled: true,
		}
		group, err := client.AddDeviceGroupCredentials(context.Background(), deviceGroup)
		gomega.Expect(err).To(gomega.Succeed())
		targetDeviceGroup = group
		// Add the device
		toAdd := grpc_authx_go.AddDeviceCredentialsRequest{
			OrganizationId: targetDeviceGroup.OrganizationId,
			DeviceGroupId: targetDeviceGroup.DeviceGroupId,
			DeviceId: uuid.New().String(),
		}
		added, err := client.AddDeviceCredentials(context.Background(), &toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(added.DeviceApiKey).NotTo(gomega.BeEmpty())
		targetDevice = added
	})

})