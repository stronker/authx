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
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

func WithServerUnaryInterceptor(config Config) grpc.ServerOption {
	return grpc.UnaryInterceptor(authxInterceptor(config))
}

func authxInterceptor(config Config) grpc.UnaryServerInterceptor {

	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		start := time.Now()

		claim, dErr := checkJWT(ctx, config)
		if dErr != nil {
			return nil, conversions.ToGRPCError(dErr)
		}
		dErr = authorize(info.FullMethod, *claim, config)

		if dErr != nil {
			return nil, conversions.ToGRPCError(dErr)
		}

		// Calls the handler
		h, err := handler(ctx, req)

		log.Info().Msgf("Request - Method:%s\tDuration:%s\tError:%v\n",
			info.FullMethod,
			time.Since(start),
			err)

		return h, err
	}

}

func checkJWT(ctx context.Context, config Config) (*token.Claim, derrors.Error) {
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
func authorize(method string, claim token.Claim, config Config) derrors.Error {
	permission, ok := config.Authorization.Permissions[method]
	if !ok {
		if config.AllowsAll {
			return nil
		}
		return derrors.NewUnauthenticatedError("unauthorized method").WithParams(method)
	}

	valid := permission.Check(claim.Primitives)
	if !valid {
		return derrors.NewUnauthenticatedError("unauthorized method").WithParams(method)
	}

	return nil
}
