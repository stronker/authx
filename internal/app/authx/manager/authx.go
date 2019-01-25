/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package manager

import (
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/authx/internal/app/authx/providers/credentials"
	"github.com/nalej/authx/internal/app/authx/providers/device"
	"github.com/nalej/authx/internal/app/authx/providers/role"
	"github.com/nalej/authx/pkg/token"
	"github.com/nalej/derrors"
	pbAuthx "github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-user-go"
	"time"
)

// DefaultExpirationDuration is the default duration used in the mockup.
const DefaultExpirationDuration = "10h"
const DefaultDeviceExpirationDuration = "10m"

// DefaultSecret is the default secret used in the mockup.
const DefaultSecret = "MyLittleSecret"

// Authx is the component that manages the business logic.
type Authx struct {
	Password            Password
	Token               Token
	CredentialsProvider credentials.BasicCredentials
	RoleProvider        role.Role
	DeviceProvider      device.Provider
	secret             string
	expirationDuration time.Duration
	DeviceToken			DeviceToken
	DeviceExpiration    time.Duration
}

// NewAuthx creates a new manager.
func NewAuthx(password Password, tokenManager Token, deviceToken DeviceToken, credentialsProvider credentials.BasicCredentials,
	roleProvide role.Role, deviceProvider device.Provider, secret string, expirationDuration time.Duration, deviceExpiration time.Duration) *Authx {

	return &Authx{
		Password: password,
		Token: tokenManager,
		CredentialsProvider: credentialsProvider,
		RoleProvider: roleProvide,
		DeviceProvider:deviceProvider,
		secret: secret,
		expirationDuration: expirationDuration,
		DeviceToken: deviceToken,
		DeviceExpiration: deviceExpiration,

	}

}

// NewAuthxMockup create a new mockup manager.
func NewAuthxMockup() *Authx {
	d, _ := time.ParseDuration(DefaultExpirationDuration)
	e, _ := time.ParseDuration(DefaultDeviceExpirationDuration)
	return NewAuthx(NewBCryptPassword(), NewJWTTokenMockup(), NewJWTDeviceTokenMockup(),
		credentials.NewBasicCredentialMockup(), role.NewRoleMockup(),
		device.NewMockupDeviceCredentialsProvider(),DefaultSecret, d, e)
}

// DeleteCredentials deletes the credential for a specific username.
func (m *Authx) DeleteCredentials(username string) derrors.Error {
	return m.CredentialsProvider.Delete(username)
}

// AddBasicCredentials generate credential for a specific user.
func (m *Authx) AddBasicCredentials(username string, organizationID string, roleID string, password string) derrors.Error {
	_, err := m.RoleProvider.Get(organizationID, roleID)
	if err != nil {
		return err
	}

	exist, err := m.CredentialsProvider.Exist(username)
	if err != nil {
		return err
	}
	if *exist {
		return derrors.NewAlreadyExistsError("credentials already exists")
	}

	hashedPassword, err := m.Password.GenerateHashedPassword(password)
	if err != nil {
		return err
	}

	entity := entities.NewBasicCredentialsData(username, hashedPassword, roleID, organizationID)
	return m.CredentialsProvider.Add(entity)
}

// LoginWithBasicCredentials check the password and returns a valid token.
func (m *Authx) LoginWithBasicCredentials(username string, password string) (*pbAuthx.LoginResponse, derrors.Error) {
	credentials, err := m.CredentialsProvider.Get(username)
	if err != nil {
		return nil, err
	}
	err = m.Password.CompareHashAndPassword(credentials.Password, password)
	if err != nil {
		return nil, err
	}
	role, err := m.RoleProvider.Get(credentials.OrganizationID, credentials.RoleID)
	if err != nil {
		return nil, err
	}
	personalClaim := token.NewPersonalClaim(username, role.Name, role.Primitives, credentials.OrganizationID)
	gToken, err := m.Token.Generate(personalClaim, m.expirationDuration, m.secret, false)
	if err != nil {

		return nil, err
	}
	response := &pbAuthx.LoginResponse{Token: gToken.Token, RefreshToken: gToken.RefreshToken}
	return response, nil
}

