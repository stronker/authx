/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package devinterceptor

import (
	"context"
	"github.com/hashicorp/golang-lru"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-cluster-api-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-login-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

const DefaultTimeout = time.Minute
const AuthHeader = "Authorization"

// ClusterApiSecretAccess structure that provides a cache over the group secrets connecting through the cluster
// api cluster. Notice that for this type of connections, it is required to log into the management cluster
// with a set of credentials, and use the associated JWT token to send the requests.
type ClusterApiSecretAccess struct {
	LoginAPI Connection
	ClusterAPI Connection
	Username string
	Password string
	cache lru.Cache
	LoginClient grpc_login_api_go.LoginClient
	DeviceManagerClient grpc_cluster_api_go.DeviceManagerClient
	Token string
}

func (sa *ClusterApiSecretAccess) login() derrors.Error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	loginRequest := &grpc_authx_go.LoginWithBasicCredentialsRequest{
		Username: sa.Username,
		Password: sa.Password,
	}
	response, lErr := sa.LoginClient.LoginWithBasicCredentials(ctx, loginRequest)
	if lErr != nil {
		return conversions.ToDerror(lErr)
	}
	sa.Token = response.Token
	return nil
}

func (sa *ClusterApiSecretAccess) Connect() derrors.Error {
	loginConn, err := sa.LoginAPI.GetConnection()
	if err != nil{
		return err
	}
	sa.LoginClient = grpc_login_api_go.NewLoginClient(loginConn)

	cErr := sa.login()
	if cErr != nil{
		return cErr
	}

	clusterConn, err := sa.ClusterAPI.GetConnection()
	if err != nil{
		return err
	}
	sa.DeviceManagerClient = grpc_cluster_api_go.NewDeviceManagerClient(clusterConn)
 	return nil
}

func (sa * ClusterApiSecretAccess) GetContext(timeout ...time.Duration) (context.Context, context.CancelFunc) {
	md := metadata.New(map[string]string{AuthHeader: sa.Token})
	if len(timeout) == 0 {
		baseContext, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		return metadata.NewOutgoingContext(baseContext, md), cancel
	}
	baseContext, cancel := context.WithTimeout(context.Background(), timeout[0])
	return metadata.NewOutgoingContext(baseContext, md), cancel
}

func (sa *ClusterApiSecretAccess) RetrieveSecret(id *grpc_device_go.DeviceGroupId) (string, derrors.Error) {

	secret, found := sa.cache.Get(DeviceGroupIdToKey(id))
	if found {
		return secret.(string), nil
	}
	ctx, cancel := sa.GetContext()
	defer cancel()

	deviceGroupSecret, err := sa.DeviceManagerClient.GetDeviceGroupSecret(ctx, id)
	if err != nil {
		st := status.Convert(err).Code()
		if st == codes.Unauthenticated {
			errLogin := sa.login()
			if errLogin != nil {
				log.Error().Str("trace", errLogin.DebugReport()).Msg("error during reauthentication")
				return "", errLogin
			}
			ctx2, cancel2 := sa.GetContext()
			defer cancel2()
			deviceGroupSecret, err = sa.DeviceManagerClient.GetDeviceGroupSecret(ctx2, id)
		} else {
			log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error obtaining device group secret")
			return "", conversions.ToDerror(err)
		}
	}

	// Put it on the cache
	_ = sa.cache.Add(DeviceGroupIdToKey(id), deviceGroupSecret.Secret)
	return deviceGroupSecret.Secret, nil
}

