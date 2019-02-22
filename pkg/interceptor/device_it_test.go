/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package interceptor

import (
	"context"
	"github.com/google/uuid"
	"github.com/nalej/authx/internal/app/authx/handler"
	"github.com/nalej/authx/internal/app/authx/manager"
	"github.com/nalej/authx/pkg/interceptor/devinterceptor"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-ping-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

// PingHandler that receives device pings and will be launched with a device interceptor. In this way, calls to the
// ping server will trigger the retrieval of the secret associated with the device group of a device request.
type PingHandler struct {
}

func (ph *PingHandler) Ping(ctx context.Context, request *grpc_ping_go.PingRequest) (*grpc_ping_go.PingResponse, error) {
	return &grpc_ping_go.PingResponse{
		RequestNumber: request.RequestNumber,
	}, nil
}

func getContextFromDeviceLogin(device grpc_authx_go.DeviceCredentials, authxClient grpc_authx_go.AuthxClient, header string) context.Context {
	loginRequest := &grpc_authx_go.DeviceLoginRequest{
		OrganizationId: device.OrganizationId,
		DeviceApiKey:   device.DeviceApiKey,
	}
	response, err := authxClient.DeviceLogin(context.Background(), loginRequest)
	gomega.Expect(err).To(gomega.Succeed())
	md := metadata.New(map[string]string{header: response.Token})
	log.Debug().Str("token", response.Token).Msg("Device token")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return ctx
}

var _ = ginkgo.Describe("Device Interceptor", func() {

	// gRPC server
	var authxServer *grpc.Server
	var targetServer *grpc.Server

	// grpc test listener
	var authxListener *bufconn.Listener
	var targetListener *bufconn.Listener

	// client
	var client grpc_authx_go.AuthxClient
	var pingClient grpc_ping_go.PingClient

	var targetDeviceGroup *grpc_authx_go.DeviceGroupCredentials
	var targetDevice *grpc_authx_go.DeviceCredentials

	cfg := NewConfig(&AuthorizationConfig{AllowsAll: false, Permissions: map[string]Permission{
		"/ping.Ping/Ping": {Must: []string{grpc_authx_go.AccessPrimitive_DEVICE.String()}},
	}}, "globalSecret", "authorization")

	// Setup the Authx server
	authxListener = test.GetDefaultListener()
	targetListener = bufconn.Listen(test.BufSize)

	authxServer = grpc.NewServer()
	// Launch Authx Server
	mgr := manager.NewAuthxMockup()
	handler := handler.NewAuthx(mgr)
	grpc_authx_go.RegisterAuthxServer(authxServer, handler)
	test.LaunchServer(authxServer, authxListener)
	conn, err := test.GetConn(*authxListener)
	if err != nil {
		ginkgo.Fail("cannot obtain connection " + err.Error())
	}
	client = grpc_authx_go.NewAuthxClient(conn)

	// Setup the ping server
	secretAccess, err := devinterceptor.NewMngtSecretAccessWithClient(client, devinterceptor.DefaultCacheEntries)
	targetServer = grpc.NewServer(WithDeviceAuthxInterceptor(secretAccess, cfg))

	// Launch Ping Server
	pingHandler := &PingHandler{}
	grpc_ping_go.RegisterPingServer(targetServer, pingHandler)
	test.LaunchServer(targetServer, targetListener)

	pingConn, err := test.GetConn(*targetListener)
	pingClient = grpc_ping_go.NewPingClient(pingConn)

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
		// Add the device
		toAdd := grpc_authx_go.AddDeviceCredentialsRequest{
			OrganizationId: targetDeviceGroup.OrganizationId,
			DeviceGroupId:  targetDeviceGroup.DeviceGroupId,
			DeviceId:       uuid.New().String(),
		}
		added, err := client.AddDeviceCredentials(context.Background(), &toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(added.DeviceApiKey).NotTo(gomega.BeEmpty())
		targetDevice = added
	})

	ginkgo.It("A valid device should be able to ping the server", func() {
		ctx := getContextFromDeviceLogin(*targetDevice, client, cfg.Header)
		request := &grpc_ping_go.PingRequest{
			RequestNumber: 1,
		}
		log.Debug().Interface("manager", mgr).Msg("Manager")
		response, err := pingClient.Ping(ctx, request)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(response.RequestNumber).Should(gomega.Equal(request.RequestNumber))
	})

})
