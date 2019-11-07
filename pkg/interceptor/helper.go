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

package interceptor

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"google.golang.org/grpc/metadata"
)

const UserIdField = "user_id"
const OrganizationIdField = "organization_id"

type RequestMetadata struct {
	UserID                 string
	OrganizationID         string
	OrgPrimitive           bool
	AppsPrimitive          bool
	ResourcePrimitive      bool
	ProfilePrimitive       bool
	AppClusterOpsPrimitive bool
}

// GetRequestMetadata extracts the request metadata from the context so that it
// can be easily consumed by upper layers.
func GetRequestMetadata(ctx context.Context) (*RequestMetadata, derrors.Error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, derrors.NewInvalidArgumentError("expecting JWT metadata")
	}
	userID, found := md[UserIdField]
	if !found {
		return nil, derrors.NewUnauthenticatedError("userID not found")
	}
	organizationID, found := md[OrganizationIdField]
	if !found {
		return nil, derrors.NewUnauthenticatedError("organizationID not found")
	}
	_, orgPrimitive := md[grpc_authx_go.AccessPrimitive_ORG.String()]
	_, appsPrimitive := md[grpc_authx_go.AccessPrimitive_APPS.String()]
	_, resourcePrimitive := md[grpc_authx_go.AccessPrimitive_RESOURCES.String()]
	_, profilePrimitive := md[grpc_authx_go.AccessPrimitive_PROFILE.String()]
	_, appClusterOpsPrimitive := md[grpc_authx_go.AccessPrimitive_APPCLUSTEROPS.String()]

	return &RequestMetadata{
		UserID:                 userID[0],
		OrganizationID:         organizationID[0],
		OrgPrimitive:           orgPrimitive,
		AppsPrimitive:          appsPrimitive,
		ResourcePrimitive:      resourcePrimitive,
		ProfilePrimitive:       profilePrimitive,
		AppClusterOpsPrimitive: appClusterOpsPrimitive,
	}, nil
}
