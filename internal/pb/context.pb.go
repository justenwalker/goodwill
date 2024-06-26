// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: context.proto

package pb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetVariablesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetVariablesRequest) Reset() {
	*x = GetVariablesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVariablesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVariablesRequest) ProtoMessage() {}

func (x *GetVariablesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVariablesRequest.ProtoReflect.Descriptor instead.
func (*GetVariablesRequest) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{0}
}

type SetVariableResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SetVariableResult) Reset() {
	*x = SetVariableResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetVariableResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetVariableResult) ProtoMessage() {}

func (x *SetVariableResult) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetVariableResult.ProtoReflect.Descriptor instead.
func (*SetVariableResult) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{1}
}

type GetVariableNameParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetVariableNameParams) Reset() {
	*x = GetVariableNameParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVariableNameParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVariableNameParams) ProtoMessage() {}

func (x *GetVariableNameParams) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVariableNameParams.ProtoReflect.Descriptor instead.
func (*GetVariableNameParams) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{2}
}

type Variables struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Parameters []*Variable `protobuf:"bytes,1,rep,name=parameters,proto3" json:"parameters,omitempty"`
}

func (x *Variables) Reset() {
	*x = Variables{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Variables) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Variables) ProtoMessage() {}

func (x *Variables) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Variables.ProtoReflect.Descriptor instead.
func (*Variables) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{3}
}

func (x *Variables) GetParameters() []*Variable {
	if x != nil {
		return x.Parameters
	}
	return nil
}

type Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value *anypb.Any `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Value) Reset() {
	*x = Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{4}
}

func (x *Value) GetValue() *anypb.Any {
	if x != nil {
		return x.Value
	}
	return nil
}

type MapValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value map[string]*Value `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *MapValue) Reset() {
	*x = MapValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapValue) ProtoMessage() {}

func (x *MapValue) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapValue.ProtoReflect.Descriptor instead.
func (*MapValue) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{5}
}

func (x *MapValue) GetValue() map[string]*Value {
	if x != nil {
		return x.Value
	}
	return nil
}

type ListValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []*Value `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
}

func (x *ListValue) Reset() {
	*x = ListValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListValue) ProtoMessage() {}

func (x *ListValue) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListValue.ProtoReflect.Descriptor instead.
func (*ListValue) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{6}
}

func (x *ListValue) GetValue() []*Value {
	if x != nil {
		return x.Value
	}
	return nil
}

type StringValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *StringValue) Reset() {
	*x = StringValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringValue) ProtoMessage() {}

func (x *StringValue) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringValue.ProtoReflect.Descriptor instead.
func (*StringValue) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{7}
}

func (x *StringValue) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type IntValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value int64 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *IntValue) Reset() {
	*x = IntValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IntValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntValue) ProtoMessage() {}

func (x *IntValue) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntValue.ProtoReflect.Descriptor instead.
func (*IntValue) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{8}
}

func (x *IntValue) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type DoubleValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value float64 `protobuf:"fixed64,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *DoubleValue) Reset() {
	*x = DoubleValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DoubleValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DoubleValue) ProtoMessage() {}

func (x *DoubleValue) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DoubleValue.ProtoReflect.Descriptor instead.
func (*DoubleValue) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{9}
}

