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
	pbUser "github.com/nalej/grpc-user-go"
)

type AuthxServer struct {
	manager manager.Authx
}

func NewAuthxServer() *AuthxServer {
	return &AuthxServer{}
}

func (h *AuthxServer) DeleteCredentials(_ context.Context, request *pbUser.UserId) (*pbCommon.Success, error) {
	panic("implement me")
}

func (h *AuthxServer) AddBasicCredentials(_ context.Context, request *pbAuthx.AddBasicCredentialRequest) (*pbCommon.Success, error) {
	if request.Username == "" {
		return nil, derrors.NewOperationError("username is mandatory")
	}
	if request.OrganizationId == "" {
		return nil, derrors.NewOperationError("organizationID is mandatory")
	}
	if request.RoleId == "" {
		return nil, derrors.NewOperationError("roleID is mandatory")
	}
	if request.Password == "" {
		return nil, derrors.NewOperationError("password is mandatory")
	}

	err := h.manager.AddBasicCredentials(request.Username, request.OrganizationId, request.RoleId, request.Password)
	if err != nil {
		return nil, err
	}
	return &pbCommon.Success{}, nil
}

func (h *AuthxServer) LoginWithBasicCredentials(_ context.Context, request *pbAuthx.LoginWithBasicCredentialsRequest) (*pbAuthx.LoginResponse, error) {
	if request.Username == "" {
		return nil, derrors.NewOperationError("username is mandatory")
	}
	if request.Password == "" {
		return nil, derrors.NewOperationError("password is mandatory")
	}
	return h.manager.LoginWithBasicCredentials(request.Username, request.Password)
}

func (h *AuthxServer) RefreshToken(_ context.Context, request *pbAuthx.RefreshTokenRequest) (*pbAuthx.LoginResponse, error) {
	panic("implement me")
}

func (h *AuthxServer) AddRole(_ context.Context, request *pbAuthx.Role) (*pbCommon.Success, error) {
	if request.RoleId == "" {
		return nil, derrors.NewOperationError("roleID is mandatory")
	}
	if request.Name == "" {
		return nil, derrors.NewOperationError("name is mandatory")
	}
	if request.OrganizationId == "" {
		return nil, derrors.NewOperationError("organizationID is mandatory")
	}
	if request.Primitives == nil || len(request.Primitives) == 0 {
		return nil, derrors.NewOperationError("primitives is mandatory")
	}

	err := h.manager.AddRole(request)
	if err != nil {
		return nil, err
	}
	return &pbCommon.Success{}, nil

}

func (h *AuthxServer) EditUserRole(_ context.Context, request *pbAuthx.EditUserRoleRequest) (*pbCommon.Success, error) {
	panic("implement me")
}
