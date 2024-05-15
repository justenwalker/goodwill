// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"google.golang.org/grpc"

	"go.justen.tech/goodwill/internal/pb"
)

// Service that manages concord secrets
type Service struct {
	orgName string
	client  pb.SecretServiceClient
}

// NewService creates a new service to manage secrets
func NewService(orgName string, conn grpc.ClientConnInterface) *Service {
	return &Service{
		orgName: orgName,
		client:  pb.NewSecretServiceClient(conn),
	}
}

func (c *Service) CreateSecretValue(ctx context.Context, name string, value []byte, opts ...CreateOption) (*CreateResponse, error) {
	req := &pb.CreateSecretValueRequest{
		Options: c.createSecretOptions(name, opts...),
		Value:   value,
	}
	resp, err := c.client.CreateSecretValue(ctx, req)
	if err != nil {
		return nil, err
	}
	return &CreateResponse{
		ID:            resp.Id,
		StorePassword: resp.StorePassword,
	}, nil
}

// CreateKeyPair creates a new KeyPair secret from the given Public and Private key data
func (c *Service) CreateKeyPair(ctx context.Context, name string, keypair KeyPair, opts ...CreateOption) (*KeyPairResponse, error) {
	req := &pb.CreateKeyPairRequest{
		Options:    c.createSecretOptions(name, opts...),
		PublicKey:  keypair.PublicKey,
		PrivateKey: keypair.PrivateKey,
	}
	resp, err := c.client.CreateKeyPair(ctx, req)
	if err != nil {
		return nil, err
	}
	return &KeyPairResponse{
		ID:            resp.Id,
		StorePassword: resp.StorePassword,
		PublicKey:     resp.PublicKey,
	}, nil
}

// GenerateKeyPair generates a new SSH Public/Private Key pair, returning the public key
func (c *Service) GenerateKeyPair(ctx context.Context, name string, opts ...CreateOption) (*KeyPairResponse, error) {
	req := c.createSecretOptions(name, opts...)
	resp, err := c.client.GenerateKeyPair(ctx, req)
	if err != nil {
		return nil, err
	}
	return &KeyPairResponse{
		ID:            resp.Id,
		StorePassword: resp.StorePassword,
		PublicKey:     resp.PublicKey,
	}, nil
}

// CreateUsernamePassword creates a new Username/Password secret
func (c *Service) CreateUsernamePassword(ctx context.Context, name string, usernamePassword UsernamePassword, opts ...CreateOption) (*CreateResponse, error) {
	req := &pb.CreateUsernamePasswordRequest{
		Options:  c.createSecretOptions(name, opts...),
		Username: usernamePassword.Username,
		Password: usernamePassword.Password,
	}
	resp, err := c.client.CreateUsernamePassword(ctx, req)
	if err != nil {
		return nil, err
	}
	return &CreateResponse{
		ID:            resp.Id,
		StorePassword: resp.StorePassword,
	}, nil
}

// Delete deletes a secret
func (c *Service) Delete(ctx context.Context, name string, opts ...Option) error {
	var options Options
	for _, opt := range opts {
		opt.ApplyOption(&options)
	}
	_, err := c.client.DeleteSecret(ctx, &pb.DeleteSecretRequest{
		Org:  options.OrgName,
		Name: name,
	})
	return err
}

// ListAccess returns the list of access rules for a secret
func (c *Service) ListAccess(ctx context.Context, name string, opts ...Option) ([]AccessEntry, error) {
	var results []AccessEntry
	resp, err := c.client.ListAccessLevels(ctx, c.getSecretRef(name, opts...))
	if err != nil {
		return nil, err
	}
	for _, access := range resp.Access {
		entry := AccessEntry{
			TeamID:   access.TeamID,
			TeamName: access.TeamName,
			OrgName:  access.OrgName,
		}
		switch access.Level {
		case pb.AccessEntry_READER:
			entry.Level = Reader
		case pb.AccessEntry_WRITER:
			entry.Level = Writer
		case pb.AccessEntry_OWNER:
			entry.Level = Owner
		}
		results = append(results, entry)
	}
	return results, nil
}

