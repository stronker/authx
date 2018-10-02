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
)

type AuthxServer struct {
	manager manager.Authx
}

func NewAuthxServer() *AuthxServer {
	return &AuthxServer{}
}

func (h *AuthxServer) DeleteCredentials(_ context.Context, request *pbAuthx.DeleteCredentialsRequest) (*pbCommon.Success, error) {
	if request.Username == "" {
		return nil, derrors.NewOperationError("username is mandatory")
	}
	err := h.manager.DeleteCredentials(request.Username)
	if err != nil {
		return nil, err
	}
	return &pbCommon.Success{}, nil
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
	if request.Username == "" {
		return nil, derrors.NewOperationError("username is mandatory")
	}
	if request.RefreshToken == "" {
		return nil, derrors.NewOperationError("refreshToken is mandatory")
	}
	if request.TokenId == "" {
		return nil, derrors.NewOperationError("tokeID is mandatory")
	}
	return h.manager.RefreshToken(request.Username, request.TokenId, request.RefreshToken)
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
	if request.Username == "" {
		return nil, derrors.NewOperationError("username is mandatory")
	}
	if request.NewRoleId == "" {
		return nil, derrors.NewOperationError("newRoleID is mandatory")
	}
	err := h.manager.EditUserRole(request.Username, request.NewRoleId)
	if err != nil {
		return nil, err
	}
	return &pbCommon.Success{}, nil
}