func (x *DoubleValue) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type BoolValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value bool `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *BoolValue) Reset() {
	*x = BoolValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoolValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoolValue) ProtoMessage() {}

func (x *BoolValue) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoolValue.ProtoReflect.Descriptor instead.
func (*BoolValue) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{10}
}

func (x *BoolValue) GetValue() bool {
	if x != nil {
		return x.Value
	}
	return false
}

type TimeValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value int64 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *TimeValue) Reset() {
	*x = TimeValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeValue) ProtoMessage() {}

func (x *TimeValue) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeValue.ProtoReflect.Descriptor instead.
func (*TimeValue) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{11}
}

func (x *TimeValue) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type NullValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NullValue) Reset() {
	*x = NullValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NullValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NullValue) ProtoMessage() {}

func (x *NullValue) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NullValue.ProtoReflect.Descriptor instead.
func (*NullValue) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{12}
}

type JSONValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Class string `protobuf:"bytes,1,opt,name=class,proto3" json:"class,omitempty"`
	Json  []byte `protobuf:"bytes,2,opt,name=json,proto3" json:"json,omitempty"`
}

func (x *JSONValue) Reset() {
	*x = JSONValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JSONValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JSONValue) ProtoMessage() {}

func (x *JSONValue) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JSONValue.ProtoReflect.Descriptor instead.
func (*JSONValue) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{13}
}

func (x *JSONValue) GetClass() string {
	if x != nil {
		return x.Class
	}
	return ""
}

func (x *JSONValue) GetJson() []byte {
	if x != nil {
		return x.Json
	}
	return nil
}

type EvaluateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Expression string      `protobuf:"bytes,1,opt,name=expression,proto3" json:"expression,omitempty"`
	Type       string      `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Parameters []*Variable `protobuf:"bytes,3,rep,name=parameters,proto3" json:"parameters,omitempty"`
}

func (x *EvaluateRequest) Reset() {
	*x = EvaluateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EvaluateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EvaluateRequest) ProtoMessage() {}

func (x *EvaluateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EvaluateRequest.ProtoReflect.Descriptor instead.
func (*EvaluateRequest) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{14}
}

func (x *EvaluateRequest) GetExpression() string {
	if x != nil {
		return x.Expression
	}
	return ""
}

func (x *EvaluateRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *EvaluateRequest) GetParameters() []*Variable {
	if x != nil {
		return x.Parameters
	}
	return nil
}

type VariableName struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *VariableName) Reset() {
	*x = VariableName{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[15]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VariableName) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VariableName) ProtoMessage() {}

func (x *VariableName) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[15]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VariableName.ProtoReflect.Descriptor instead.
func (*VariableName) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{15}
}

func (x *VariableName) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type VariableNameList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name []string `protobuf:"bytes,1,rep,name=name,proto3" json:"name,omitempty"`
}

func (x *VariableNameList) Reset() {
	*x = VariableNameList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[16]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VariableNameList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VariableNameList) ProtoMessage() {}

func (x *VariableNameList) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[16]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VariableNameList.ProtoReflect.Descriptor instead.
func (*VariableNameList) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{16}
}

func (x *VariableNameList) GetName() []string {
	if x != nil {
		return x.Name
	}
	return nil
}