// UpdateAccess performs a bulk update of a secret's access rules
func (c *Service) UpdateAccess(ctx context.Context, name string, entries []AccessEntry, opts ...Option) error {
	var options Options
	for _, opt := range opts {
		opt.ApplyOption(&options)
	}
	raes := make([]*pb.AccessEntry, len(entries))
	for i, l := range entries {
		raes[i] = &pb.AccessEntry{
			TeamID:   l.TeamID,
			TeamName: l.TeamName,
			OrgName:  l.OrgName,
		}
		switch l.Level {
		case Reader:
			raes[i].Level = pb.AccessEntry_READER
		case Writer:
			raes[i].Level = pb.AccessEntry_WRITER
		case Owner:
			raes[i].Level = pb.AccessEntry_OWNER
		}
	}
	_, err := c.client.UpdateAccessLevels(ctx, &pb.UpdateSecretAccessRequest{
		OrgName: options.OrgName,
		Name:    name,
		Entries: raes,
	})
	return err
}

// KeyPairFiles exports a KeyPair to the filesystem, returning a KeyPairFiles struct
// containing the path to each of those files
func (c *Service) KeyPairFiles(ctx context.Context, name string, opts ...ExportOption) (*KeyPairFiles, error) {
	kp, err := c.client.ExportKeyPairAsFiles(ctx, c.getSecretRequest(name, opts...))
	if err != nil {
		return nil, err
	}
	return &KeyPairFiles{
		PublicKeyFile:  kp.PublicKeyFile,
		PrivateKeyFile: kp.PrivateKeyFile,
	}, nil
}

// UsernamePassword exports a Username+Password secret
func (c *Service) UsernamePassword(ctx context.Context, name string, opts ...ExportOption) (*UsernamePassword, error) {
	up, err := c.client.GetUsernamePassword(ctx, c.getSecretRequest(name, opts...))
	if err != nil {
		return nil, err
	}
	return &UsernamePassword{
		Username: up.Username,
		Password: up.Password,
	}, nil
}

// SecretString exports a single-value data secret as a regular string
func (c *Service) SecretString(ctx context.Context, name string, opts ...ExportOption) (string, error) {
	sf, err := c.client.ExportAsString(ctx, c.getSecretRequest(name, opts...))
	if err != nil {
		return "", err
	}
	return sf.Str, nil
}

// SecretFile exports a single-value data secret into a file and returns the path to that file
func (c *Service) SecretFile(ctx context.Context, name string, opts ...ExportOption) (string, error) {
	sf, err := c.client.ExportAsFile(ctx, c.getSecretRequest(name, opts...))
	if err != nil {
		return "", err
	}
	return sf.File, nil
}

// EncryptString encrypts a string which can be later decrypted by Concord
func (c *Service) EncryptString(ctx context.Context, value string, opts ...EncryptOption) (string, error) {
	var options EncryptOptions
	for _, opt := range opts {
		opt.ApplyEncryptOption(&options)
	}
	enc, err := c.client.EncryptString(ctx, &pb.EncryptStringRequest{
		Org:     options.OrgName,
		Project: options.Project,
		Value:   value,
	})
	if err != nil {
		return "", err
	}
	return enc.Str, nil
}

// DecryptString takes a string previously encrypted by EncryptString and returns the decrypted value
func (c *Service) DecryptString(ctx context.Context, value string) (string, error) {
	dec, err := c.client.DecryptString(ctx, &pb.SecretString{
		Str: value,
	})
	if err != nil {
		return "", err
	}
	return dec.Str, nil
}

func (c *Service) getSecretRef(name string, opts ...Option) *pb.SecretRef {
	var options Options
	options.OrgName = c.orgName
	for _, opt := range opts {
		opt.ApplyOption(&options)
	}
	return &pb.SecretRef{
		OrgName: options.OrgName,
		Name:    name,
	}
}

func (c *Service) getSecretRequest(name string, opts ...ExportOption) *pb.GetSecretRequest {
	var options ExportOptions
	options.OrgName = c.orgName
	for _, opt := range opts {
		opt.ApplyExportOption(&options)
	}
	return &pb.GetSecretRequest{
		Org:           options.OrgName,
		Name:          name,
		StorePassword: options.StorePassword,
	}
}

func (c *Service) createSecretOptions(name string, opts ...CreateOption) *pb.SecretParams {
	var req pb.SecretParams
	req.Name = name
	var options CreateOptions
	options.OrgName = c.orgName
	for _, opt := range opts {
		opt.ApplyCreateOption(&options)
	}
	req.OrgName = options.OrgName
	req.StorePassword = options.StorePassword
	req.Project = options.Project
	req.GeneratePassword = options.GeneratePassword
	switch options.Visibility {
	case Private:
		req.Visibility = pb.SecretParams_PRIVATE
	case Public:
		req.Visibility = pb.SecretParams_PUBLIC
	}
	return &req
}
