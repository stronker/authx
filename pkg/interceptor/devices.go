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
	"github.com/dgrijalva/jwt-go"
	"github.com/nalej/authx/pkg/interceptor/devinterceptor"
	"github.com/nalej/authx/pkg/token"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// WithServerDeviceAuthxInterceptor is a gRPC option. If this option is included, the interceptor verifies that the device
// is authorized to use the method, using the JWT token.
func WithDeviceAuthxInterceptor(secretAccess devinterceptor.SecretAccess, config *Config) grpc.ServerOption {
	return grpc.UnaryInterceptor(managementDeviceInterceptor(secretAccess, config))
}

// deviceInterceptor to create metadata entries for device users.
func managementDeviceInterceptor(secretAccess devinterceptor.SecretAccess, config *Config) grpc.UnaryServerInterceptor {

	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		_, ok := config.Authorization.Permissions[info.FullMethod]

		if ok {
			// Extract the raw token to be able to obtain the device group id required to retrieve the secret.
			tk, err := extractDeviceRawToken(ctx, config)
			if err != nil {
				return nil, conversions.ToGRPCError(derrors.NewUnauthenticatedError("token is not supplied"))
			}
			// Check the claim using the extracted device group id.
			claim, dErr := checkDeviceJWT(*tk, secretAccess, config)
			if dErr != nil {
				return nil, conversions.ToGRPCError(dErr)
			}
			dErr = authorizePrimitive(info.FullMethod, claim.Primitives, config)

			if dErr != nil {
				return nil, conversions.ToGRPCError(dErr)
			}
			// Add the new metadata to the context.
			values := make([]string, 0)
			values = append(values, "organization_id", claim.OrganizationID, "device_id", claim.DeviceID, "device_group_id", claim.DeviceGroupID)
			for _, p := range claim.Primitives {
				values = append(values, p, "true")
			}
			newMD := metadata.Pairs(values...)
			oldMD, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, derrors.NewInternalError("impossible to extract metadata")
			}
			newContext := metadata.NewIncomingContext(ctx, metadata.Join(oldMD, newMD))
			return handler(newContext, req)

		} else {
			if !config.Authorization.AllowsAll {
				return nil, conversions.ToGRPCError(
					derrors.NewUnauthenticatedError("unauthorized method").
						WithParams(info.FullMethod))
			}
		}
		log.Warn().Msg("auth metadata has not been added")
		return handler(ctx, req)
	}

}

// RawDeviceToken is a structure that permits to store the raw token as well as the parsed one before checking for
// the signature.
type RawDeviceToken struct {
	RawToken    string
	DeviceClaim token.DeviceClaim
}

func (rdt *RawDeviceToken) GetDeviceGroupId() *grpc_device_go.DeviceGroupId {
	return &grpc_device_go.DeviceGroupId{
		OrganizationId: rdt.DeviceClaim.OrganizationID,
		DeviceGroupId:  rdt.DeviceClaim.DeviceGroupID,
	}
}

// extractDeviceRawToken extracts the claim without checking the secret.
func extractDeviceRawToken(ctx context.Context, config *Config) (*RawDeviceToken, derrors.Error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, derrors.NewInternalError("impossible to extract metadata")
	}
	authHeader, ok := md[config.Header]
	if !ok {
		return nil, derrors.NewUnauthenticatedError("token is not supplied")
	}
	rawToken := authHeader[0]
	parser := new(jwt.Parser)
	tk, _, err := parser.ParseUnverified(rawToken, &token.DeviceClaim{})
	if err != nil {
		return nil, derrors.NewUnauthenticatedError("token is not valid", err)
	}
	return &RawDeviceToken{
		RawToken:    rawToken,
		DeviceClaim: *tk.Claims.(*token.DeviceClaim),
	}, nil
}

// CheckDeviceJWT checks the validity of the device JWT token and returns the DeviceClaim
func checkDeviceJWT(rawToken RawDeviceToken, secretAccess devinterceptor.SecretAccess, config *Config) (*token.DeviceClaim, derrors.Error) {
	secret, rErr := secretAccess.RetrieveSecret(rawToken.GetDeviceGroupId())
	if rErr != nil {
		return nil, rErr
	}
	// validateToken function validates the token
	tk, err := jwt.ParseWithClaims(rawToken.RawToken, &token.DeviceClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, derrors.NewUnauthenticatedError("token is not valid", err)
	}

	return tk.Claims.(*token.DeviceClaim), nil
}

// authorizePrimitive checks for the set of required primitives.
func authorizePrimitive(method string, primitives []string, config *Config) derrors.Error {
	permission, ok := config.Authorization.Permissions[method]
	if !ok {
		if config.Authorization.AllowsAll {
			return nil
		}
		return derrors.NewUnauthenticatedError("unauthorized method").WithParams(method)
	}

	valid := permission.Valid(primitives)
	if !valid {
		return derrors.NewUnauthenticatedError("unauthorized method").WithParams(method)
	}

	return nil
}
