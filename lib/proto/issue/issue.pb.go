// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.27.0--rc1
// source: lib/proto/issue.proto

package issue

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

type IssueCategory int32

const (
	IssueCategory_Unknown    IssueCategory = 0
	IssueCategory_Bug        IssueCategory = 1
	IssueCategory_Feature    IssueCategory = 2
	IssueCategory_Operations IssueCategory = 3
	IssueCategory_Question   IssueCategory = 4
)

// Enum value maps for IssueCategory.
var (
	IssueCategory_name = map[int32]string{
		0: "Unknown",
		1: "Bug",
		2: "Feature",
		3: "Operations",
		4: "Question",
	}
	IssueCategory_value = map[string]int32{
		"Unknown":    0,
		"Bug":        1,
		"Feature":    2,
		"Operations": 3,
		"Question":   4,
	}
)

func (x IssueCategory) Enum() *IssueCategory {
	p := new(IssueCategory)
	*p = x
	return p
}

func (x IssueCategory) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (IssueCategory) Descriptor() protoreflect.EnumDescriptor {
	return file_lib_proto_issue_proto_enumTypes[0].Descriptor()
}

func (IssueCategory) Type() protoreflect.EnumType {
	return &file_lib_proto_issue_proto_enumTypes[0]
}

func (x IssueCategory) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use IssueCategory.Descriptor instead.
func (IssueCategory) EnumDescriptor() ([]byte, []int) {
	return file_lib_proto_issue_proto_rawDescGZIP(), []int{0}
}

type Tag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Label     string `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty"`
	ColorCode string `protobuf:"bytes,2,opt,name=color_code,json=colorCode,proto3" json:"color_code,omitempty"`
}

func (x *Tag) Reset() {
	*x = Tag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_proto_issue_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tag) ProtoMessage() {}

func (x *Tag) ProtoReflect() protoreflect.Message {
	mi := &file_lib_proto_issue_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tag.ProtoReflect.Descriptor instead.
func (*Tag) Descriptor() ([]byte, []int) {
	return file_lib_proto_issue_proto_rawDescGZIP(), []int{0}
}

func (x *Tag) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *Tag) GetColorCode() string {
	if x != nil {
		return x.ColorCode
	}
	return ""
}

type ViewStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number uint64 `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	SetAt  int64  `protobuf:"varint,2,opt,name=set_at,json=setAt,proto3" json:"set_at,omitempty"`
}

func (x *ViewStatus) Reset() {
	*x = ViewStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_proto_issue_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ViewStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ViewStatus) ProtoMessage() {}

func (x *ViewStatus) ProtoReflect() protoreflect.Message {
	mi := &file_lib_proto_issue_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ViewStatus.ProtoReflect.Descriptor instead.
func (*ViewStatus) Descriptor() ([]byte, []int) {
	return file_lib_proto_issue_proto_rawDescGZIP(), []int{1}
}

func (x *ViewStatus) GetNumber() uint64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *ViewStatus) GetSetAt() int64 {
	if x != nil {
		return x.SetAt
	}
	return 0
}

type FlowStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number uint64 `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	SetAt  int64  `protobuf:"varint,2,opt,name=set_at,json=setAt,proto3" json:"set_at,omitempty"`
}

func (x *FlowStatus) Reset() {
	*x = FlowStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_proto_issue_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlowStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlowStatus) ProtoMessage() {}

func (x *FlowStatus) ProtoReflect() protoreflect.Message {
	mi := &file_lib_proto_issue_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlowStatus.ProtoReflect.Descriptor instead.
func (*FlowStatus) Descriptor() ([]byte, []int) {
	return file_lib_proto_issue_proto_rawDescGZIP(), []int{2}
}

func (x *FlowStatus) GetNumber() uint64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *FlowStatus) GetSetAt() int64 {
	if x != nil {
		return x.SetAt
	}
	return 0
}

type Issue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number    uint64        `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	CreatedAt int64         `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Title     string        `protobuf:"bytes,20,opt,name=title,proto3" json:"title,omitempty"`
	Body      string        `protobuf:"bytes,21,opt,name=body,proto3" json:"body,omitempty"`
	Category  IssueCategory `protobuf:"varint,22,opt,name=category,proto3,enum=brainslurp.IssueCategory" json:"category,omitempty"`
	Tags      []*Tag        `protobuf:"bytes,40,rep,name=tags,proto3" json:"tags,omitempty"`
	Views     []*ViewStatus `protobuf:"bytes,50,rep,name=views,proto3" json:"views,omitempty"`
	Flows     []*FlowStatus `protobuf:"bytes,51,rep,name=flows,proto3" json:"flows,omitempty"`
}

func (x *Issue) Reset() {
	*x = Issue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lib_proto_issue_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Issue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Issue) ProtoMessage() {}

func (x *Issue) ProtoReflect() protoreflect.Message {
	mi := &file_lib_proto_issue_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Issue.ProtoReflect.Descriptor instead.
func (*Issue) Descriptor() ([]byte, []int) {
	return file_lib_proto_issue_proto_rawDescGZIP(), []int{3}
}

func (x *Issue) GetNumber() uint64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *Issue) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Issue) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Issue) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *Issue) GetCategory() IssueCategory {
	if x != nil {
		return x.Category
	}
	return IssueCategory_Unknown
}

func (x *Issue) GetTags() []*Tag {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *Issue) GetViews() []*ViewStatus {
	if x != nil {
		return x.Views
	}
	return nil
}

func (x *Issue) GetFlows() []*FlowStatus {
	if x != nil {
		return x.Flows
	}
	return nil
}

