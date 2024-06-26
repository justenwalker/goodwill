// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: docker.proto

package pb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DockerContainerResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Result:
	//
	//	*DockerContainerResult_Status
	//	*DockerContainerResult_Stdout
	//	*DockerContainerResult_Stderr
	Result isDockerContainerResult_Result `protobuf_oneof:"result"`
}

func (x *DockerContainerResult) Reset() {
	*x = DockerContainerResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_docker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DockerContainerResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DockerContainerResult) ProtoMessage() {}

func (x *DockerContainerResult) ProtoReflect() protoreflect.Message {
	mi := &file_docker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DockerContainerResult.ProtoReflect.Descriptor instead.
func (*DockerContainerResult) Descriptor() ([]byte, []int) {
	return file_docker_proto_rawDescGZIP(), []int{0}
}

func (m *DockerContainerResult) GetResult() isDockerContainerResult_Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (x *DockerContainerResult) GetStatus() uint32 {
	if x, ok := x.GetResult().(*DockerContainerResult_Status); ok {
		return x.Status
	}
	return 0
}

func (x *DockerContainerResult) GetStdout() string {
	if x, ok := x.GetResult().(*DockerContainerResult_Stdout); ok {
		return x.Stdout
	}
	return ""
}

func (x *DockerContainerResult) GetStderr() string {
	if x, ok := x.GetResult().(*DockerContainerResult_Stderr); ok {
		return x.Stderr
	}
	return ""
}

type isDockerContainerResult_Result interface {
	isDockerContainerResult_Result()
}

type DockerContainerResult_Status struct {
	Status uint32 `protobuf:"varint,1,opt,name=status,proto3,oneof"`
}

type DockerContainerResult_Stdout struct {
	Stdout string `protobuf:"bytes,2,opt,name=stdout,proto3,oneof"`
}

type DockerContainerResult_Stderr struct {
	Stderr string `protobuf:"bytes,3,opt,name=stderr,proto3,oneof"`
}

func (*DockerContainerResult_Status) isDockerContainerResult_Result() {}

func (*DockerContainerResult_Stdout) isDockerContainerResult_Result() {}

func (*DockerContainerResult_Stderr) isDockerContainerResult_Result() {}

type DockerContainerSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Image            string            `protobuf:"bytes,1,opt,name=image,proto3" json:"image,omitempty"`
	Name             string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	User             string            `protobuf:"bytes,3,opt,name=user,proto3" json:"user,omitempty"`
	WorkDir          string            `protobuf:"bytes,4,opt,name=workDir,proto3" json:"workDir,omitempty"`
	EntryPoint       string            `protobuf:"bytes,5,opt,name=entryPoint,proto3" json:"entryPoint,omitempty"`
	Command          []string          `protobuf:"bytes,6,rep,name=command,proto3" json:"command,omitempty"`
	Cpu              string            `protobuf:"bytes,7,opt,name=cpu,proto3" json:"cpu,omitempty"`
	Memory           string            `protobuf:"bytes,8,opt,name=memory,proto3" json:"memory,omitempty"`
	Env              map[string]string `protobuf:"bytes,9,rep,name=env,proto3" json:"env,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	EnvFile          string            `protobuf:"bytes,10,opt,name=envFile,proto3" json:"envFile,omitempty"`
	Labels           map[string]string `protobuf:"bytes,11,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ForcePull        bool              `protobuf:"varint,12,opt,name=forcePull,proto3" json:"forcePull,omitempty"`
	Hosts            []string          `protobuf:"bytes,13,rep,name=hosts,proto3" json:"hosts,omitempty"`
	StdoutFilePath   string            `protobuf:"bytes,14,opt,name=stdoutFilePath,proto3" json:"stdoutFilePath,omitempty"`
	RedirectStdError bool              `protobuf:"varint,15,opt,name=redirectStdError,proto3" json:"redirectStdError,omitempty"`
}

func (x *DockerContainerSpec) Reset() {
	*x = DockerContainerSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_docker_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DockerContainerSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DockerContainerSpec) ProtoMessage() {}

func (x *DockerContainerSpec) ProtoReflect() protoreflect.Message {
	mi := &file_docker_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DockerContainerSpec.ProtoReflect.Descriptor instead.
func (*DockerContainerSpec) Descriptor() ([]byte, []int) {
	return file_docker_proto_rawDescGZIP(), []int{1}
}

func (x *DockerContainerSpec) GetImage() string {
	if x != nil {
		return x.Image
	}
	return ""
}

func (x *DockerContainerSpec) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DockerContainerSpec) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *DockerContainerSpec) GetWorkDir() string {
	if x != nil {
		return x.WorkDir
	}
	return ""
}

func (x *DockerContainerSpec) GetEntryPoint() string {
	if x != nil {
		return x.EntryPoint
	}
	return ""
}

func (x *DockerContainerSpec) GetCommand() []string {
	if x != nil {
		return x.Command
	}
	return nil
}

func (x *DockerContainerSpec) GetCpu() string {
	if x != nil {
		return x.Cpu
	}
	return ""
}

func (x *DockerContainerSpec) GetMemory() string {
	if x != nil {
		return x.Memory
	}
	return ""
}

func (x *DockerContainerSpec) GetEnv() map[string]string {
	if x != nil {
		return x.Env
	}
	return nil
}

func (x *DockerContainerSpec) GetEnvFile() string {
	if x != nil {
		return x.EnvFile
	}
	return ""
}

func (x *DockerContainerSpec) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *DockerContainerSpec) GetForcePull() bool {
	if x != nil {
		return x.ForcePull
	}
	return false
}

func (x *DockerContainerSpec) GetHosts() []string {
	if x != nil {
		return x.Hosts
	}
	return nil
}

