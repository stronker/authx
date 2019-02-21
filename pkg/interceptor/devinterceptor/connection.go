/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package devinterceptor

import (
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// Connection is a helper structure to establish the gRPC connection.
type Connection struct {
	Address    string
}

func (c *Connection) GetInsecureConnection() (*grpc.ClientConn, derrors.Error) {
	log.Debug().Str("address", c.Address).Msg("creating connection")
	conn, err := grpc.Dial(c.Address, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the auth proxy")
	}
	return conn, nil
}