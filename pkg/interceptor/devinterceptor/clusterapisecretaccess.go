/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package devinterceptor


import (
	"github.com/hashicorp/golang-lru"
	"github.com/nalej/grpc-cluster-api-go"
)

type ClusterApiSecretAccess struct {
	cache lru.Cache
	Username string
	Password string
	client grpc_cluster_api_go.DeviceManagerClient
}