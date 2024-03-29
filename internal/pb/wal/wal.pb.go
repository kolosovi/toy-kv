// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.1
// source: wal.proto

package wal

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

type KV struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	K []byte `protobuf:"bytes,1,opt,name=k,proto3" json:"k,omitempty"`
	V []byte `protobuf:"bytes,2,opt,name=v,proto3" json:"v,omitempty"`
}

func (x *KV) Reset() {
	*x = KV{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wal_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KV) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KV) ProtoMessage() {}

func (x *KV) ProtoReflect() protoreflect.Message {
	mi := &file_wal_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KV.ProtoReflect.Descriptor instead.
func (*KV) Descriptor() ([]byte, []int) {
	return file_wal_proto_rawDescGZIP(), []int{0}
}

func (x *KV) GetK() []byte {
	if x != nil {
		return x.K
	}
	return nil
}

func (x *KV) GetV() []byte {
	if x != nil {
		return x.V
	}
	return nil
}

type Insert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kv *KV `protobuf:"bytes,1,opt,name=kv,proto3" json:"kv,omitempty"`
}

func (x *Insert) Reset() {
	*x = Insert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wal_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Insert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Insert) ProtoMessage() {}

func (x *Insert) ProtoReflect() protoreflect.Message {
	mi := &file_wal_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Insert.ProtoReflect.Descriptor instead.
func (*Insert) Descriptor() ([]byte, []int) {
	return file_wal_proto_rawDescGZIP(), []int{1}
}

func (x *Insert) GetKv() *KV {
	if x != nil {
		return x.Kv
	}
	return nil
}

type Delete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	K []byte `protobuf:"bytes,1,opt,name=k,proto3" json:"k,omitempty"`
}

func (x *Delete) Reset() {
	*x = Delete{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wal_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Delete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Delete) ProtoMessage() {}

func (x *Delete) ProtoReflect() protoreflect.Message {
	mi := &file_wal_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Delete.ProtoReflect.Descriptor instead.
func (*Delete) Descriptor() ([]byte, []int) {
	return file_wal_proto_rawDescGZIP(), []int{2}
}

func (x *Delete) GetK() []byte {
	if x != nil {
		return x.K
	}
	return nil
}

type Log struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Log:
	//
	//	*Log_Insert
	//	*Log_Delete
	Log isLog_Log `protobuf_oneof:"log"`
}

func (x *Log) Reset() {
	*x = Log{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wal_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Log) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Log) ProtoMessage() {}

func (x *Log) ProtoReflect() protoreflect.Message {
	mi := &file_wal_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Log.ProtoReflect.Descriptor instead.
func (*Log) Descriptor() ([]byte, []int) {
	return file_wal_proto_rawDescGZIP(), []int{3}
}

func (m *Log) GetLog() isLog_Log {
	if m != nil {
		return m.Log
	}
	return nil
}

func (x *Log) GetInsert() *Insert {
	if x, ok := x.GetLog().(*Log_Insert); ok {
		return x.Insert
	}
	return nil
}

func (x *Log) GetDelete() *Delete {
	if x, ok := x.GetLog().(*Log_Delete); ok {
		return x.Delete
	}
	return nil
}

type isLog_Log interface {
	isLog_Log()
}

type Log_Insert struct {
	Insert *Insert `protobuf:"bytes,1,opt,name=insert,proto3,oneof"`
}

type Log_Delete struct {
	Delete *Delete `protobuf:"bytes,2,opt,name=delete,proto3,oneof"`
}

func (*Log_Insert) isLog_Log() {}

func (*Log_Delete) isLog_Log() {}

var File_wal_proto protoreflect.FileDescriptor

var file_wal_proto_rawDesc = []byte{
	0x0a, 0x09, 0x77, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x74, 0x6f, 0x79,
	0x5f, 0x6b, 0x76, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x77, 0x61, 0x6c,
	0x22, 0x20, 0x0a, 0x02, 0x4b, 0x56, 0x12, 0x0c, 0x0a, 0x01, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x01, 0x6b, 0x12, 0x0c, 0x0a, 0x01, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x01, 0x76, 0x22, 0x31, 0x0a, 0x06, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x12, 0x27, 0x0a, 0x02,
	0x6b, 0x76, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x74, 0x6f, 0x79, 0x5f, 0x6b,
	0x76, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x77, 0x61, 0x6c, 0x2e, 0x4b,
	0x56, 0x52, 0x02, 0x6b, 0x76, 0x22, 0x16, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12,
	0x0c, 0x0a, 0x01, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x01, 0x6b, 0x22, 0x7a, 0x0a,
	0x03, 0x4c, 0x6f, 0x67, 0x12, 0x35, 0x0a, 0x06, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x6f, 0x79, 0x5f, 0x6b, 0x76, 0x2e, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x77, 0x61, 0x6c, 0x2e, 0x49, 0x6e, 0x73, 0x65, 0x72,
	0x74, 0x48, 0x00, 0x52, 0x06, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x12, 0x35, 0x0a, 0x06, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x6f,
	0x79, 0x5f, 0x6b, 0x76, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x77, 0x61,
	0x6c, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x48, 0x00, 0x52, 0x06, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x42, 0x05, 0x0a, 0x03, 0x6c, 0x6f, 0x67, 0x42, 0x11, 0x5a, 0x0f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62, 0x2f, 0x77, 0x61, 0x6c, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_wal_proto_rawDescOnce sync.Once
	file_wal_proto_rawDescData = file_wal_proto_rawDesc
)

func file_wal_proto_rawDescGZIP() []byte {
	file_wal_proto_rawDescOnce.Do(func() {
		file_wal_proto_rawDescData = protoimpl.X.CompressGZIP(file_wal_proto_rawDescData)
	})
	return file_wal_proto_rawDescData
}

var file_wal_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_wal_proto_goTypes = []interface{}{
	(*KV)(nil),     // 0: toy_kv.internal.wal.KV
	(*Insert)(nil), // 1: toy_kv.internal.wal.Insert
	(*Delete)(nil), // 2: toy_kv.internal.wal.Delete
	(*Log)(nil),    // 3: toy_kv.internal.wal.Log
}
var file_wal_proto_depIdxs = []int32{
	0, // 0: toy_kv.internal.wal.Insert.kv:type_name -> toy_kv.internal.wal.KV
	1, // 1: toy_kv.internal.wal.Log.insert:type_name -> toy_kv.internal.wal.Insert
	2, // 2: toy_kv.internal.wal.Log.delete:type_name -> toy_kv.internal.wal.Delete
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_wal_proto_init() }
func file_wal_proto_init() {
	if File_wal_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_wal_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KV); i {
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
		file_wal_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Insert); i {
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
		file_wal_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Delete); i {
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
		file_wal_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Log); i {
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
	file_wal_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Log_Insert)(nil),
		(*Log_Delete)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_wal_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_wal_proto_goTypes,
		DependencyIndexes: file_wal_proto_depIdxs,
		MessageInfos:      file_wal_proto_msgTypes,
	}.Build()
	File_wal_proto = out.File
	file_wal_proto_rawDesc = nil
	file_wal_proto_goTypes = nil
	file_wal_proto_depIdxs = nil
}
