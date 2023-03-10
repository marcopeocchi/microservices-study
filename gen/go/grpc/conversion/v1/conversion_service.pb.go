// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: grpc/conversion/v1/conversion_service.proto

package conversionpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RunRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Job *ConversionJob `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty"`
}

func (x *RunRequest) Reset() {
	*x = RunRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_conversion_v1_conversion_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RunRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RunRequest) ProtoMessage() {}

func (x *RunRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_conversion_v1_conversion_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RunRequest.ProtoReflect.Descriptor instead.
func (*RunRequest) Descriptor() ([]byte, []int) {
	return file_grpc_conversion_v1_conversion_service_proto_rawDescGZIP(), []int{0}
}

func (x *RunRequest) GetJob() *ConversionJob {
	if x != nil {
		return x.Job
	}
	return nil
}

type RunResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Job       *ConversionJob `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty"`
	Converted int32          `protobuf:"varint,2,opt,name=converted,proto3" json:"converted,omitempty"`
}

func (x *RunResponse) Reset() {
	*x = RunResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_conversion_v1_conversion_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RunResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RunResponse) ProtoMessage() {}

func (x *RunResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_conversion_v1_conversion_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RunResponse.ProtoReflect.Descriptor instead.
func (*RunResponse) Descriptor() ([]byte, []int) {
	return file_grpc_conversion_v1_conversion_service_proto_rawDescGZIP(), []int{1}
}

func (x *RunResponse) GetJob() *ConversionJob {
	if x != nil {
		return x.Job
	}
	return nil
}

func (x *RunResponse) GetConverted() int32 {
	if x != nil {
		return x.Converted
	}
	return 0
}

var File_grpc_conversion_v1_conversion_service_proto protoreflect.FileDescriptor

var file_grpc_conversion_v1_conversion_service_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x63,
	0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x23, 0x67, 0x72,
	0x70, 0x63, 0x2f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x3c, 0x0a, 0x0a, 0x52, 0x75, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2e, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63,
	0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x4a, 0x6f, 0x62, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x22,
	0x5b, 0x0a, 0x0b, 0x52, 0x75, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e,
	0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x6f,
	0x6e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x4a, 0x6f, 0x62, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x12, 0x1c,
	0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x65, 0x64, 0x32, 0x53, 0x0a, 0x11,
	0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x3e, 0x0a, 0x03, 0x52, 0x75, 0x6e, 0x12, 0x19, 0x2e, 0x63, 0x6f, 0x6e, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x75, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x75, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x29, 0x5a, 0x27, 0x66, 0x75, 0x75, 0x2f, 0x76, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67,
	0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x3b,
	0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_conversion_v1_conversion_service_proto_rawDescOnce sync.Once
	file_grpc_conversion_v1_conversion_service_proto_rawDescData = file_grpc_conversion_v1_conversion_service_proto_rawDesc
)

func file_grpc_conversion_v1_conversion_service_proto_rawDescGZIP() []byte {
	file_grpc_conversion_v1_conversion_service_proto_rawDescOnce.Do(func() {
		file_grpc_conversion_v1_conversion_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_conversion_v1_conversion_service_proto_rawDescData)
	})
	return file_grpc_conversion_v1_conversion_service_proto_rawDescData
}

var file_grpc_conversion_v1_conversion_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_grpc_conversion_v1_conversion_service_proto_goTypes = []interface{}{
	(*RunRequest)(nil),    // 0: conversion.v1.RunRequest
	(*RunResponse)(nil),   // 1: conversion.v1.RunResponse
	(*ConversionJob)(nil), // 2: conversion.v1.ConversionJob
}
var file_grpc_conversion_v1_conversion_service_proto_depIdxs = []int32{
	2, // 0: conversion.v1.RunRequest.job:type_name -> conversion.v1.ConversionJob
	2, // 1: conversion.v1.RunResponse.job:type_name -> conversion.v1.ConversionJob
	0, // 2: conversion.v1.ConversionService.Run:input_type -> conversion.v1.RunRequest
	1, // 3: conversion.v1.ConversionService.Run:output_type -> conversion.v1.RunResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_grpc_conversion_v1_conversion_service_proto_init() }
func file_grpc_conversion_v1_conversion_service_proto_init() {
	if File_grpc_conversion_v1_conversion_service_proto != nil {
		return
	}
	file_grpc_conversion_v1_conversion_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_grpc_conversion_v1_conversion_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RunRequest); i {
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
		file_grpc_conversion_v1_conversion_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RunResponse); i {
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
			RawDescriptor: file_grpc_conversion_v1_conversion_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_conversion_v1_conversion_service_proto_goTypes,
		DependencyIndexes: file_grpc_conversion_v1_conversion_service_proto_depIdxs,
		MessageInfos:      file_grpc_conversion_v1_conversion_service_proto_msgTypes,
	}.Build()
	File_grpc_conversion_v1_conversion_service_proto = out.File
	file_grpc_conversion_v1_conversion_service_proto_rawDesc = nil
	file_grpc_conversion_v1_conversion_service_proto_goTypes = nil
	file_grpc_conversion_v1_conversion_service_proto_depIdxs = nil
}