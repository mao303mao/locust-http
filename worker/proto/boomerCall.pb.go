// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v4.0.0
// source: boomerCall.proto

package boomerCall

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type InitBommerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsSession    bool                             `protobuf:"varint,1,opt,name=isSession,proto3" json:"isSession,omitempty"`
	PreTask      []*InitBommerRequest_HttpRequest `protobuf:"bytes,2,rep,name=PreTask,proto3" json:"PreTask,omitempty"`   // 前置任务
	MainTask     []*InitBommerRequest_TestTask    `protobuf:"bytes,3,rep,name=MainTask,proto3" json:"MainTask,omitempty"` // 测试任务
	LocustMaster string                           `protobuf:"bytes,4,opt,name=LocustMaster,proto3" json:"LocustMaster,omitempty"`
	HttpProxy    string                           `protobuf:"bytes,5,opt,name=HttpProxy,proto3" json:"HttpProxy,omitempty"`
}

func (x *InitBommerRequest) Reset() {
	*x = InitBommerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_boomerCall_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitBommerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitBommerRequest) ProtoMessage() {}

func (x *InitBommerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_boomerCall_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitBommerRequest.ProtoReflect.Descriptor instead.
func (*InitBommerRequest) Descriptor() ([]byte, []int) {
	return file_boomerCall_proto_rawDescGZIP(), []int{0}
}

func (x *InitBommerRequest) GetIsSession() bool {
	if x != nil {
		return x.IsSession
	}
	return false
}

func (x *InitBommerRequest) GetPreTask() []*InitBommerRequest_HttpRequest {
	if x != nil {
		return x.PreTask
	}
	return nil
}

func (x *InitBommerRequest) GetMainTask() []*InitBommerRequest_TestTask {
	if x != nil {
		return x.MainTask
	}
	return nil
}

func (x *InitBommerRequest) GetLocustMaster() string {
	if x != nil {
		return x.LocustMaster
	}
	return ""
}

func (x *InitBommerRequest) GetHttpProxy() string {
	if x != nil {
		return x.HttpProxy
	}
	return ""
}

type EndBommerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EndBommerRequest) Reset() {
	*x = EndBommerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_boomerCall_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EndBommerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EndBommerRequest) ProtoMessage() {}

