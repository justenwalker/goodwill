// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package secret

type KeyPairFiles struct {
	PublicKeyFile  string
	PrivateKeyFile string
}

type KeyPair struct {
	PublicKey  []byte
	PrivateKey []byte
}

type UsernamePassword struct {
	Username string
	Password string
}

type CreateResponse struct {
	ID            string
	StorePassword string
}

type KeyPairResponse struct {
	ID            string
	StorePassword string
	PublicKey     string
}

type Visibility int

const (
	Private Visibility = iota
	Public
)

type AccessEntry struct {
	TeamID   string
	TeamName string
	OrgName  string
	Level    AccessLevel
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
