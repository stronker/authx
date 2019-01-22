/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package handler

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/nalej/authx/internal/app/authx/manager"
	"github.com/nalej/authx/pkg/token"
	pbAuthx "github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

var _ = ginkgo.Describe("Applications", func() {

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client pbAuthx.AuthxClient

	var mgr *manager.Authx

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		mgr = manager.NewAuthxMockup()
		handler := NewAuthx(mgr)

		pbAuthx.RegisterAuthxServer(server, handler)

		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = pbAuthx.NewAuthxClient(conn)
	})

	ginkgo.Context("with a role", func() {
		userName := "u1"
		organizationID := "o1"
		roleID := "r1"
		pass := "MyLittlePassword"

		ginkgo.BeforeEach(func() {
			role := &pbAuthx.Role{
				OrganizationId: organizationID,
				RoleId:         roleID,
				Name:           "rName1",
				Primitives:     []pbAuthx.AccessPrimitive{pbAuthx.AccessPrimitive_ORG},
			}
			success, err := client.AddRole(context.Background(), role)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())
		})

		ginkgo.It("add basic credentials with correct roleID", func() {
			success, err := client.AddBasicCredentials(context.Background(),
				&pbAuthx.AddBasicCredentialRequest{OrganizationId: organizationID,
					RoleId:   roleID,
					Username: userName,
					Password: pass,
				})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())
		})

		ginkgo.It("add basic credentials with incorrect roleID", func() {
			success, err := client.AddBasicCredentials(context.Background(),
				&pbAuthx.AddBasicCredentialRequest{OrganizationId: organizationID,
					RoleId:   roleID + "wrong",
					Username: userName,
					Password: pass,
				})
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(success).To(gomega.BeNil())
		})

		ginkgo.It("add basic credentials two times should fail", func() {
			success, err := client.AddBasicCredentials(context.Background(),
				&pbAuthx.AddBasicCredentialRequest{OrganizationId: organizationID,
					RoleId:   roleID,
					Username: userName,
					Password: pass,
				})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())

			success, err = client.AddBasicCredentials(context.Background(),
				&pbAuthx.AddBasicCredentialRequest{OrganizationId: organizationID,
					RoleId:   roleID,
					Username: userName,
					Password: pass,
				})
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(success).To(gomega.BeNil())
		})

		ginkgo.It("should be able to retrieve the user role", func(){
			success, err := client.AddBasicCredentials(context.Background(),
				&pbAuthx.AddBasicCredentialRequest{OrganizationId: organizationID,
					RoleId:   roleID,
					Username: userName,
					Password: pass,
				})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())

			userID := &grpc_user_go.UserId{
				OrganizationId:       organizationID,
				Email:                userName,
			}

			role, err := client.GetUserRole(context.Background(), userID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(role.OrganizationId).Should(gomega.Equal(organizationID))
			gomega.Expect(role.RoleId).Should(gomega.Equal(roleID))
		})

	})

	ginkgo.Context("with a basic credentials and two roleIDs", func() {
		userName := "u1"
		organizationID := "o1"
		roleID := "r1"
		roleID2 := "r2"
		pass := "MyLittlePassword"

		ginkgo.BeforeEach(func() {
			role := &pbAuthx.Role{
				OrganizationId: organizationID,
				RoleId:         roleID,
				Name:           "rName1",
				Primitives:     []pbAuthx.AccessPrimitive{pbAuthx.AccessPrimitive_ORG},
			}
			success, err := client.AddRole(context.Background(), role)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())

			role2 := &pbAuthx.Role{
				OrganizationId: organizationID,
				RoleId:         roleID2,
				Name:           "rName2",
				Primitives:     []pbAuthx.AccessPrimitive{pbAuthx.AccessPrimitive_ORG},
			}
			success, err = client.AddRole(context.Background(), role2)

			gomega.Expect(err).To(gomega.Succeed())

			success, err = client.AddBasicCredentials(context.Background(),
				&pbAuthx.AddBasicCredentialRequest{OrganizationId: organizationID,
					RoleId:   roleID,
					Username: userName,
					Password: pass,
				})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())
		})

		ginkgo.It("should login with correct password", func() {
			response, err := client.LoginWithBasicCredentials(context.Background(),
				&pbAuthx.LoginWithBasicCredentialsRequest{Username: userName, Password: pass})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(response).NotTo(gomega.BeNil())
		})

		ginkgo.It("should login with incorrect password", func() {
			response, err := client.LoginWithBasicCredentials(context.Background(),
				&pbAuthx.LoginWithBasicCredentialsRequest{Username: userName, Password: pass + "wrong"})

			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(status.Convert(err).Code()).Should(gomega.Equal(codes.Unauthenticated))
			gomega.Expect(response).To(gomega.BeNil())
		})

		ginkgo.It("should change password with correct password", func() {
			newPassword:=pass+"New"
			response,err := client.ChangePassword(context.Background(),
				&pbAuthx.ChangePasswordRequest{Username:userName,Password:pass,NewPassword:newPassword})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(response).NotTo(gomega.BeNil())

			loginResponse, err := client.LoginWithBasicCredentials(context.Background(),
				&pbAuthx.LoginWithBasicCredentialsRequest{Username: userName, Password: newPassword})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(loginResponse).NotTo(gomega.BeNil())
		})
		ginkgo.It("should change password with incorrect password", func() {
			newPassword:=pass+"New"
			response,err := client.ChangePassword(context.Background(),
				&pbAuthx.ChangePasswordRequest{Username:userName,Password:pass+"wrong",NewPassword:newPassword})

			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(status.Convert(err).Code()).Should(gomega.Equal(codes.Unauthenticated))
			gomega.Expect(response).To(gomega.BeNil())
		})

		ginkgo.It("should change password with incorrect username", func() {
			newPassword:=pass+"New"
			response,err := client.ChangePassword(context.Background(),
				&pbAuthx.ChangePasswordRequest{Username:userName+"wrong",Password:pass,NewPassword:newPassword})

			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(response).To(gomega.BeNil())
		})

		ginkgo.It("should change to a valid roleID", func() {
			success, err := client.EditUserRole(context.Background(),
				&pbAuthx.EditUserRoleRequest{Username: userName, NewRoleId: roleID2})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())
		})

		ginkgo.It("should change to a invalid roleID", func() {
			success, err := client.EditUserRole(context.Background(),
				&pbAuthx.EditUserRoleRequest{Username: userName, NewRoleId: roleID2+"wrong"})
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(success).To(gomega.BeNil())
		})
		ginkgo.It("should delete credentials", func() {
			success,err := client.DeleteCredentials(context.Background(),&pbAuthx.DeleteCredentialsRequest{Username:userName})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())
		})
		ginkgo.It("should delete wrong credentials", func() {
			success,err := client.DeleteCredentials(context.Background(),&pbAuthx.DeleteCredentialsRequest{Username:userName+"wrong"})
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(success).To(gomega.BeNil())
		})

		ginkgo.It("should refresh token", func() {
			response, err := client.LoginWithBasicCredentials(context.Background(),
				&pbAuthx.LoginWithBasicCredentialsRequest{Username: userName, Password: pass})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(response).NotTo(gomega.BeNil())

			tk, jwtErr := jwt.ParseWithClaims(response.Token, &token.Claim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(manager.DefaultSecret), nil
			})
			gomega.Expect(jwtErr).To(gomega.Succeed())
			gomega.Expect(tk).NotTo(gomega.BeNil())

			cl, ok := tk.Claims.(*token.Claim)
			gomega.Expect(ok).To(gomega.BeTrue())
			gomega.Expect(cl).NotTo(gomega.BeNil())

			newResponse, err := client.RefreshToken(context.Background(),
				&pbAuthx.RefreshTokenRequest{Token:response.Token,
				RefreshToken:response.RefreshToken})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(newResponse).NotTo(gomega.BeNil())

		})

		ginkgo.It("should reject invalid refresh token", func() {
			response, err := client.LoginWithBasicCredentials(context.Background(),
				&pbAuthx.LoginWithBasicCredentialsRequest{Username: userName, Password: pass})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(response).NotTo(gomega.BeNil())

			tk, jwtErr := jwt.ParseWithClaims(response.Token, &token.Claim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(manager.DefaultSecret), nil
			})
			gomega.Expect(jwtErr).To(gomega.Succeed())
			gomega.Expect(tk).NotTo(gomega.BeNil())

			cl, ok := tk.Claims.(*token.Claim)
			gomega.Expect(ok).To(gomega.BeTrue())
			gomega.Expect(cl).NotTo(gomega.BeNil())

			newResponse, err := client.RefreshToken(context.Background(),
				&pbAuthx.RefreshTokenRequest{Token:response.Token,
					RefreshToken:response.RefreshToken+"wrong"})
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(newResponse).To(gomega.BeNil())

		})

		ginkgo.It("should be able to retrieve the list roles", func(){
		    orgID := &grpc_organization_go.OrganizationId{
				OrganizationId:       organizationID,
			}
		    roles, err := client.ListRoles(context.Background(), orgID)
		    gomega.Expect(err).To(gomega.Succeed())
		    gomega.Expect(len(roles.Roles)).Should(gomega.Equal(2))
		})

	})

	ginkgo.Context("with device credentials", func() {

		var targetDeviceGroup * pbAuthx.DeviceGroupCredentials

		ginkgo.BeforeEach(func() {
			deviceGroup := &pbAuthx.AddDeviceGroupCredentialsRequest{
				OrganizationId: uuid.New().String(),
				DeviceGroupId: uuid.New().String(),
			}
			group, err := client.AddDeviceGroupCredentials(context.Background(), deviceGroup)
			gomega.Expect(err).To(gomega.Succeed())
			targetDeviceGroup = group
		})

		ginkgo.It("should be able to add a device credentials", func() {
			toAdd := pbAuthx.AddDeviceCredentialsRequest{
				OrganizationId: targetDeviceGroup.OrganizationId,
				DeviceGroupId: targetDeviceGroup.DeviceGroupId,
				DeviceId: uuid.New().String(),
			}
			added, err := client.AddDeviceCredentials(context.Background(), &toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added.DeviceApiKey).NotTo(gomega.BeEmpty())
		})
		ginkgo.It("should not be able to add a device credentials of an non existing group ", func() {
			toAdd := pbAuthx.AddDeviceCredentialsRequest{
				OrganizationId: targetDeviceGroup.OrganizationId,
				DeviceGroupId: uuid.New().String(),
				DeviceId: uuid.New().String(),
			}
			_, err := client.AddDeviceCredentials(context.Background(), &toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("should be able to update a device credentials", func() {
			toAdd := pbAuthx.AddDeviceCredentialsRequest{
				OrganizationId: targetDeviceGroup.OrganizationId,
				DeviceGroupId: targetDeviceGroup.DeviceGroupId,
				DeviceId: uuid.New().String(),
			}
			added, err := client.AddDeviceCredentials(context.Background(), &toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added.DeviceApiKey).NotTo(gomega.BeEmpty())

			toUpdate := pbAuthx.UpdateDeviceCredentialsRequest{
				OrganizationId: added.OrganizationId,
				DeviceGroupId: added.DeviceGroupId,
				DeviceId: added.DeviceId,
				Enabled: true,
			}
			success, err := client.UpdateDeviceCredentials(context.Background(), &toUpdate)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())

			// TODO: Check the credential has been updated

		})
		ginkgo.It("should not be able to update a non existing credentials", func() {
			toUpdate := pbAuthx.UpdateDeviceCredentialsRequest{
				OrganizationId: targetDeviceGroup.OrganizationId,
				DeviceGroupId: targetDeviceGroup.DeviceGroupId,
				DeviceId: uuid.New().String(),
				Enabled: true,
			}
			success, err := client.UpdateDeviceCredentials(context.Background(), &toUpdate)
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())

		})
		ginkgo.It("should be able to remove  a device credentials", func() {
			toAdd := pbAuthx.AddDeviceCredentialsRequest{
				OrganizationId: targetDeviceGroup.OrganizationId,
				DeviceGroupId: targetDeviceGroup.DeviceGroupId,
				DeviceId: uuid.New().String(),
			}
			added, err := client.AddDeviceCredentials(context.Background(), &toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added.DeviceApiKey).NotTo(gomega.BeEmpty())

			toRemove := grpc_device_go.DeviceId{
				OrganizationId: added.OrganizationId,
				DeviceGroupId: added.DeviceGroupId,
				DeviceId: added.DeviceId,
			}
			success, err := client.RemoveDeviceCredentials(context.Background(), &toRemove)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())

		})
		ginkgo.It("should not be able to remove a non existing credentials", func() {

			toRemove := grpc_device_go.DeviceId{
				OrganizationId: targetDeviceGroup.OrganizationId,
				DeviceGroupId: targetDeviceGroup.DeviceGroupId,
				DeviceId: uuid.New().String(),
			}
			success, err := client.RemoveDeviceCredentials(context.Background(), &toRemove)
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())

		})


	})

	ginkgo.Context("with device group credentials", func(){
		ginkgo.FIt("should be able to add a device group credentials", func() {
			toAdd := pbAuthx.AddDeviceGroupCredentialsRequest{
				OrganizationId:  uuid.New().String(),
				DeviceGroupId:  uuid.New().String(),
			}
			added, err := client.AddDeviceGroupCredentials(context.Background(), &toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added.DeviceGroupApiKey).NotTo(gomega.BeEmpty())
		})
		ginkgo.FIt("should not be able to add a device group credentials without organization_id", func() {
			toAdd := pbAuthx.AddDeviceGroupCredentialsRequest{
				DeviceGroupId:  uuid.New().String(),
			}
			_, err := client.AddDeviceGroupCredentials(context.Background(), &toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})

		ginkgo.FIt("should be able to update a device group credentials", func() {
			toAdd := pbAuthx.AddDeviceGroupCredentialsRequest{
				OrganizationId:  uuid.New().String(),
				DeviceGroupId:  uuid.New().String(),
			}
			added, err := client.AddDeviceGroupCredentials(context.Background(), &toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added.DeviceGroupApiKey).NotTo(gomega.BeEmpty())

			toUpdate := pbAuthx.UpdateDeviceGroupCredentialsRequest{
				OrganizationId:  toAdd.OrganizationId,
				DeviceGroupId:  toAdd.DeviceGroupId,
				UpdateEnabled: true,
				Enabled:true,
			}
			success, err := client.UpdateDeviceGroupCredentials(context.Background(), &toUpdate)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())
		})
		ginkgo.FIt("should not be able to update a device group credentials with all flags to false", func() {
			toAdd := pbAuthx.AddDeviceGroupCredentialsRequest{
				OrganizationId:  uuid.New().String(),
				DeviceGroupId:  uuid.New().String(),
			}
			added, err := client.AddDeviceGroupCredentials(context.Background(), &toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added.DeviceGroupApiKey).NotTo(gomega.BeEmpty())

			toUpdate := pbAuthx.UpdateDeviceGroupCredentialsRequest{
				OrganizationId:  toAdd.OrganizationId,
				DeviceGroupId:  toAdd.DeviceGroupId,
				UpdateEnabled: false,
				UpdateDeviceConnectivity: false,
			}
			success, err := client.UpdateDeviceGroupCredentials(context.Background(), &toUpdate)
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())
		})
		ginkgo.FIt("should not be able to update a non existing device group credentials", func() {

			toUpdate := pbAuthx.UpdateDeviceGroupCredentialsRequest{
				OrganizationId:  uuid.New().String(),
				DeviceGroupId:  uuid.New().String(),
				UpdateEnabled: true,
				Enabled: true,
			}
			success, err := client.UpdateDeviceGroupCredentials(context.Background(), &toUpdate)
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())
		})

		ginkgo.FIt("should be able to remove a device group credentials", func() {
			toAdd := pbAuthx.AddDeviceGroupCredentialsRequest{
				OrganizationId:  uuid.New().String(),
				DeviceGroupId:  uuid.New().String(),
			}
			added, err := client.AddDeviceGroupCredentials(context.Background(), &toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added.DeviceGroupApiKey).NotTo(gomega.BeEmpty())

			toRemove := grpc_device_go.DeviceGroupId{
				OrganizationId:  toAdd.OrganizationId,
				DeviceGroupId:  toAdd.DeviceGroupId,
			}
			success, err := client.RemoveDeviceGroupCredentials(context.Background(), &toRemove)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())
		})
		ginkgo.FIt("should not be able to remove a non existing device group credentials", func() {

			toRemove := grpc_device_go.DeviceGroupId{
				OrganizationId:   uuid.New().String(),
				DeviceGroupId:   uuid.New().String(),
			}
			success, err := client.RemoveDeviceGroupCredentials(context.Background(), &toRemove)
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())
		})

	})

	ginkgo.AfterEach(func() {
		err := mgr.Clean()
		gomega.Expect(err).To(gomega.Succeed())
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})
})