func (x *EndBommerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_boomerCall_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EndBommerRequest.ProtoReflect.Descriptor instead.
func (*EndBommerRequest) Descriptor() ([]byte, []int) {
	return file_boomerCall_proto_rawDescGZIP(), []int{1}
}

type BoomerCallResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  bool   `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *BoomerCallResponse) Reset() {
	*x = BoomerCallResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_boomerCall_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoomerCallResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoomerCallResponse) ProtoMessage() {}

func (x *BoomerCallResponse) ProtoReflect() protoreflect.Message {
	mi := &file_boomerCall_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoomerCallResponse.ProtoReflect.Descriptor instead.
func (*BoomerCallResponse) Descriptor() ([]byte, []int) {
	return file_boomerCall_proto_rawDescGZIP(), []int{2}
}

func (x *BoomerCallResponse) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

func (x *BoomerCallResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type InitBommerRequest_SaveParamAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SaveType  int32  `protobuf:"varint,1,opt,name=SaveType,proto3" json:"SaveType,omitempty"`  // 0-HTML-XPATH解析；1-JSON-解析; 2-文本正则匹配; 3-保存固定字符串
	ParamName string `protobuf:"bytes,2,opt,name=ParamName,proto3" json:"ParamName,omitempty"` // 全局变量名称
	RuleValue string `protobuf:"bytes,3,opt,name=RuleValue,proto3" json:"RuleValue,omitempty"`
}

func (x *InitBommerRequest_SaveParamAction) Reset() {
	*x = InitBommerRequest_SaveParamAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_boomerCall_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitBommerRequest_SaveParamAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitBommerRequest_SaveParamAction) ProtoMessage() {}

func (x *InitBommerRequest_SaveParamAction) ProtoReflect() protoreflect.Message {
	mi := &file_boomerCall_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitBommerRequest_SaveParamAction.ProtoReflect.Descriptor instead.
func (*InitBommerRequest_SaveParamAction) Descriptor() ([]byte, []int) {
	return file_boomerCall_proto_rawDescGZIP(), []int{0, 0}
}

func (x *InitBommerRequest_SaveParamAction) GetSaveType() int32 {
	if x != nil {
		return x.SaveType
	}
	return 0
}

func (x *InitBommerRequest_SaveParamAction) GetParamName() string {
	if x != nil {
		return x.ParamName
	}
	return ""
}

func (x *InitBommerRequest_SaveParamAction) GetRuleValue() string {
	if x != nil {
		return x.RuleValue
	}
	return ""
}

type InitBommerRequest_AssertAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AssertType int32 `protobuf:"varint,1,opt,name=AssertType,proto3" json:"AssertType,omitempty"` // 0-状态码等于; 1-响应内容字节长度小于；2-响应内容直接长度等于；3-响应内容直接长度大于
	RuleValue  int32 `protobuf:"varint,2,opt,name=RuleValue,proto3" json:"RuleValue,omitempty"`
}

func (x *InitBommerRequest_AssertAction) Reset() {
	*x = InitBommerRequest_AssertAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_boomerCall_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitBommerRequest_AssertAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitBommerRequest_AssertAction) ProtoMessage() {}

func (x *InitBommerRequest_AssertAction) ProtoReflect() protoreflect.Message {
	mi := &file_boomerCall_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitBommerRequest_AssertAction.ProtoReflect.Descriptor instead.
func (*InitBommerRequest_AssertAction) Descriptor() ([]byte, []int) {
	return file_boomerCall_proto_rawDescGZIP(), []int{0, 1}
}

func (x *InitBommerRequest_AssertAction) GetAssertType() int32 {
	if x != nil {
		return x.AssertType
	}
	return 0
}

func (x *InitBommerRequest_AssertAction) GetRuleValue() int32 {
	if x != nil {
		return x.RuleValue
	}
	return 0
}

type InitBommerRequest_HttpRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UrlPath        string                               `protobuf:"bytes,1,opt,name=UrlPath,proto3" json:"UrlPath,omitempty"`
	Method         string                               `protobuf:"bytes,2,opt,name=Method,proto3" json:"Method,omitempty"`
	Headers        map[string]string                    `protobuf:"bytes,3,rep,name=Headers,proto3" json:"Headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	DictData       map[string]string                    `protobuf:"bytes,4,rep,name=DictData,proto3" json:"DictData,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Params         map[string]string                    `protobuf:"bytes,5,rep,name=Params,proto3" json:"Params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	RawData        string                               `protobuf:"bytes,6,opt,name=RawData,proto3" json:"RawData,omitempty"`
	JsonData       string                               `protobuf:"bytes,7,opt,name=JsonData,proto3" json:"JsonData,omitempty"`
	SaveParamChain []*InitBommerRequest_SaveParamAction `protobuf:"bytes,8,rep,name=SaveParamChain,proto3" json:"SaveParamChain,omitempty"`
	AssertChain    []*InitBommerRequest_AssertAction    `protobuf:"bytes,9,rep,name=AssertChain,proto3" json:"AssertChain,omitempty"`
}

func (x *InitBommerRequest_HttpRequest) Reset() {
	*x = InitBommerRequest_HttpRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_boomerCall_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitBommerRequest_HttpRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitBommerRequest_HttpRequest) ProtoMessage() {}

func (x *InitBommerRequest_HttpRequest) ProtoReflect() protoreflect.Message {
	mi := &file_boomerCall_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitBommerRequest_HttpRequest.ProtoReflect.Descriptor instead.
func (*InitBommerRequest_HttpRequest) Descriptor() ([]byte, []int) {
	return file_boomerCall_proto_rawDescGZIP(), []int{0, 2}
}

func (x *InitBommerRequest_HttpRequest) GetUrlPath() string {
	if x != nil {
		return x.UrlPath
	}
	return ""
}

func (x *InitBommerRequest_HttpRequest) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *InitBommerRequest_HttpRequest) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *InitBommerRequest_HttpRequest) GetDictData() map[string]string {
	if x != nil {
		return x.DictData
	}
	return nil
}

func (x *InitBommerRequest_HttpRequest) GetParams() map[string]string {
	if x != nil {
		return x.Params
	}
	return nil
}

func (x *InitBommerRequest_HttpRequest) GetRawData() string {
	if x != nil {
		return x.RawData
	}
	return ""
}

func (x *InitBommerRequest_HttpRequest) GetJsonData() string {
	if x != nil {
		return x.JsonData
	}
	return ""
}

func (x *InitBommerRequest_HttpRequest) GetSaveParamChain() []*InitBommerRequest_SaveParamAction {
	if x != nil {
		return x.SaveParamChain
	}
	return nil
}

func (x *InitBommerRequest_HttpRequest) GetAssertChain() []*InitBommerRequest_AssertAction {
	if x != nil {
		return x.AssertChain
	}
	return nil
}

type InitBommerRequest_TestTask struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskName   string                           `protobuf:"bytes,1,opt,name=TaskName,proto3" json:"TaskName,omitempty"`
	TaskWeight int32                            `protobuf:"varint,2,opt,name=TaskWeight,proto3" json:"TaskWeight,omitempty"`
	PreWork    []*InitBommerRequest_HttpRequest `protobuf:"bytes,3,rep,name=PreWork,proto3" json:"PreWork,omitempty"`   // 准备工作请求，比如获取token
	TestWork   *InitBommerRequest_HttpRequest   `protobuf:"bytes,4,opt,name=TestWork,proto3" json:"TestWork,omitempty"` // 主要性能测试任务
}

func (x *InitBommerRequest_TestTask) Reset() {
	*x = InitBommerRequest_TestTask{}
	if protoimpl.UnsafeEnabled {
		mi := &file_boomerCall_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitBommerRequest_TestTask) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitBommerRequest_TestTask) ProtoMessage() {}

func (x *InitBommerRequest_TestTask) ProtoReflect() protoreflect.Message {
	mi := &file_boomerCall_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitBommerRequest_TestTask.ProtoReflect.Descriptor instead.
func (*InitBommerRequest_TestTask) Descriptor() ([]byte, []int) {
	return file_boomerCall_proto_rawDescGZIP(), []int{0, 3}
}

func (x *InitBommerRequest_TestTask) GetTaskName() string {
	if x != nil {
		return x.TaskName
	}
	return ""
}

func (x *InitBommerRequest_TestTask) GetTaskWeight() int32 {
	if x != nil {
		return x.TaskWeight
	}
	return 0
}

func (x *InitBommerRequest_TestTask) GetPreWork() []*InitBommerRequest_HttpRequest {
	if x != nil {
		return x.PreWork
	}
	return nil
}

func (x *InitBommerRequest_TestTask) GetTestWork() *InitBommerRequest_HttpRequest {
	if x != nil {
		return x.TestWork
	}
	return nil
}

var File_boomerCall_proto protoreflect.FileDescriptor

var file_boomerCall_proto_rawDesc = []byte{
	0x0a, 0x10, 0x62, 0x6f, 0x6f, 0x6d, 0x65, 0x72, 0x43, 0x61, 0x6c, 0x6c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xee, 0x09, 0x0a, 0x11, 0x49, 0x6e, 0x69, 0x74, 0x42, 0x6f, 0x6d, 0x6d, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x73, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x54, 0x61, 0x73,
	0x6b, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x42, 0x6f,
	0x6d, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x48, 0x74, 0x74, 0x70,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x07, 0x50, 0x72, 0x65, 0x54, 0x61, 0x73, 0x6b,
	0x12, 0x37, 0x0a, 0x08, 0x4d, 0x61, 0x69, 0x6e, 0x54, 0x61, 0x73, 0x6b, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x42, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x52,
	0x08, 0x4d, 0x61, 0x69, 0x6e, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x22, 0x0a, 0x0c, 0x4c, 0x6f, 0x63,
	0x75, 0x73, 0x74, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x4c, 0x6f, 0x63, 0x75, 0x73, 0x74, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x12, 0x1c, 0x0a,
	0x09, 0x48, 0x74, 0x74, 0x70, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x48, 0x74, 0x74, 0x70, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x1a, 0x69, 0x0a, 0x0f, 0x53,
	0x61, 0x76, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a,
	0x0a, 0x08, 0x53, 0x61, 0x76, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x08, 0x53, 0x61, 0x76, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x52, 0x75, 0x6c, 0x65,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x52, 0x75, 0x6c,
	0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x4c, 0x0a, 0x0c, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74,
	0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x41, 0x73, 0x73, 0x65,
	0x72, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x52, 0x75, 0x6c, 0x65, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x52, 0x75, 0x6c, 0x65, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x1a, 0x8d, 0x05, 0x0a, 0x0b, 0x48, 0x74, 0x74, 0x70, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x55, 0x72, 0x6c, 0x50, 0x61, 0x74, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x55, 0x72, 0x6c, 0x50, 0x61, 0x74, 0x68, 0x12, 0x16,
	0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x45, 0x0a, 0x07, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x42, 0x6f,
	0x6d, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x48, 0x74, 0x74, 0x70,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x48, 0x0a,
	0x08, 0x44, 0x69, 0x63, 0x74, 0x44, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x2c, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x42, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x2e, 0x48, 0x74, 0x74, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x44, 0x69, 0x63, 0x74, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x44,
	0x69, 0x63, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x42, 0x0a, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x42, 0x6f,
	0x6d, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x48, 0x74, 0x74, 0x70,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x52,
	0x61, 0x77, 0x44, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x52, 0x61,
	0x77, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x4a, 0x73, 0x6f, 0x6e, 0x44, 0x61, 0x74,
	0x61, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4a, 0x73, 0x6f, 0x6e, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x4a, 0x0a, 0x0e, 0x53, 0x61, 0x76, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x43, 0x68,
	0x61, 0x69, 0x6e, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x49, 0x6e, 0x69, 0x74,
	0x42, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x61,
	0x76, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0e, 0x53,
	0x61, 0x76, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x12, 0x41, 0x0a,
	0x0b, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x18, 0x09, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x42, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x43, 0x68, 0x61, 0x69, 0x6e,
	0x1a, 0x3a, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3b, 0x0a, 0x0d,
	0x44, 0x69, 0x63, 0x74, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x39, 0x0a, 0x0b, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x1a, 0xbc, 0x01, 0x0a, 0x08, 0x54, 0x65, 0x73, 0x74, 0x54, 0x61, 0x73,
	0x6b, 0x12, 0x1a, 0x0a, 0x08, 0x54, 0x61, 0x73, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x54, 0x61, 0x73, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a,
	0x0a, 0x54, 0x61, 0x73, 0x6b, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0a, 0x54, 0x61, 0x73, 0x6b, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x38, 0x0a,
	0x07, 0x50, 0x72, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x49, 0x6e, 0x69, 0x74, 0x42, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x2e, 0x48, 0x74, 0x74, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x07,
	0x50, 0x72, 0x65, 0x57, 0x6f, 0x72, 0x6b, 0x12, 0x3a, 0x0a, 0x08, 0x54, 0x65, 0x73, 0x74, 0x57,
	0x6f, 0x72, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x49, 0x6e, 0x69, 0x74,
	0x42, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x48, 0x74,
	0x74, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x08, 0x54, 0x65, 0x73, 0x74, 0x57,
	0x6f, 0x72, 0x6b, 0x22, 0x12, 0x0a, 0x10, 0x45, 0x6e, 0x64, 0x42, 0x6f, 0x6d, 0x6d, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x46, 0x0a, 0x12, 0x42, 0x6f, 0x6f, 0x6d, 0x65,
	0x72, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32,
	0x83, 0x01, 0x0a, 0x11, 0x42, 0x6f, 0x6f, 0x6d, 0x65, 0x72, 0x43, 0x61, 0x6c, 0x6c, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x37, 0x0a, 0x0a, 0x49, 0x6e, 0x69, 0x74, 0x42, 0x6f, 0x6d,
	0x6d, 0x65, 0x72, 0x12, 0x12, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x42, 0x6f, 0x6d, 0x6d, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x42, 0x6f, 0x6f, 0x6d, 0x65, 0x72,
	0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x35,
	0x0a, 0x09, 0x45, 0x6e, 0x64, 0x42, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x12, 0x11, 0x2e, 0x45, 0x6e,
	0x64, 0x42, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13,
	0x2e, 0x42, 0x6f, 0x6f, 0x6d, 0x65, 0x72, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0e, 0x5a, 0x0c, 0x2e, 0x3b, 0x62, 0x6f, 0x6f, 0x6d, 0x65,
	0x72, 0x43, 0x61, 0x6c, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_boomerCall_proto_rawDescOnce sync.Once
	file_boomerCall_proto_rawDescData = file_boomerCall_proto_rawDesc
)

func file_boomerCall_proto_rawDescGZIP() []byte {
	file_boomerCall_proto_rawDescOnce.Do(func() {
		file_boomerCall_proto_rawDescData = protoimpl.X.CompressGZIP(file_boomerCall_proto_rawDescData)
	})
	return file_boomerCall_proto_rawDescData
}

var file_boomerCall_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_boomerCall_proto_goTypes = []interface{}{
	(*InitBommerRequest)(nil),                 // 0: InitBommerRequest
	(*EndBommerRequest)(nil),                  // 1: EndBommerRequest
	(*BoomerCallResponse)(nil),                // 2: BoomerCallResponse
	(*InitBommerRequest_SaveParamAction)(nil), // 3: InitBommerRequest.SaveParamAction
	(*InitBommerRequest_AssertAction)(nil),    // 4: InitBommerRequest.AssertAction
	(*InitBommerRequest_HttpRequest)(nil),     // 5: InitBommerRequest.HttpRequest
	(*InitBommerRequest_TestTask)(nil),        // 6: InitBommerRequest.TestTask
	nil,                                       // 7: InitBommerRequest.HttpRequest.HeadersEntry
	nil,                                       // 8: InitBommerRequest.HttpRequest.DictDataEntry
	nil,                                       // 9: InitBommerRequest.HttpRequest.ParamsEntry
}
var file_boomerCall_proto_depIdxs = []int32{
	5,  // 0: InitBommerRequest.PreTask:type_name -> InitBommerRequest.HttpRequest
	6,  // 1: InitBommerRequest.MainTask:type_name -> InitBommerRequest.TestTask
	7,  // 2: InitBommerRequest.HttpRequest.Headers:type_name -> InitBommerRequest.HttpRequest.HeadersEntry
	8,  // 3: InitBommerRequest.HttpRequest.DictData:type_name -> InitBommerRequest.HttpRequest.DictDataEntry
	9,  // 4: InitBommerRequest.HttpRequest.Params:type_name -> InitBommerRequest.HttpRequest.ParamsEntry
	3,  // 5: InitBommerRequest.HttpRequest.SaveParamChain:type_name -> InitBommerRequest.SaveParamAction
	4,  // 6: InitBommerRequest.HttpRequest.AssertChain:type_name -> InitBommerRequest.AssertAction
	5,  // 7: InitBommerRequest.TestTask.PreWork:type_name -> InitBommerRequest.HttpRequest
	5,  // 8: InitBommerRequest.TestTask.TestWork:type_name -> InitBommerRequest.HttpRequest
	0,  // 9: BoomerCallService.InitBommer:input_type -> InitBommerRequest
	1,  // 10: BoomerCallService.EndBommer:input_type -> EndBommerRequest
	2,  // 11: BoomerCallService.InitBommer:output_type -> BoomerCallResponse
	2,  // 12: BoomerCallService.EndBommer:output_type -> BoomerCallResponse
	11, // [11:13] is the sub-list for method output_type
	9,  // [9:11] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_boomerCall_proto_init() }
func file_boomerCall_proto_init() {
	if File_boomerCall_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_boomerCall_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitBommerRequest); i {
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
		file_boomerCall_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EndBommerRequest); i {
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
		file_boomerCall_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoomerCallResponse); i {
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
		file_boomerCall_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitBommerRequest_SaveParamAction); i {
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
		file_boomerCall_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitBommerRequest_AssertAction); i {
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
		file_boomerCall_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitBommerRequest_HttpRequest); i {
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
		file_boomerCall_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitBommerRequest_TestTask); i {
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
			RawDescriptor: file_boomerCall_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_boomerCall_proto_goTypes,
		DependencyIndexes: file_boomerCall_proto_depIdxs,
		MessageInfos:      file_boomerCall_proto_msgTypes,
	}.Build()
	File_boomerCall_proto = out.File
	file_boomerCall_proto_rawDesc = nil
	file_boomerCall_proto_goTypes = nil
	file_boomerCall_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// BoomerCallServiceClient is the client API for BoomerCallService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BoomerCallServiceClient interface {
	InitBommer(ctx context.Context, in *InitBommerRequest, opts ...grpc.CallOption) (*BoomerCallResponse, error)
	EndBommer(ctx context.Context, in *EndBommerRequest, opts ...grpc.CallOption) (*BoomerCallResponse, error)
}

type boomerCallServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBoomerCallServiceClient(cc grpc.ClientConnInterface) BoomerCallServiceClient {
	return &boomerCallServiceClient{cc}
}

func (c *boomerCallServiceClient) InitBommer(ctx context.Context, in *InitBommerRequest, opts ...grpc.CallOption) (*BoomerCallResponse, error) {
	out := new(BoomerCallResponse)
	err := c.cc.Invoke(ctx, "/BoomerCallService/InitBommer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boomerCallServiceClient) EndBommer(ctx context.Context, in *EndBommerRequest, opts ...grpc.CallOption) (*BoomerCallResponse, error) {
	out := new(BoomerCallResponse)
	err := c.cc.Invoke(ctx, "/BoomerCallService/EndBommer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BoomerCallServiceServer is the server API for BoomerCallService service.
type BoomerCallServiceServer interface {
	InitBommer(context.Context, *InitBommerRequest) (*BoomerCallResponse, error)
	EndBommer(context.Context, *EndBommerRequest) (*BoomerCallResponse, error)
}

// UnimplementedBoomerCallServiceServer can be embedded to have forward compatible implementations.
type UnimplementedBoomerCallServiceServer struct {
}

func (*UnimplementedBoomerCallServiceServer) InitBommer(context.Context, *InitBommerRequest) (*BoomerCallResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitBommer not implemented")
}
func (*UnimplementedBoomerCallServiceServer) EndBommer(context.Context, *EndBommerRequest) (*BoomerCallResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EndBommer not implemented")
}

func RegisterBoomerCallServiceServer(s *grpc.Server, srv BoomerCallServiceServer) {
	s.RegisterService(&_BoomerCallService_serviceDesc, srv)
}

func _BoomerCallService_InitBommer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitBommerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoomerCallServiceServer).InitBommer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/BoomerCallService/InitBommer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoomerCallServiceServer).InitBommer(ctx, req.(*InitBommerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BoomerCallService_EndBommer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndBommerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoomerCallServiceServer).EndBommer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/BoomerCallService/EndBommer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoomerCallServiceServer).EndBommer(ctx, req.(*EndBommerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _BoomerCallService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "BoomerCallService",
	HandlerType: (*BoomerCallServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InitBommer",
			Handler:    _BoomerCallService_InitBommer_Handler,
		},
		{
			MethodName: "EndBommer",
			Handler:    _BoomerCallService_EndBommer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "boomerCall.proto",
}