func (m *Authx) ChangePassword(username string, password string, newPassword string) derrors.Error {
	credentials, err := m.CredentialsProvider.Get(username)
	if err != nil {
		return err
	}
	err = m.Password.CompareHashAndPassword(credentials.Password, password)
	if err != nil {
		return err
	}
	hashedPassword, err := m.Password.GenerateHashedPassword(newPassword)
	if err != nil {
		return err
	}
	edit := entities.NewEditBasicCredentialsData().WithPassword(hashedPassword)
	return m.CredentialsProvider.Edit(username, edit)
}

// RefreshToken renew an old token.
func (m *Authx) RefreshToken(oldToken string, refreshToken string) (*pbAuthx.LoginResponse, derrors.Error) {
	gToken, err := m.Token.Refresh(oldToken, refreshToken, m.expirationDuration, m.secret)
	if err != nil {
		return nil, err
	}
	response := &pbAuthx.LoginResponse{Token: gToken.Token, RefreshToken: gToken.RefreshToken}
	return response, nil
}

// AddRole add a new role to the authorization system.
func (m *Authx) AddRole(role *pbAuthx.Role) derrors.Error {
	entity := entities.NewRoleData(role.OrganizationId, role.RoleId, role.Name, role.Internal, PrimitivesToString(role.Primitives))
	return m.RoleProvider.Add(entity)
}

// EditUserRole change the RoleID to a specific user.
func (m *Authx) EditUserRole(username string, roleID string) derrors.Error {
	credentials, err := m.CredentialsProvider.Get(username)
	if err != nil {
		return err
	}
	role, err := m.RoleProvider.Get(credentials.OrganizationID, roleID)
	if err != nil {
		return err
	}

	if role == nil {
		return derrors.NewNotFoundError("role not found")
	}

	edit := entities.NewEditBasicCredentialsData().WithRoleID(roleID)
	return m.CredentialsProvider.Edit(username, edit)
}

func (m * Authx) ListRoles(organizationID * grpc_organization_go.OrganizationId) ([]entities.RoleData, derrors.Error){
	return m.RoleProvider.List(organizationID.OrganizationId)
}

func (m * Authx) GetUserRole(userID * grpc_user_go.UserId)( * entities.RoleData, derrors.Error){
	cred, err := m.CredentialsProvider.Get(userID.Email)
	if err != nil{
		return nil, err
	}
	return m.RoleProvider.Get(userID.OrganizationId, cred.RoleID)
}

// Clean removes all the data.
func (m *Authx) Clean() derrors.Error {
	err := m.Token.Clean()
	if err != nil {
		return err
	}
	err = m.CredentialsProvider.Truncate()
	if err != nil {
		return err
	}
	err = m.RoleProvider.Truncate()
	if err != nil {
		return err
	}
	err = m.DeviceProvider.Truncate()
	if err != nil {
		return err
	}

	return nil
}

// PrimitivesToString transform the primitive to strings.
func PrimitivesToString(primitives [] pbAuthx.AccessPrimitive) [] string {
	strPrimitives := make([] string, 0, len(primitives))
	for _, p := range primitives {
		strPrimitives = append(strPrimitives, p.String())
	}
	return strPrimitives
}

// -- Device Credentials -- //
func (m * Authx) AddDeviceCredentials(deviceCredentials * pbAuthx.AddDeviceCredentialsRequest) (*entities.DeviceCredentials, derrors.Error) {

	// Check if the group exists
	exists, err := m.DeviceProvider.ExistsDeviceGroup(deviceCredentials.OrganizationId, deviceCredentials.DeviceGroupId)
	if err != nil {
		return nil, err
	}
	if ! exists{
		return nil, derrors.NewNotFoundError("deviceGroupID").WithParams(deviceCredentials.OrganizationId, deviceCredentials.DeviceGroupId)
	}

	// Get the group to review if it is enable
	group, err := m.DeviceProvider.GetDeviceGroup(deviceCredentials.OrganizationId, deviceCredentials.DeviceGroupId)
	if err != nil {
		return nil, err
	}
	if ! group.Enabled {
		return nil, derrors.NewPermissionDeniedError("the group is temporarily disabled").WithParams(deviceCredentials.OrganizationId, deviceCredentials.DeviceGroupId)
	}

	toAdd := entities.NewDeviceCredentialsFromGRPC(deviceCredentials)
	err =  m.DeviceProvider.AddDeviceCredentials(toAdd)
	if err != nil{
		return nil, err
	}
	return toAdd, nil
}

