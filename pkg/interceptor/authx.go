/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package interceptor

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/nalej/authx/pkg/token"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// WithServerAuthxInterceptor is a gRPC option. If this option is included, the interceptor verifies that the user is
// is authorized to use the method, using the JWT token.
func WithServerAuthxInterceptor(config *Config) grpc.ServerOption {
	return grpc.UnaryInterceptor(authxInterceptor(config))
}

func authxInterceptor(config *Config) grpc.UnaryServerInterceptor {

	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		_, ok := config.Authorization.Permissions[info.FullMethod]

		if ok {
			claim, dErr := checkJWT(ctx, config)
			if dErr != nil {
				return nil, conversions.ToGRPCError(dErr)
			}
			dErr = authorize(info.FullMethod, claim, config)

			if dErr != nil {
				return nil, conversions.ToGRPCError(dErr)
			}

		} else {
			if !config.Authorization.AllowsAll {
				return nil, conversions.ToGRPCError(
					derrors.NewUnauthenticatedError("unauthorized method").
						WithParams(info.FullMethod))
			}
		}

		return handler(ctx, req)
	}

}

func checkJWT(ctx context.Context, config *Config) (*token.Claim, derrors.Error) {
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
	tk, err := jwt.ParseWithClaims(t, &token.Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Secret), nil
	})

	if err != nil {
		return nil, derrors.NewUnauthenticatedError("token is not valid", err)
	}

	return tk.Claims.(*token.Claim), nil
}

// authorize function authorizes the token received from Metadata
func authorize(method string, claim *token.Claim, config *Config) derrors.Error {
	permission, ok := config.Authorization.Permissions[method]
	if !ok {
		if config.Authorization.AllowsAll {
			return nil
		}
		return derrors.NewUnauthenticatedError("unauthorized method").WithParams(method)
	}

	valid := permission.Valid(claim.Primitives)
	if !valid {
		return derrors.NewUnauthenticatedError("unauthorized method").WithParams(method)
	}

	return nil
}
