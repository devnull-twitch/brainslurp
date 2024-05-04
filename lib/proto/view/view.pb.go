// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.27.0--rc1
// source: lib/proto/view.proto

package view

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

type View struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number    uint64 `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	CreatedAt int64  `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Title     string `protobuf:"bytes,20,opt,name=title,proto3" json:"title,omitempty"`
}

func (x *View) Reset() {
	*x = View{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_proto_view_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *View) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*View) ProtoMessage() {}

func (x *View) ProtoReflect() protoreflect.Message {
	mi := &file_lib_proto_view_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use View.ProtoReflect.Descriptor instead.
func (*View) Descriptor() ([]byte, []int) {
	return file_lib_proto_view_proto_rawDescGZIP(), []int{0}
}

func (x *View) GetNumber() uint64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *View) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *View) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

var File_lib_proto_view_proto protoreflect.FileDescriptor

var file_lib_proto_view_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6c, 0x69, 0x62, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x69, 0x65, 0x77,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x62, 0x72, 0x61, 0x69, 0x6e, 0x73, 0x6c, 0x75,
	0x72, 0x70, 0x22, 0x53, 0x0a, 0x04, 0x56, 0x69, 0x65, 0x77, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x65, 0x76, 0x6e, 0x75, 0x6c, 0x6c, 0x2d, 0x74, 0x77,
	0x69, 0x74, 0x63, 0x68, 0x2f, 0x62, 0x72, 0x61, 0x69, 0x6e, 0x73, 0x6c, 0x75, 0x72, 0x70, 0x2f,
	0x6c, 0x69, 0x62, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x69, 0x65, 0x77, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_lib_proto_view_proto_rawDescOnce sync.Once
	file_lib_proto_view_proto_rawDescData = file_lib_proto_view_proto_rawDesc
)

func file_lib_proto_view_proto_rawDescGZIP() []byte {
	file_lib_proto_view_proto_rawDescOnce.Do(func() {
		file_lib_proto_view_proto_rawDescData = protoimpl.X.CompressGZIP(file_lib_proto_view_proto_rawDescData)
	})
	return file_lib_proto_view_proto_rawDescData
}

var file_lib_proto_view_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_lib_proto_view_proto_goTypes = []interface{}{
	(*View)(nil), // 0: brainslurp.View
}
var file_lib_proto_view_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_lib_proto_view_proto_init() }
func file_lib_proto_view_proto_init() {
	if File_lib_proto_view_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_lib_proto_view_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*View); i {
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
			RawDescriptor: file_lib_proto_view_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_lib_proto_view_proto_goTypes,
		DependencyIndexes: file_lib_proto_view_proto_depIdxs,
		MessageInfos:      file_lib_proto_view_proto_msgTypes,
	}.Build()
	File_lib_proto_view_proto = out.File
	file_lib_proto_view_proto_rawDesc = nil
	file_lib_proto_view_proto_goTypes = nil
	file_lib_proto_view_proto_depIdxs = nil
}
