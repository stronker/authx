/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package handler

import (
	"context"
	pbAuthx "github.com/nalej/grpc-authx-go"
	pbCommon "github.com/nalej/grpc-common-go"
	pbUser "github.com/nalej/grpc-user-go"
)

type AuthxServer struct {
}

func NewAuthxServer() *AuthxServer {
	return &AuthxServer{}
}

func (*AuthxServer) DeleteCredentials(context.Context, *pbUser.UserId) (*pbCommon.Success, error) {
	panic("implement me")
}

func (*AuthxServer) AddBasicCredentials(context.Context, *pbAuthx.AddBasicCredentialRequest) (*pbCommon.Success, error) {
	panic("implement me")
}

func (*AuthxServer) LoginWithBasicCredentials(context.Context, *pbAuthx.LoginWithBasicCredentialsRequest) (*pbAuthx.LoginResponse, error) {
	panic("implement me")
}

func (*AuthxServer) RefreshToken(context.Context, *pbAuthx.RefreshTokenRequest) (*pbAuthx.LoginResponse, error) {
	panic("implement me")
}

func (*AuthxServer) AddRole(context.Context, *pbAuthx.Role) (*pbCommon.Success, error) {
	panic("implement me")
}

func (*AuthxServer) EditUserRole(context.Context, *pbAuthx.EditUserRoleRequest) (*pbCommon.Success, error) {
	panic("implement me")
}
