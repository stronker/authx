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

type Authx struct {
	Manager *manager.Authx
}

func NewAuthx(manager *manager.Authx) *Authx {
	return &Authx{Manager: manager}
}

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

func (h *Authx) AddBasicCredentials(_ context.Context, request *pbAuthx.AddBasicCredentialRequest) (*pbCommon.Success, error) {
	if request.Username == "" {
		return nil, derrors.NewInvalidArgumentError("username is mandatory")
	}
	if request.OrganizationId == "" {
		return nil, derrors.NewInvalidArgumentError("organizationID is mandatory")
	}
	if request.RoleId == "" {
		return nil, derrors.NewInvalidArgumentError("roleID is mandatory")
	}
	if request.Password == "" {
		return nil, derrors.NewInvalidArgumentError("password is mandatory")
	}

	err := h.Manager.AddBasicCredentials(request.Username, request.OrganizationId, request.RoleId, request.Password)
	if err != nil {
		return nil, err
	}
	return &pbCommon.Success{}, nil
}

func (h *Authx) LoginWithBasicCredentials(_ context.Context, request *pbAuthx.LoginWithBasicCredentialsRequest) (*pbAuthx.LoginResponse, error) {
	if request.Username == "" {
		return nil, derrors.NewInvalidArgumentError("username is mandatory")
	}
	if request.Password == "" {
		return nil, derrors.NewInvalidArgumentError("password is mandatory")
	}
	return h.Manager.LoginWithBasicCredentials(request.Username, request.Password)
}

func (h *Authx) RefreshToken(_ context.Context, request *pbAuthx.RefreshTokenRequest) (*pbAuthx.LoginResponse, error) {
	if request.Username == "" {
		return nil, derrors.NewInvalidArgumentError("username is mandatory")
	}
	if request.RefreshToken == "" {
		return nil, derrors.NewInvalidArgumentError("refreshToken is mandatory")
	}
	if request.TokenId == "" {
		return nil, derrors.NewInvalidArgumentError("tokeID is mandatory")
	}
	return h.Manager.RefreshToken(request.Username, request.TokenId, request.RefreshToken)
}

func (h *Authx) AddRole(_ context.Context, request *pbAuthx.Role) (*pbCommon.Success, error) {
	if request.RoleId == "" {
		return nil, derrors.NewInvalidArgumentError("roleID is mandatory")
	}
	if request.Name == "" {
		return nil, derrors.NewInvalidArgumentError("name is mandatory")
	}
	if request.OrganizationId == "" {
		return nil, derrors.NewInvalidArgumentError("organizationID is mandatory")
	}
	if request.Primitives == nil || len(request.Primitives) == 0 {
		return nil, derrors.NewInvalidArgumentError("primitives is mandatory")
	}

	err := h.Manager.AddRole(request)
	if err != nil {
		return nil, err
	}
	return &pbCommon.Success{}, nil

}

func (h *Authx) EditUserRole(_ context.Context, request *pbAuthx.EditUserRoleRequest) (*pbCommon.Success, error) {
	if request.Username == "" {
		return nil, derrors.NewInvalidArgumentError("username is mandatory")
	}
	if request.NewRoleId == "" {
		return nil, derrors.NewInvalidArgumentError("newRoleID is mandatory")
	}
	err := h.Manager.EditUserRole(request.Username, request.NewRoleId)
	if err != nil {
		return nil, err
	}
	return &pbCommon.Success{}, nil
}
