/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package interceptor

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/nalej/authx/pkg/token"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"github.com/hashicorp/golang-lru"
)

// TODO, two interceptors are required, a local one with grpc_authx_go Clients and another with grpc_app_cluster...

// WithServerDeviceAuthxInterceptor is a gRPC option. If this option is included, the interceptor verifies that the device
// is authorized to use the method, using the JWT token.
func WithServerDeviceAuthxInterceptor(client grpc_authx_go.AuthxClient, config *Config) grpc.ServerOption {
	return grpc.UnaryInterceptor(deviceInterceptor(client, config))
}

// deviceInterceptor to create metadata entries for device users.
func deviceInterceptor(client grpc_authx_go.AuthxClient, config *Config) grpc.UnaryServerInterceptor {

	groupSecretCache, err := lru.New(config.NumCacheEntries)
	if err != nil{
		log.Fatal().Err(err).Msg("cannot create LRU cache for devicegroup secrets")
	}

	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		_, ok := config.Authorization.Permissions[info.FullMethod]

		if ok {
			//deviceGroupId := extractDeviceGroupId(ctx)
			claim, dErr := checkDeviceJWT(ctx, groupSecretCache, config)
			if dErr != nil {
				return nil, conversions.ToGRPCError(dErr)
			}
			dErr = authorizePrimitive(info.FullMethod, claim.Primitives, config)

			if dErr != nil {
				return nil, conversions.ToGRPCError(dErr)
			}

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

type RawDeviceToken struct {
	RawToken string
	DeviceClaim token.DeviceClaim
}

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
	if err != nil{
		return nil, derrors.NewUnauthenticatedError("token is not valid", err)
	}
	return &RawDeviceToken{
		RawToken:      rawToken,
		DeviceClaim:   *tk.Claims.(*token.DeviceClaim),
	}, nil
}

// CheckDeviceJWT checks the validity of the device JWT token and returns the DeviceClaim
func checkDeviceJWT(ctx context.Context, cache *lru.Cache, config *Config) (*token.DeviceClaim, derrors.Error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, derrors.NewInternalError("impossible to extract metadata")
	}

	authHeader, ok := md[config.Header]
	if !ok {
		return nil, derrors.NewUnauthenticatedError("token is not supplied")
	}
	t := authHeader[0]
	// validateToken function validates the token
	tk, err := jwt.ParseWithClaims(t, &token.DeviceClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Secret), nil
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