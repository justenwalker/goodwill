// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package secret

// KeyPairFiles are a pair a file names for a keypair
type KeyPairFiles struct {
	// PublicKeyFile is the path to the public key
	PublicKeyFile string
	// PrivateKeyFile is the path to the private key
	PrivateKeyFile string
}

// KeyPair contains key pair data
type KeyPair struct {
	// PublicKey is the public key data
	// Usually this is an OpenSSH public key: ssh-<alg> KEY-DATA COMMENT
	PublicKey []byte
	// PrivateKey is the private key data.
	// Usually this is a PEM encoded OPENSSH PRIVATE KEY
	PrivateKey []byte
}

// UsernamePassword is a username/password pair
type UsernamePassword struct {
	// Username is the username
	Username string
	// Password is the password
	Password string
}

// CreateResponse is the response of a secret creation request
type CreateResponse struct {
	// ID is the unique ID of the new secret
	ID string
	// StorePassword is the StorePassword of the new secret, if one was generated by Concord
	StorePassword string
}

// KeyPairResponse  is the response of key-pair generation request
type KeyPairResponse struct {
	// ID is the unique ID of the new secret
	ID string
	// StorePassword is the StorePassword of the new secret, if one was generated by Concord
	StorePassword string
	// PublicKey is the public key data
	// Usually this is an OpenSSH public key: ssh-<alg> KEY-DATA COMMENT
	PublicKey string
}

// Visibility represents the secrets visibility in Concord
type Visibility int

const (
	Private Visibility = iota
	Public
)

// AccessEntry is an entry in the secret's ACL table
type AccessEntry struct {
	// TeamID is the unique ID of the team being granted permission
	TeamID string
	// TeamName is the name of the team being granted permissions.
	// If this is provided, OrgName must also be provided.
	// If TeamID is set, this does nothing.
	TeamName string
	// OrgName is the name of the org containing the TeamName to be granted access
	// Does nothing if TeamID is set
	OrgName string
	// Lavel is the level of access being granted. See AccessLevel for the types of access.
	Level AccessLevel
}

type AccessLevel int

const (
	Reader AccessLevel = iota
	Writer
	Owner
)

func (o AccessLevel) String() string {
	switch o {
	case Reader:
		return "Reader"
	case Writer:
		return "Writer"
	case Owner:
		return "Owner"
	}
	return "Unknown"
}