func (m * Authx) UpdateDeviceCredentials (deviceCredentials * pbAuthx.UpdateDeviceCredentialsRequest) derrors.Error {
	// Check if the credentials exists
	exists, err := m.DeviceProvider.ExistsDevice(deviceCredentials.OrganizationId, deviceCredentials.DeviceGroupId, deviceCredentials.DeviceId)
	if err != nil {
		return err
	}
	if !exists{
		return derrors.NewNotFoundError("device credentials").WithParams(deviceCredentials.OrganizationId, deviceCredentials.DeviceGroupId, deviceCredentials.DeviceId)
	}

	toUpdate, err := m.DeviceProvider.GetDevice(deviceCredentials.OrganizationId, deviceCredentials.DeviceGroupId, deviceCredentials.DeviceId)
	if err != nil {
		return err
	}

	toUpdate.Enabled = deviceCredentials.Enabled

	err = m.DeviceProvider.UpdateDeviceCredentials(toUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (m * Authx) GetDeviceCredentials(request *grpc_device_go.DeviceId) (*entities.DeviceCredentials, derrors.Error) {
	// Check if the credentials group exist
	exists, err := m.DeviceProvider.ExistsDeviceGroup(request.OrganizationId, request.DeviceGroupId)
	if err != nil {
		return nil, err
	}
	if !exists{
		return nil, derrors.NewNotFoundError("device group credentials").WithParams(request.OrganizationId, request.DeviceGroupId)
	}
	// Check if the credentials device exist
	exists, err = m.DeviceProvider.ExistsDevice(request.OrganizationId, request.DeviceGroupId, request.DeviceId)
	if err != nil {
		return nil, err
	}
	if !exists{
		return nil, derrors.NewNotFoundError("device credentials").WithParams(request.OrganizationId, request.DeviceGroupId, request.DeviceId)
	}

	credentials, err := m.DeviceProvider.GetDevice(request.OrganizationId, request.DeviceGroupId, request.DeviceId)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

func (m * Authx) RemoveDeviceCredentials (deviceCredentials * grpc_device_go.DeviceId) derrors.Error {

	exists, err := m.DeviceProvider.ExistsDevice(deviceCredentials.OrganizationId, deviceCredentials.DeviceGroupId, deviceCredentials.DeviceId)
	if err != nil {
		return err
	}
	if !exists{
		return derrors.NewNotFoundError("device credentials").WithParams(deviceCredentials.OrganizationId, deviceCredentials.DeviceGroupId, deviceCredentials.DeviceId)
	}

	err = m.DeviceProvider.RemoveDevice(deviceCredentials.OrganizationId, deviceCredentials.DeviceGroupId, deviceCredentials.DeviceId)
	if err != nil {
		return err
	}

	return nil
}

func (m * Authx) LoginDeviceCredentials (loginRequest * pbAuthx.DeviceLoginRequest) (*pbAuthx.LoginResponse, derrors.Error) {

	credentials, err := m.DeviceProvider.GetDeviceByApiKey(loginRequest.DeviceApiKey)
	if err != nil {
		return nil, err
	}

	if credentials.OrganizationID != loginRequest.OrganizationId {
		return nil, derrors.NewUnauthenticatedError("Invalid credentials")
	}

	// Get the group to review if it is enable
	group, err := m.DeviceProvider.GetDeviceGroup(credentials.OrganizationID, credentials.DeviceGroupID)
	if err != nil {
		return nil, err
	}
	// if the group is disabled, the login is not allowed
	if ! group.Enabled {
		return nil, derrors.NewPermissionDeniedError("the group is temporarily disabled").WithParams(credentials.OrganizationID, credentials.DeviceGroupID)
	}

	deviceClaim := token.NewDeviceClaim(credentials.OrganizationID, credentials.DeviceGroupID, credentials.DeviceGroupID)

	gToken, err := m.DeviceToken.Generate(deviceClaim, m.DeviceExpiration, m.secret, false)
	if err != nil {

		return nil, err
	}
	response := &pbAuthx.LoginResponse{Token: gToken.Token, RefreshToken: gToken.RefreshToken}

	return response, nil

}

func (m * Authx) AddDeviceGroupCredentials(groupCredentials *pbAuthx.AddDeviceGroupCredentialsRequest) (*entities.DeviceGroupCredentials, derrors.Error){

	toAdd := entities.NewDeviceGroupCredentialsFromGRPC(groupCredentials)
	err := m.DeviceProvider.AddDeviceGroupCredentials(toAdd)
	if err != nil{
		return nil, err
	}
	return toAdd, nil
}

func (m * Authx) UpdateDeviceGroupCredentials(groupCredentials * pbAuthx.UpdateDeviceGroupCredentialsRequest) derrors.Error {
	// Check if the credentials exists
	exists, err := m.DeviceProvider.ExistsDeviceGroup(groupCredentials.OrganizationId, groupCredentials.DeviceGroupId)
	if err != nil {
		return err
	}
	if !exists{
		return derrors.NewNotFoundError("device group credentials").WithParams(groupCredentials.OrganizationId, groupCredentials.DeviceGroupId)
	}

	toUpdate, err := m.DeviceProvider.GetDeviceGroup(groupCredentials.OrganizationId, groupCredentials.DeviceGroupId)
	if err != nil {
		return err
	}
	if groupCredentials.UpdateEnabled {
		toUpdate.Enabled = groupCredentials.Enabled
	}
	if groupCredentials.UpdateDeviceConnectivity {
		toUpdate.DefaultDeviceConnectivity = groupCredentials.DefaultDeviceConnectivity
	}

	err = m.DeviceProvider.UpdateDeviceGroupCredentials(toUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (m * Authx) GetDeviceGroupCredentials(request *grpc_device_go.DeviceGroupId) (*entities.DeviceGroupCredentials, derrors.Error) {
	// Check if the credentials exists
	exists, err := m.DeviceProvider.ExistsDeviceGroup(request.OrganizationId, request.DeviceGroupId)
	if err != nil {
		return nil, err
	}
	if !exists{
		return nil, derrors.NewNotFoundError("device group credentials").WithParams(request.OrganizationId, request.DeviceGroupId)
	}

	group, err := m.DeviceProvider.GetDeviceGroup(request.OrganizationId, request.DeviceGroupId)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (m * Authx) RemoveDeviceGroupCredentials(groupCredentials * grpc_device_go.DeviceGroupId) derrors.Error {

	exists, err := m.DeviceProvider.ExistsDeviceGroup(groupCredentials.OrganizationId, groupCredentials.DeviceGroupId)
	if err != nil {
		return err
	}
	if !exists{
		return derrors.NewNotFoundError("device group credentials").WithParams(groupCredentials.OrganizationId, groupCredentials.DeviceGroupId)
	}

	err = m.DeviceProvider.RemoveDeviceGroup(groupCredentials.OrganizationId, groupCredentials.DeviceGroupId)
	if err != nil {
		return err
	}

	return nil
}

func (m * Authx) LoginDeviceGroup (credentials *pbAuthx.DeviceGroupLoginRequest) derrors.Error  {

	group, err := m.DeviceProvider.GetDeviceGroupByApiKey(credentials.DeviceGroupApiKey)
	if err != nil {
		return err
	}
	if group.OrganizationID != credentials.OrganizationId{
		return derrors.NewUnauthenticatedError("Invalid credentials")
	}
	// if the group is disabled, the login is not allowed
	if ! group.Enabled {
		return derrors.NewPermissionDeniedError("the group is temporarily disabled").WithParams(group.OrganizationID, group.DeviceGroupID)
	}
	return nil

}

// RefreshDeviceToken renew an old token.
func (m *Authx) RefreshDeviceToken(oldToken string, refreshToken string) (*pbAuthx.LoginResponse, derrors.Error) {

	claim, err := m.DeviceToken.GetTokenInfo(oldToken, m.secret)

	// get the group to check if it is enabled
	group, err := m.DeviceProvider.GetDeviceGroup(claim.OrganizationID, claim.DeviceGroupID)
	if err != nil {
		return nil, err
	}
	if ! group.Enabled {
		return nil, derrors.NewPermissionDeniedError("the group is temporarily disabled").WithParams(group.OrganizationID, group.DeviceGroupID)
	}


	gToken, err := m.DeviceToken.Refresh(oldToken, refreshToken, m.expirationDuration, m.secret)
	if err != nil {
		return nil, err
	}
	response := &pbAuthx.LoginResponse{Token: gToken.Token, RefreshToken: gToken.RefreshToken}
	return response, nil
}