/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package devinterceptor

import (
	"context"
	"fmt"
	"github.com/hashicorp/golang-lru"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
)

const DefaultCacheEntries = 100

// MngtSecretAccess structure that provides a cache over the group secrets using a client that connects
// directly with the Authx component. This implementation is intended to be used from within the management
// cluster.
type MngtSecretAccess struct{
	Connection
	cache lru.Cache
	Client grpc_authx_go.AuthxClient
}

func DeviceGroupIdToKey(id *grpc_device_go.DeviceGroupId) string {
	return fmt.Sprintf("%s-%s", id.OrganizationId, id.DeviceGroupId)
}

func NewMngtSecretAccess(address string, numCachedEntries int) (SecretAccess, derrors.Error) {
	lruCache, err := lru.New(numCachedEntries)
	if err != nil{
		return nil, derrors.AsError(err, "cannot create cache")
	}

	var access SecretAccess = &MngtSecretAccess{
		Connection: Connection{Address:address},
		cache:  *lruCache,
	}
	return access, nil
}

func NewMngtSecretAccessWithClient(client grpc_authx_go.AuthxClient, numCachedEntries int) (SecretAccess, derrors.Error){
	lruCache, err := lru.New(numCachedEntries)
	if err != nil{
		return nil, derrors.AsError(err, "cannot create cache")
	}

	var access SecretAccess = &MngtSecretAccess{
		Client: client,
		cache:  *lruCache,
	}
	return access, nil
}

func (sa * MngtSecretAccess) Connect() derrors.Error{
	log.Debug().Msg("connecting to authx")
	conn, err := sa.GetInsecureConnection()
	if err != nil{
		return err
	}
	sa.Client = grpc_authx_go.NewAuthxClient(conn)
	return nil
}

func (sa *MngtSecretAccess) RetrieveSecret(id *grpc_device_go.DeviceGroupId) (string, derrors.Error) {
	secret, found := sa.cache.Get(DeviceGroupIdToKey(id))
	if found {
		return secret.(string), nil
	}
	// Put it on the cache
	deviceGroupSecret, aErr := sa.Client.GetDeviceGroupSecret(context.Background(), id)
	if aErr != nil{
		log.Warn().Msg("cannot retrieve secret from authx")
		return "", conversions.ToDerror(aErr)
	}
	_ = sa.cache.Add(DeviceGroupIdToKey(id), deviceGroupSecret.Secret)
	return deviceGroupSecret.Secret, nil

}