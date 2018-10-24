/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package handler

import (
	"context"
	"github.com/nalej/authx/internal/app/authx/manager"
	"github.com/nalej/derrors"
	pbAuthx "github.com/nalej/grpc-authx-go"
	pbCommon "github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
)

// Authx is the struct that handles the gRPC service.
type Authx struct {
	// Manager is the struct responsible of the service business logic.
	Manager *manager.Authx
}

// NewAuthx creates a new handler.
func NewAuthx(manager *manager.Authx) *Authx {
	return &Authx{Manager: manager}
}

// DeleteCredentials remove an existing credential using the username.
func (h *Authx) DeleteCredentials(_ context.Context, request *pbAuthx.DeleteCredentialsRequest) (*pbCommon.Success, error) {
	if request.Username == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("username is mandatory"))
	}
	err := h.Manager.DeleteCredentials(request.Username)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &pbCommon.Success{}, nil
}

// AddBasicCredentials adds a new credential specifying a password.
func (h *Authx) AddBasicCredentials(_ context.Context, request *pbAuthx.AddBasicCredentialRequest) (*pbCommon.Success, error) {
	if request.Username == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("username is mandatory"))
	}
	if request.OrganizationId == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("organizationID is mandatory"))
	}
	if request.RoleId == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("roleID is mandatory"))
	}
	if request.Password == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("password is mandatory"))
	}

	err := h.Manager.AddBasicCredentials(request.Username, request.OrganizationId, request.RoleId, request.Password)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &pbCommon.Success{}, nil
}

// LoginWithBasicCredentials login in the system and recovers a auth token.
func (h *Authx) LoginWithBasicCredentials(_ context.Context, request *pbAuthx.LoginWithBasicCredentialsRequest) (*pbAuthx.LoginResponse, error) {
	if request.Username == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("username is mandatory"))
	}
	if request.Password == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("password is mandatory"))
	}

	response, err := h.Manager.LoginWithBasicCredentials(request.Username, request.Password)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return response, nil
}

// ChangePassword update an existing password.
func (h *Authx) ChangePassword(ctx context.Context, request *pbAuthx.ChangePasswordRequest) (*pbCommon.Success, error){
	if request.Username == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("username is mandatory"))
	}
	if request.Password == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("password is mandatory"))
	}

	if request.NewPassword == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("newPassword is mandatory"))
	}

	err := h.Manager.ChangePassword(request.Username, request.Password, request.NewPassword)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &pbCommon.Success{}, nil
}

// RefreshToken renews an existing token.
func (h *Authx) RefreshToken(_ context.Context, request *pbAuthx.RefreshTokenRequest) (*pbAuthx.LoginResponse, error) {
	if request.Username == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("username is mandatory"))
	}
	if request.RefreshToken == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("refreshToken is mandatory"))
	}
	//TODO maybe we should revisit this, to send the old token instead of the tokenID could be a good idea.
	if request.TokenId == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("tokenID is mandatory"))
	}

	response, err := h.Manager.RefreshToken(request.Username, request.TokenId, request.RefreshToken)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return response, nil
}

// AddRole adds a role with a authorization properties.
func (h *Authx) AddRole(_ context.Context, request *pbAuthx.Role) (*pbCommon.Success, error) {
	if request.RoleId == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("roleID is mandatory"))
	}
	if request.Name == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("name is mandatory"))
	}
	if request.OrganizationId == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("organizationID is mandatory"))
	}
	if request.Primitives == nil || len(request.Primitives) == 0 {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("primitives is mandatory"))
	}

	err := h.Manager.AddRole(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &pbCommon.Success{}, nil

}

// EditUserRole change the roleID to a specific user.
func (h *Authx) EditUserRole(_ context.Context, request *pbAuthx.EditUserRoleRequest) (*pbCommon.Success, error) {
	if request.Username == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("username is mandatory"))
	}
	if request.NewRoleId == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("newRoleID is mandatory"))
	}
	err := h.Manager.EditUserRole(request.Username, request.NewRoleId)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &pbCommon.Success{}, nil
}