func (x *DockerContainerSpec) GetStdoutFilePath() string {
	if x != nil {
		return x.StdoutFilePath
	}
	return ""
}

func (x *DockerContainerSpec) GetRedirectStdError() bool {
	if x != nil {
		return x.RedirectStdError
	}
	return false
}

var File_docker_proto protoreflect.FileDescriptor

var file_docker_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x22, 0x6f, 0x0a, 0x15, 0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72,
	0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12,
	0x18, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x48,
	0x00, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x06, 0x73, 0x74, 0x64,
	0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x06, 0x73, 0x74, 0x64,
	0x6f, 0x75, 0x74, 0x12, 0x18, 0x0a, 0x06, 0x73, 0x74, 0x64, 0x65, 0x72, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x06, 0x73, 0x74, 0x64, 0x65, 0x72, 0x72, 0x42, 0x08, 0x0a,
	0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0xdf, 0x04, 0x0a, 0x13, 0x44, 0x6f, 0x63, 0x6b,
	0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x12,
	0x14, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x18, 0x0a,
	0x07, 0x77, 0x6f, 0x72, 0x6b, 0x44, 0x69, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x77, 0x6f, 0x72, 0x6b, 0x44, 0x69, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x6e, 0x74, 0x72, 0x79,
	0x50, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x6e, 0x74,
	0x72, 0x79, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x70, 0x75, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x63, 0x70, 0x75, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x12, 0x36, 0x0a, 0x03, 0x65,
	0x6e, 0x76, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x64, 0x6f, 0x63, 0x6b, 0x65,
	0x72, 0x2e, 0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x72, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x45, 0x6e, 0x76, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x03,
	0x65, 0x6e, 0x76, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x76, 0x46, 0x69, 0x6c, 0x65, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x76, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x3f, 0x0a,
	0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e,
	0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x43, 0x6f, 0x6e,
	0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x1c,
	0x0a, 0x09, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x50, 0x75, 0x6c, 0x6c, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x09, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x50, 0x75, 0x6c, 0x6c, 0x12, 0x14, 0x0a, 0x05,
	0x68, 0x6f, 0x73, 0x74, 0x73, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x68, 0x6f, 0x73,
	0x74, 0x73, 0x12, 0x26, 0x0a, 0x0e, 0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65,
	0x50, 0x61, 0x74, 0x68, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x74, 0x64, 0x6f,
	0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x2a, 0x0a, 0x10, 0x72, 0x65,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x53, 0x74, 0x64, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x0f,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x10, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x53, 0x74,
	0x64, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x1a, 0x36, 0x0a, 0x08, 0x45, 0x6e, 0x76, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x39,
	0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x5f, 0x0a, 0x0d, 0x44, 0x6f, 0x63,
	0x6b, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4e, 0x0a, 0x0c, 0x52, 0x75,
	0x6e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x12, 0x1b, 0x2e, 0x64, 0x6f, 0x63,
	0x6b, 0x65, 0x72, 0x2e, 0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69,
	0x6e, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x1a, 0x1d, 0x2e, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72,
	0x2e, 0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x30, 0x01, 0x42, 0x55, 0x0a, 0x21, 0x74, 0x65,
	0x63, 0x68, 0x2e, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x6e, 0x2e, 0x63, 0x6f, 0x6e, 0x63, 0x6f, 0x72,
	0x64, 0x2e, 0x67, 0x6f, 0x6f, 0x64, 0x77, 0x69, 0x6c, 0x6c, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x42,
	0x0b, 0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5a, 0x23, 0x67, 0x6f,
	0x2e, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x6e, 0x2e, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x67, 0x6f, 0x6f,
	0x64, 0x77, 0x69, 0x6c, 0x6c, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_docker_proto_rawDescOnce sync.Once
	file_docker_proto_rawDescData = file_docker_proto_rawDesc
)

func file_docker_proto_rawDescGZIP() []byte {
	file_docker_proto_rawDescOnce.Do(func() {
		file_docker_proto_rawDescData = protoimpl.X.CompressGZIP(file_docker_proto_rawDescData)
	})
	return file_docker_proto_rawDescData
}

var file_docker_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_docker_proto_goTypes = []interface{}{
	(*DockerContainerResult)(nil), // 0: docker.DockerContainerResult
	(*DockerContainerSpec)(nil),   // 1: docker.DockerContainerSpec
	nil,                           // 2: docker.DockerContainerSpec.EnvEntry
	nil,                           // 3: docker.DockerContainerSpec.LabelsEntry
}
var file_docker_proto_depIdxs = []int32{
	2, // 0: docker.DockerContainerSpec.env:type_name -> docker.DockerContainerSpec.EnvEntry
	3, // 1: docker.DockerContainerSpec.labels:type_name -> docker.DockerContainerSpec.LabelsEntry
	1, // 2: docker.DockerService.RunContainer:input_type -> docker.DockerContainerSpec
	0, // 3: docker.DockerService.RunContainer:output_type -> docker.DockerContainerResult
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_docker_proto_init() }
func file_docker_proto_init() {
	if File_docker_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_docker_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DockerContainerResult); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_docker_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DockerContainerSpec); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_docker_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*DockerContainerResult_Status)(nil),
		(*DockerContainerResult_Stdout)(nil),
		(*DockerContainerResult_Stderr)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_docker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_docker_proto_goTypes,
		DependencyIndexes: file_docker_proto_depIdxs,
		MessageInfos:      file_docker_proto_msgTypes,
	}.Build()
	File_docker_proto = out.File
	file_docker_proto_rawDesc = nil
	file_docker_proto_goTypes = nil
	file_docker_proto_depIdxs = nil
}
