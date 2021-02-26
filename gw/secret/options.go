// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package secret

// Options are optional values for a secret request
type Options struct {
	OrgName string
}

func (o Options) ApplyOption(t *Options) {
	*t = o
}

type Option interface {
	ApplyOption(opts *Options)

	// A private method to prevent users implementing the
	// interface and so future additions to it will not
	// violate compatibility.
	private()
}

// ExportOptions are optional parameters for export requests
type ExportOptions struct {
	OrgName       string
	StorePassword string
}

type ExportOption interface {
	ApplyExportOption(opts *ExportOptions)

	// A private method to prevent users implementing the
	// interface and so future additions to it will not
	// violate compatibility.
	private()
}

// OrgName is the name of the org containing the secret
type OrgName string

func (o OrgName) ApplyOption(spec *Options) {
	spec.OrgName = string(o)
}

func (o OrgName) ApplyExportOption(spec *ExportOptions) {
	spec.OrgName = string(o)
}

func (o OrgName) ApplyCreateOption(opts *CreateOptions) {
	opts.OrgName = string(o)
}

func (o OrgName) private() {}

// StorePassword is the password used to encrypt the secret
type StorePassword string

func (o StorePassword) ApplyExportOption(spec *ExportOptions) {
	spec.StorePassword = string(o)
}

func (o StorePassword) ApplyCreateOption(opts *CreateOptions) {
	opts.StorePassword = string(o)
}

func (o StorePassword) private() {}

type CreateOption interface {
	ApplyCreateOption(opts *CreateOptions)
}

type EncryptOptions struct {
	OrgName string
	Project string
}

type EncryptOption interface {
	ApplyEncryptOption(opts *EncryptOptions)

	// A private method to prevent users implementing the
	// interface and so future additions to it will not
	// violate compatibility.
	private()
}

type CreateOptions struct {
	OrgName          string
	StorePassword    string
	Project          string
	Visibility       Visibility
	GeneratePassword bool
}

// Project is an option to set the project name
type Project string

func (o Project) ApplyCreateOption(opts *CreateOptions) {
	opts.Project = string(o)
}

func (o Project) ApplyEncryptOption(opts *EncryptOptions) {
	opts.Project = string(o)
}

func (o Project) private() {}

func (o Visibility) ApplyCreateOption(opts *CreateOptions) {
	opts.Visibility = o
}

// GeneratePassword sets the CreateOption to generate a StorePassword
const GeneratePassword = StorePasswordGeneration(true)

// StorePasswordGeneration is a CreateOption to generate a StorePassword during secret creation
type StorePasswordGeneration bool

func (o StorePasswordGeneration) ApplyCreateOption(opts *CreateOptions) {
	opts.GeneratePassword = true
}