type Variable struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value *Value `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Variable) Reset() {
	*x = Variable{}
	if protoimpl.UnsafeEnabled {
		mi := &file_context_proto_msgTypes[17]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Variable) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Variable) ProtoMessage() {}

func (x *Variable) ProtoReflect() protoreflect.Message {
	mi := &file_context_proto_msgTypes[17]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Variable.ProtoReflect.Descriptor instead.
func (*Variable) Descriptor() ([]byte, []int) {
	return file_context_proto_rawDescGZIP(), []int{17}
}

func (x *Variable) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Variable) GetValue() *Value {
	if x != nil {
		return x.Value
	}
	return nil
}

var File_context_proto protoreflect.FileDescriptor

var file_context_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x15, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62,
	0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x13, 0x0a, 0x11, 0x53, 0x65,
	0x74, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22,
	0x17, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61,
	0x6d, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x22, 0x3e, 0x0a, 0x09, 0x56, 0x61, 0x72, 0x69,
	0x61, 0x62, 0x6c, 0x65, 0x73, 0x12, 0x31, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74,
	0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x78, 0x74, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x0a, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x22, 0x33, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x2a, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x88, 0x01,
	0x0a, 0x08, 0x4d, 0x61, 0x70, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x32, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x78, 0x74, 0x2e, 0x4d, 0x61, 0x70, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x48,
	0x0a, 0x0a, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x24,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x31, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x24, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x23, 0x0a, 0x0b, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x20, 0x0a, 0x08, 0x49, 0x6e, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x23, 0x0a, 0x0b, 0x44, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x21, 0x0a, 0x09, 0x42, 0x6f, 0x6f, 0x6c, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x21, 0x0a, 0x09, 0x54, 0x69,
	0x6d, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x0b, 0x0a,
	0x09, 0x4e, 0x75, 0x6c, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x35, 0x0a, 0x09, 0x4a, 0x53,
	0x4f, 0x4e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6c, 0x61, 0x73, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x6a, 0x73, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x6a, 0x73, 0x6f,
	0x6e, 0x22, 0x78, 0x0a, 0x0f, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x72, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x31, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x61,
	0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x52,
	0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x22, 0x22, 0x0a, 0x0c, 0x56,
	0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0x26, 0x0a, 0x10, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x4c,
	0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x44, 0x0a, 0x08, 0x56, 0x61, 0x72, 0x69, 0x61,
	0x62, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74,
	0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x32, 0xd9, 0x03,
	0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x36, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x12,
	0x15, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62,
	0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x1a, 0x0e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74,
	0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x0b, 0x53, 0x65, 0x74, 0x56,
	0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x11, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78,
	0x74, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x1a, 0x1a, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x78, 0x74, 0x2e, 0x53, 0x65, 0x74, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x40, 0x0a, 0x0c, 0x53, 0x65, 0x74, 0x56,
	0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x12, 0x12, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x78, 0x74, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x1a, 0x1a, 0x2e, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x53, 0x65, 0x74, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62,
	0x6c, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x0d, 0x53, 0x65,
	0x74, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x12, 0x2e, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x1a,
	0x1a, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x53, 0x65, 0x74, 0x56, 0x61, 0x72,
	0x69, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x4f, 0x0a,
	0x10, 0x47, 0x65, 0x74, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x73, 0x12, 0x1e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x56,
	0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x1a, 0x19, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x56, 0x61, 0x72, 0x69,
	0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x41,
	0x0a, 0x0c, 0x47, 0x65, 0x74, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x12, 0x1c,
	0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x61, 0x72, 0x69,
	0x61, 0x62, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x4d, 0x61, 0x70, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22,
	0x00, 0x12, 0x36, 0x0a, 0x08, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x12, 0x18, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78,
	0x74, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x00, 0x42, 0x56, 0x0a, 0x21, 0x74, 0x65, 0x63,
	0x68, 0x2e, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x6e, 0x2e, 0x63, 0x6f, 0x6e, 0x63, 0x6f, 0x72, 0x64,
	0x2e, 0x67, 0x6f, 0x6f, 0x64, 0x77, 0x69, 0x6c, 0x6c, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x42, 0x0c,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5a, 0x23, 0x67, 0x6f,
	0x2e, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x6e, 0x2e, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x67, 0x6f, 0x6f,
	0x64, 0x77, 0x69, 0x6c, 0x6c, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_context_proto_rawDescOnce sync.Once
	file_context_proto_rawDescData = file_context_proto_rawDesc
)

func file_context_proto_rawDescGZIP() []byte {
	file_context_proto_rawDescOnce.Do(func() {
		file_context_proto_rawDescData = protoimpl.X.CompressGZIP(file_context_proto_rawDescData)
	})
	return file_context_proto_rawDescData
}

var file_context_proto_msgTypes = make([]protoimpl.MessageInfo, 19)
var file_context_proto_goTypes = []interface{}{
	(*GetVariablesRequest)(nil),   // 0: context.GetVariablesRequest
	(*SetVariableResult)(nil),     // 1: context.SetVariableResult
	(*GetVariableNameParams)(nil), // 2: context.GetVariableNameParams
	(*Variables)(nil),             // 3: context.Variables
	(*Value)(nil),                 // 4: context.Value
	(*MapValue)(nil),              // 5: context.MapValue
	(*ListValue)(nil),             // 6: context.ListValue
	(*StringValue)(nil),           // 7: context.StringValue
	(*IntValue)(nil),              // 8: context.IntValue
	(*DoubleValue)(nil),           // 9: context.DoubleValue
	(*BoolValue)(nil),             // 10: context.BoolValue
	(*TimeValue)(nil),             // 11: context.TimeValue
	(*NullValue)(nil),             // 12: context.NullValue
	(*JSONValue)(nil),             // 13: context.JSONValue
	(*EvaluateRequest)(nil),       // 14: context.EvaluateRequest
	(*VariableName)(nil),          // 15: context.VariableName
	(*VariableNameList)(nil),      // 16: context.VariableNameList
	(*Variable)(nil),              // 17: context.Variable
	nil,                           // 18: context.MapValue.ValueEntry
	(*anypb.Any)(nil),             // 19: google.protobuf.Any
}
var file_context_proto_depIdxs = []int32{
	17, // 0: context.Variables.parameters:type_name -> context.Variable
	19, // 1: context.Value.value:type_name -> google.protobuf.Any
	18, // 2: context.MapValue.value:type_name -> context.MapValue.ValueEntry
	4,  // 3: context.ListValue.value:type_name -> context.Value
	17, // 4: context.EvaluateRequest.parameters:type_name -> context.Variable
	4,  // 5: context.Variable.value:type_name -> context.Value
	4,  // 6: context.MapValue.ValueEntry.value:type_name -> context.Value
	15, // 7: context.ContextService.GetVariable:input_type -> context.VariableName
	17, // 8: context.ContextService.SetVariable:input_type -> context.Variable
	3,  // 9: context.ContextService.SetVariables:input_type -> context.Variables
	3,  // 10: context.ContextService.SetTaskResult:input_type -> context.Variables
	2,  // 11: context.ContextService.GetVariableNames:input_type -> context.GetVariableNameParams
	0,  // 12: context.ContextService.GetVariables:input_type -> context.GetVariablesRequest
	14, // 13: context.ContextService.Evaluate:input_type -> context.EvaluateRequest
	4,  // 14: context.ContextService.GetVariable:output_type -> context.Value
	1,  // 15: context.ContextService.SetVariable:output_type -> context.SetVariableResult
	1,  // 16: context.ContextService.SetVariables:output_type -> context.SetVariableResult
	1,  // 17: context.ContextService.SetTaskResult:output_type -> context.SetVariableResult
	16, // 18: context.ContextService.GetVariableNames:output_type -> context.VariableNameList
	5,  // 19: context.ContextService.GetVariables:output_type -> context.MapValue
	4,  // 20: context.ContextService.Evaluate:output_type -> context.Value
	14, // [14:21] is the sub-list for method output_type
	7,  // [7:14] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_context_proto_init() }
func file_context_proto_init() {
	if File_context_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_context_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVariablesRequest); i {
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
		file_context_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetVariableResult); i {
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
		file_context_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVariableNameParams); i {
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
		file_context_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Variables); i {
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
		file_context_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Value); i {
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
		file_context_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapValue); i {
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
		file_context_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListValue); i {
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
		file_context_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringValue); i {
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
		file_context_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IntValue); i {
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
		file_context_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DoubleValue); i {
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
		file_context_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoolValue); i {
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
		file_context_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeValue); i {
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
		file_context_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NullValue); i {
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
		file_context_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JSONValue); i {
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
		file_context_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EvaluateRequest); i {
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
		file_context_proto_msgTypes[15].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VariableName); i {
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
		file_context_proto_msgTypes[16].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VariableNameList); i {
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
		file_context_proto_msgTypes[17].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Variable); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_context_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   19,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_context_proto_goTypes,
		DependencyIndexes: file_context_proto_depIdxs,
		MessageInfos:      file_context_proto_msgTypes,
	}.Build()
	File_context_proto = out.File
	file_context_proto_rawDesc = nil
	file_context_proto_goTypes = nil
	file_context_proto_depIdxs = nil
}
