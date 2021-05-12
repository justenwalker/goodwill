// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package gw

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"os"
)

const (
	EnvGRPCAddress = "GRPC_ADDR"
	EnvMagicKey    = "GRPC_MAGIC_KEY"
	EnvMagicValue  = "d0c08ee0-a663-4a6b-ad5e-00a5fca1e5cf"
	EnvCAFile      = "GRPC_CA_CERT_FILE"
	EnvCertFile    = "GRPC_CLIENT_CERT_FILE"
	EnvKeyFile     = "GRPC_CLIENT_KEY_FILE"
	EnvOrgName     = "CONCORD_ORG_NAME"
	EnvProcessID   = "CONCORD_PROCESS_ID"
	EnvWorkingDir  = "CONCORD_WORKING_DIRECTORY"
)

const (
	ErrNoMagicKey = Error("no magic key provided")
	ErrNoGRPCAddr = Error("no grpc address provided")
)

var ranOnce bool

// Run a function
func Run(ctx context.Context, tr TaskRunner) error {
	// Guard to prevent task from executing this function
	if ranOnce {
		return Error("Run called more than once")
	}
	ranOnce = true
	var err error
	if os.Getenv(EnvMagicKey) != EnvMagicValue {
		return ErrNoMagicKey
	}
	addr := os.Getenv(EnvGRPCAddress)
	if addr == "" {
		return ErrNoGRPCAddr
	}
	orgName := DefaultOrgName
	if name := os.Getenv(EnvOrgName); name != "" {
		orgName = name
	}
	var opts []grpc.DialOption
	opts, err = transportSecurity(opts)
	if err != nil {
		return fmt.Errorf("fail to configure TLS: %w", err)
	}
	c, err := newTask(addr, runtime{
		OrgName:    orgName,
		WorkingDir: os.Getenv(EnvWorkingDir),
		ProcessID:  os.Getenv(EnvProcessID),
	}, opts...)
	if err != nil {
		return fmt.Errorf("could not create context: %w", err)
	}
	defer c.Close()
	vars, err := tr.Run(ctx, c)
	if err != nil {
		return fmt.Errorf("task failed: %w", err)
	}
	err = c.Context().SetVariables(ctx, vars)
	if err != nil {
		return fmt.Errorf("failed to set output variables: %w", err)
	}
	return nil
}

func transportSecurity(opts []grpc.DialOption) ([]grpc.DialOption, error) {
	caCertFile := os.Getenv(EnvCAFile)
	certFile := os.Getenv(EnvCertFile)
	keyFile := os.Getenv(EnvKeyFile)
	if caCertFile == "" {
		return append(opts, grpc.WithInsecure()), nil
	}
	caBytes, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("could not read %q: %w", caCertFile, err)
	}
	caPem, _ := pem.Decode(caBytes)
	if caPem.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("ca is not a CERTIFICATE: was %s", caPem.Type)
	}
	caCert, err := x509.ParseCertificate(caPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("could parse ca cert %q: %w", caCertFile, err)
	}
	pool := x509.NewCertPool()
	pool.AddCert(caCert)
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("could not load client keypair (%s,%s): %w", certFile, keyFile, err)
	}
	tp := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		ServerName:   "localhost",
	})

	return append(opts, grpc.WithTransportCredentials(tp)), nil
}
