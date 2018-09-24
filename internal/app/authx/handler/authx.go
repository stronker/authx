/*
 * Copyright 2018 Nalej
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