var File_lib_proto_issue_proto protoreflect.FileDescriptor

var file_lib_proto_issue_proto_rawDesc = []byte{
	0x0a, 0x15, 0x6c, 0x69, 0x62, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x73, 0x73, 0x75,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x62, 0x72, 0x61, 0x69, 0x6e, 0x73, 0x6c,
	0x75, 0x72, 0x70, 0x22, 0x3a, 0x0a, 0x03, 0x54, 0x61, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x61,
	0x62, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c,
	0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x22,
	0x3b, 0x0a, 0x0a, 0x56, 0x69, 0x65, 0x77, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a,
	0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x15, 0x0a, 0x06, 0x73, 0x65, 0x74, 0x5f, 0x61, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x65, 0x74, 0x41, 0x74, 0x22, 0x3b, 0x0a, 0x0a,
	0x46, 0x6c, 0x6f, 0x77, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x15, 0x0a, 0x06, 0x73, 0x65, 0x74, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x73, 0x65, 0x74, 0x41, 0x74, 0x22, 0xa0, 0x02, 0x0a, 0x05, 0x49, 0x73,
	0x73, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x15, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x62, 0x6f, 0x64, 0x79, 0x12, 0x35, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
	0x18, 0x16, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x62, 0x72, 0x61, 0x69, 0x6e, 0x73, 0x6c,
	0x75, 0x72, 0x70, 0x2e, 0x49, 0x73, 0x73, 0x75, 0x65, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72,
	0x79, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x23, 0x0a, 0x04, 0x74,
	0x61, 0x67, 0x73, 0x18, 0x28, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x62, 0x72, 0x61, 0x69,
	0x6e, 0x73, 0x6c, 0x75, 0x72, 0x70, 0x2e, 0x54, 0x61, 0x67, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73,
	0x12, 0x2c, 0x0a, 0x05, 0x76, 0x69, 0x65, 0x77, 0x73, 0x18, 0x32, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x62, 0x72, 0x61, 0x69, 0x6e, 0x73, 0x6c, 0x75, 0x72, 0x70, 0x2e, 0x56, 0x69, 0x65,
	0x77, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x05, 0x76, 0x69, 0x65, 0x77, 0x73, 0x12, 0x2c,
	0x0a, 0x05, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x18, 0x33, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x62, 0x72, 0x61, 0x69, 0x6e, 0x73, 0x6c, 0x75, 0x72, 0x70, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x05, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2a, 0x50, 0x0a, 0x0d,
	0x49, 0x73, 0x73, 0x75, 0x65, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x0b, 0x0a,
	0x07, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x42, 0x75,
	0x67, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x10, 0x02,
	0x12, 0x0e, 0x0a, 0x0a, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x10, 0x03,
	0x12, 0x0c, 0x0a, 0x08, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0x04, 0x42, 0x36,
	0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x65, 0x76,
	0x6e, 0x75, 0x6c, 0x6c, 0x2d, 0x74, 0x77, 0x69, 0x74, 0x63, 0x68, 0x2f, 0x62, 0x72, 0x61, 0x69,
	0x6e, 0x73, 0x6c, 0x75, 0x72, 0x70, 0x2f, 0x6c, 0x69, 0x62, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x69, 0x73, 0x73, 0x75, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_lib_proto_issue_proto_rawDescOnce sync.Once
	file_lib_proto_issue_proto_rawDescData = file_lib_proto_issue_proto_rawDesc
)

func file_lib_proto_issue_proto_rawDescGZIP() []byte {
	file_lib_proto_issue_proto_rawDescOnce.Do(func() {
		file_lib_proto_issue_proto_rawDescData = protoimpl.X.CompressGZIP(file_lib_proto_issue_proto_rawDescData)
	})
	return file_lib_proto_issue_proto_rawDescData
}

var file_lib_proto_issue_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_lib_proto_issue_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_lib_proto_issue_proto_goTypes = []interface{}{
	(IssueCategory)(0), // 0: brainslurp.IssueCategory
	(*Tag)(nil),        // 1: brainslurp.Tag
	(*ViewStatus)(nil), // 2: brainslurp.ViewStatus
	(*FlowStatus)(nil), // 3: brainslurp.FlowStatus
	(*Issue)(nil),      // 4: brainslurp.Issue
}
var file_lib_proto_issue_proto_depIdxs = []int32{
	0, // 0: brainslurp.Issue.category:type_name -> brainslurp.IssueCategory
	1, // 1: brainslurp.Issue.tags:type_name -> brainslurp.Tag
	2, // 2: brainslurp.Issue.views:type_name -> brainslurp.ViewStatus
	3, // 3: brainslurp.Issue.flows:type_name -> brainslurp.FlowStatus
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_lib_proto_issue_proto_init() }
func file_lib_proto_issue_proto_init() {
	if File_lib_proto_issue_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_lib_proto_issue_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tag); i {
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
		file_lib_proto_issue_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ViewStatus); i {
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
		file_lib_proto_issue_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlowStatus); i {
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
		file_lib_proto_issue_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Issue); i {
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
			RawDescriptor: file_lib_proto_issue_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_lib_proto_issue_proto_goTypes,
		DependencyIndexes: file_lib_proto_issue_proto_depIdxs,
		EnumInfos:         file_lib_proto_issue_proto_enumTypes,
		MessageInfos:      file_lib_proto_issue_proto_msgTypes,
	}.Build()
	File_lib_proto_issue_proto = out.File
	file_lib_proto_issue_proto_rawDesc = nil
	file_lib_proto_issue_proto_goTypes = nil
	file_lib_proto_issue_proto_depIdxs = nil
}
