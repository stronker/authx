/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package devinterceptor

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"strconv"
	"strings"
)

// Connection is a helper structure to establish the gRPC connection.
type Connection struct {
	Address    string
	UseTLS bool
	CACertPath string
	SkipCAValidation bool
}

func (c*Connection) SplitAddress() (string, int, derrors.Error){
	splits := strings.Split(c.Address, ":")
	if len(splits) != 2{
		return "", -1, derrors.NewInvalidArgumentError("expecting address as hostname:port")
	}
	port, err := strconv.Atoi(splits[1])
	if err != nil{
		return "", -1, derrors.AsError(err, "expecting address as hostname:port")
	}
	return splits[0], port, nil
}

func (c *Connection) GetInsecureConnection() (*grpc.ClientConn, derrors.Error) {
	return c.GetInsecureGRPCConnection(c.Address)
}

func (c *Connection) GetInsecureGRPCConnection(address string) (*grpc.ClientConn, derrors.Error) {
	log.Debug().Str("address", address).Msg("creating connection")
	conn, err := grpc.Dial(c.Address, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the auth proxy")
	}
	return conn, nil
}

func (c* Connection) GetSecureConnection() (*grpc.ClientConn, derrors.Error){
	targetHostname, _, err := c.SplitAddress()
	if err !=nil{
		return nil, err
	}
	rootCAs := x509.NewCertPool()
	tlsConfig := &tls.Config{
		ServerName:   targetHostname,
	}

	if c.CACertPath != "" {
		log.Debug().Str("caCertPath", c.CACertPath).Msg("loading CA cert")
		caCert, err := ioutil.ReadFile(c.CACertPath)
		if err != nil {
			return nil, derrors.NewInternalError("Error loading CA certificate")
		}
		added := rootCAs.AppendCertsFromPEM(caCert)
		if !added {
			return nil, derrors.NewInternalError("cannot add CA certificate to the pool")
		}
		tlsConfig.RootCAs = rootCAs
	}

	if c.SkipCAValidation {
		log.Warn().Msg("CA will not be verified")
		tlsConfig.InsecureSkipVerify = true
	}

	creds := credentials.NewTLS(tlsConfig)
	log.Debug().Interface("creds", creds.Info()).Msg("Secure credentials")
	log.Debug().Str("address", c.Address).Msg("creating connection")
	sConn, dErr := grpc.Dial(c.Address, grpc.WithTransportCredentials(creds))
	if dErr != nil {
		return nil, derrors.AsError(dErr, "cannot create connection with the signup service")
	}
	return sConn, nil
}

func (c *Connection) GetConnection() (*grpc.ClientConn, derrors.Error) {
	if c.UseTLS {
		return c.GetSecureConnection()
	}
	return c.GetInsecureConnection()
}