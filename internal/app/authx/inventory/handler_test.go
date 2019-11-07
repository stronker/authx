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
	"context"
	"github.com/google/uuid"
	"github.com/nalej/authx/internal/app/authx/config"
	inventoryProv "github.com/nalej/authx/internal/app/authx/providers/inventory"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"time"
)

func createTestConfig() *config.Config {
	return &config.Config{
		ManagementClusterCert: "ManagementCertContent",
		EdgeControllerExpTime: time.Minute,
	}
}

var _ = ginkgo.Describe("Asset service", func() {
	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_authx_go.InventoryClient

	// Providers
	var inventoryProvider inventoryProv.Provider

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		// Register the service
		inventoryProvider = inventoryProv.NewMockupInventoryProvider()
		manager := NewManager(inventoryProvider, *createTestConfig())
		handler := NewHandler(manager)
		grpc_authx_go.RegisterInventoryServer(server, handler)
		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_authx_go.NewInventoryClient(conn)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})

	ginkgo.BeforeEach(func() {
		ginkgo.By("cleaning the mockups", func() {
			inventoryProvider.Clear()
		})
	})

	ginkgo.It("should be able to create a token", func() {
		orgID := &grpc_organization_go.OrganizationId{
			OrganizationId: uuid.New().String(),
		}
		token, err := client.CreateEICJoinToken(context.Background(), orgID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(token.OrganizationId).Should(gomega.Equal(orgID.OrganizationId))
		gomega.Expect(token.Token).ShouldNot(gomega.BeEmpty())
		gomega.Expect(token.Cacert).Should(gomega.Equal(createTestConfig().ManagementClusterCert))
	})

	ginkgo.It("should be able to use a valid join token", func() {
		orgID := &grpc_organization_go.OrganizationId{
			OrganizationId: uuid.New().String(),
		}
		token, err := client.CreateEICJoinToken(context.Background(), orgID)
		gomega.Expect(err).To(gomega.Succeed())
		joinRequest := &grpc_authx_go.EICJoinRequest{
			OrganizationId: token.OrganizationId,
			Token:          token.Token,
		}
		success, err := client.ValidEICJoinToken(context.Background(), joinRequest)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(success).ShouldNot(gomega.BeNil())
	})
})
