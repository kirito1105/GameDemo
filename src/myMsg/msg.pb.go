// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.28.0--rc1
// source: msg.proto

package myMsg

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

type BlockType int32

const (
	BlockType_Null   BlockType = 0
	BlockType_Ground BlockType = 1
)

// Enum value maps for BlockType.
var (
	BlockType_name = map[int32]string{
		0: "Null",
		1: "Ground",
	}
	BlockType_value = map[string]int32{
		"Null":   0,
		"Ground": 1,
	}
)

func (x BlockType) Enum() *BlockType {
	p := new(BlockType)
	*p = x
	return p
}

func (x BlockType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BlockType) Descriptor() protoreflect.EnumDescriptor {
	return file_msg_proto_enumTypes[0].Descriptor()
}

func (BlockType) Type() protoreflect.EnumType {
	return &file_msg_proto_enumTypes[0]
}

func (x BlockType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BlockType.Descriptor instead.
func (BlockType) EnumDescriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{0}
}

type Cmd int32

const (
	Cmd_Pong           Cmd = 0
	Cmd_Authentication Cmd = 1
	Cmd_Move           Cmd = 2
)

// Enum value maps for Cmd.
var (
	Cmd_name = map[int32]string{
		0: "Pong",
		1: "Authentication",
		2: "Move",
	}
	Cmd_value = map[string]int32{
		"Pong":           0,
		"Authentication": 1,
		"Move":           2,
	}
)

func (x Cmd) Enum() *Cmd {
	p := new(Cmd)
	*p = x
	return p
}

func (x Cmd) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Cmd) Descriptor() protoreflect.EnumDescriptor {
	return file_msg_proto_enumTypes[1].Descriptor()
}

func (Cmd) Type() protoreflect.EnumType {
	return &file_msg_proto_enumTypes[1]
}

func (x Cmd) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Cmd.Descriptor instead.
func (Cmd) EnumDescriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{1}
}

type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  BlockType     `protobuf:"varint,1,opt,name=type,proto3,enum=myMsg.BlockType" json:"type,omitempty"`
	Index *LocationInfo `protobuf:"bytes,2,opt,name=index,proto3" json:"index,omitempty"`
	List  []*Obj        `protobuf:"bytes,3,rep,name=list,proto3" json:"list,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{0}
}

func (x *Block) GetType() BlockType {
	if x != nil {
		return x.Type
	}
	return BlockType_Null
}

func (x *Block) GetIndex() *LocationInfo {
	if x != nil {
		return x.Index
	}
	return nil
}

func (x *Block) GetList() []*Obj {
	if x != nil {
		return x.List
	}
	return nil
}

type Obj struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ObjType string        `protobuf:"bytes,1,opt,name=objType,proto3" json:"objType,omitempty"`
	ObjId   string        `protobuf:"bytes,2,opt,name=objId,proto3" json:"objId,omitempty"`
	Index   *LocationInfo `protobuf:"bytes,3,opt,name=index,proto3" json:"index,omitempty"`
}

func (x *Obj) Reset() {
	*x = Obj{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Obj) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Obj) ProtoMessage() {}

func (x *Obj) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Obj.ProtoReflect.Descriptor instead.
func (*Obj) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{1}
}

func (x *Obj) GetObjType() string {
	if x != nil {
		return x.ObjType
	}
	return ""
}

func (x *Obj) GetObjId() string {
	if x != nil {
		return x.ObjId
	}
	return ""
}

func (x *Obj) GetIndex() *LocationInfo {
	if x != nil {
		return x.Index
	}
	return nil
}

type Msg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Authentication *MsgAuthentication `protobuf:"bytes,1,opt,name=Authentication,proto3" json:"Authentication,omitempty"`
	Scene          *MsgScene          `protobuf:"bytes,2,opt,name=scene,proto3" json:"scene,omitempty"`
}

func (x *Msg) Reset() {
	*x = Msg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Msg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Msg) ProtoMessage() {}

func (x *Msg) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Msg.ProtoReflect.Descriptor instead.
func (*Msg) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{2}
}

func (x *Msg) GetAuthentication() *MsgAuthentication {
	if x != nil {
		return x.Authentication
	}
	return nil
}

func (x *Msg) GetScene() *MsgScene {
	if x != nil {
		return x.Scene
	}
	return nil
}

type MsgFromClient struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd            Cmd                `protobuf:"varint,1,opt,name=cmd,proto3,enum=myMsg.Cmd" json:"cmd,omitempty"`
	Authentication *MsgAuthentication `protobuf:"bytes,2,opt,name=Authentication,proto3" json:"Authentication,omitempty"`
	Move           *MsgMove           `protobuf:"bytes,3,opt,name=move,proto3" json:"move,omitempty"`
}

func (x *MsgFromClient) Reset() {
	*x = MsgFromClient{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgFromClient) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgFromClient) ProtoMessage() {}

func (x *MsgFromClient) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgFromClient.ProtoReflect.Descriptor instead.
func (*MsgFromClient) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{3}
}

func (x *MsgFromClient) GetCmd() Cmd {
	if x != nil {
		return x.Cmd
	}
	return Cmd_Pong
}

func (x *MsgFromClient) GetAuthentication() *MsgAuthentication {
	if x != nil {
		return x.Authentication
	}
	return nil
}

func (x *MsgFromClient) GetMove() *MsgMove {
	if x != nil {
		return x.Move
	}
	return nil
}

type MsgMove struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X float32 `protobuf:"fixed32,1,opt,name=x,proto3" json:"x,omitempty"`
	Y float32 `protobuf:"fixed32,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (x *MsgMove) Reset() {
	*x = MsgMove{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgMove) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgMove) ProtoMessage() {}

func (x *MsgMove) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgMove.ProtoReflect.Descriptor instead.
func (*MsgMove) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{4}
}

func (x *MsgMove) GetX() float32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *MsgMove) GetY() float32 {
	if x != nil {
		return x.Y
	}
	return 0
}

type MsgFromService struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Scene *MsgScene `protobuf:"bytes,1,opt,name=scene,proto3" json:"scene,omitempty"`
}

func (x *MsgFromService) Reset() {
	*x = MsgFromService{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgFromService) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgFromService) ProtoMessage() {}

func (x *MsgFromService) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgFromService.ProtoReflect.Descriptor instead.
func (*MsgFromService) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{5}
}

func (x *MsgFromService) GetScene() *MsgScene {
	if x != nil {
		return x.Scene
	}
	return nil
}

type MsgAuthentication struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Addr     string `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	RoomId   string `protobuf:"bytes,3,opt,name=roomId,proto3" json:"roomId,omitempty"`
	Token    string `protobuf:"bytes,4,opt,name=Token,proto3" json:"Token,omitempty"`
}

func (x *MsgAuthentication) Reset() {
	*x = MsgAuthentication{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgAuthentication) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgAuthentication) ProtoMessage() {}

func (x *MsgAuthentication) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgAuthentication.ProtoReflect.Descriptor instead.
func (*MsgAuthentication) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{6}
}

func (x *MsgAuthentication) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *MsgAuthentication) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *MsgAuthentication) GetRoomId() string {
	if x != nil {
		return x.RoomId
	}
	return ""
}

func (x *MsgAuthentication) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type MsgScene struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Blocks []*Block    `protobuf:"bytes,1,rep,name=blocks,proto3" json:"blocks,omitempty"`
	Chars  []*CharInfo `protobuf:"bytes,2,rep,name=chars,proto3" json:"chars,omitempty"`
}

func (x *MsgScene) Reset() {
	*x = MsgScene{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgScene) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgScene) ProtoMessage() {}

func (x *MsgScene) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgScene.ProtoReflect.Descriptor instead.
func (*MsgScene) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{7}
}

func (x *MsgScene) GetBlocks() []*Block {
	if x != nil {
		return x.Blocks
	}
	return nil
}

func (x *MsgScene) GetChars() []*CharInfo {
	if x != nil {
		return x.Chars
	}
	return nil
}

type CharInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string        `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Index    *LocationInfo `protobuf:"bytes,2,opt,name=index,proto3" json:"index,omitempty"`
	IsUser   bool          `protobuf:"varint,3,opt,name=isUser,proto3" json:"isUser,omitempty"`
}

func (x *CharInfo) Reset() {
	*x = CharInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CharInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CharInfo) ProtoMessage() {}

func (x *CharInfo) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CharInfo.ProtoReflect.Descriptor instead.
func (*CharInfo) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{8}
}

func (x *CharInfo) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *CharInfo) GetIndex() *LocationInfo {
	if x != nil {
		return x.Index
	}
	return nil
}

func (x *CharInfo) GetIsUser() bool {
	if x != nil {
		return x.IsUser
	}
	return false
}

type LocationInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X float32 `protobuf:"fixed32,1,opt,name=x,proto3" json:"x,omitempty"`
	Y float32 `protobuf:"fixed32,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (x *LocationInfo) Reset() {
	*x = LocationInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocationInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocationInfo) ProtoMessage() {}

func (x *LocationInfo) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocationInfo.ProtoReflect.Descriptor instead.
func (*LocationInfo) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{9}
}

func (x *LocationInfo) GetX() float32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *LocationInfo) GetY() float32 {
	if x != nil {
		return x.Y
	}
	return 0
}

var File_msg_proto protoreflect.FileDescriptor

var file_msg_proto_rawDesc = []byte{
	0x0a, 0x09, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d, 0x79, 0x4d,
	0x73, 0x67, 0x22, 0x78, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x24, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x6d, 0x79, 0x4d, 0x73,
	0x67, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x29, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x6d, 0x79, 0x4d, 0x73, 0x67, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1e, 0x0a, 0x04,
	0x6c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x6d, 0x79, 0x4d,
	0x73, 0x67, 0x2e, 0x4f, 0x62, 0x6a, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x22, 0x60, 0x0a, 0x03,
	0x4f, 0x62, 0x6a, 0x12, 0x18, 0x0a, 0x07, 0x6f, 0x62, 0x6a, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x62, 0x6a, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x6f, 0x62, 0x6a, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x62,
	0x6a, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x79, 0x4d, 0x73, 0x67, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x6e,
	0x0a, 0x03, 0x4d, 0x73, 0x67, 0x12, 0x40, 0x0a, 0x0e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x6d, 0x79, 0x4d, 0x73, 0x67, 0x2e, 0x4d, 0x73, 0x67, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x05, 0x73, 0x63, 0x65, 0x6e, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x79, 0x4d, 0x73, 0x67, 0x2e, 0x4d,
	0x73, 0x67, 0x53, 0x63, 0x65, 0x6e, 0x65, 0x52, 0x05, 0x73, 0x63, 0x65, 0x6e, 0x65, 0x22, 0x93,
	0x01, 0x0a, 0x0d, 0x4d, 0x73, 0x67, 0x46, 0x72, 0x6f, 0x6d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x12, 0x1c, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e,
	0x6d, 0x79, 0x4d, 0x73, 0x67, 0x2e, 0x43, 0x6d, 0x64, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x40,
	0x0a, 0x0e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6d, 0x79, 0x4d, 0x73, 0x67, 0x2e, 0x4d,
	0x73, 0x67, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x0e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x22, 0x0a, 0x04, 0x6d, 0x6f, 0x76, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e,
	0x2e, 0x6d, 0x79, 0x4d, 0x73, 0x67, 0x2e, 0x4d, 0x73, 0x67, 0x4d, 0x6f, 0x76, 0x65, 0x52, 0x04,
	0x6d, 0x6f, 0x76, 0x65, 0x22, 0x25, 0x0a, 0x07, 0x4d, 0x73, 0x67, 0x4d, 0x6f, 0x76, 0x65, 0x12,
	0x0c, 0x0a, 0x01, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a,
	0x01, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x79, 0x22, 0x37, 0x0a, 0x0e, 0x4d,
	0x73, 0x67, 0x46, 0x72, 0x6f, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x25, 0x0a,
	0x05, 0x73, 0x63, 0x65, 0x6e, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d,
	0x79, 0x4d, 0x73, 0x67, 0x2e, 0x4d, 0x73, 0x67, 0x53, 0x63, 0x65, 0x6e, 0x65, 0x52, 0x05, 0x73,
	0x63, 0x65, 0x6e, 0x65, 0x22, 0x71, 0x0a, 0x11, 0x4d, 0x73, 0x67, 0x41, 0x75, 0x74, 0x68, 0x65,
	0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x6f, 0x6f,
	0x6d, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x57, 0x0a, 0x08, 0x4d, 0x73, 0x67, 0x53, 0x63,
	0x65, 0x6e, 0x65, 0x12, 0x24, 0x0a, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x79, 0x4d, 0x73, 0x67, 0x2e, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x52, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12, 0x25, 0x0a, 0x05, 0x63, 0x68, 0x61,
	0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x79, 0x4d, 0x73, 0x67,
	0x2e, 0x43, 0x68, 0x61, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x63, 0x68, 0x61, 0x72, 0x73,
	0x22, 0x69, 0x0a, 0x08, 0x43, 0x68, 0x61, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1a, 0x0a, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x79, 0x4d, 0x73, 0x67, 0x2e,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x73, 0x55, 0x73, 0x65, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x55, 0x73, 0x65, 0x72, 0x22, 0x2a, 0x0a, 0x0c, 0x4c,
	0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0c, 0x0a, 0x01, 0x78,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x79, 0x2a, 0x21, 0x0a, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x75, 0x6c, 0x6c, 0x10, 0x00, 0x12, 0x0a,
	0x0a, 0x06, 0x47, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0x01, 0x2a, 0x2d, 0x0a, 0x03, 0x43, 0x6d,
	0x64, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x6f, 0x6e, 0x67, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x41,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0x01, 0x12,
	0x08, 0x0a, 0x04, 0x4d, 0x6f, 0x76, 0x65, 0x10, 0x02, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2e, 0x2f,
	0x6d, 0x79, 0x4d, 0x73, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msg_proto_rawDescOnce sync.Once
	file_msg_proto_rawDescData = file_msg_proto_rawDesc
)

func file_msg_proto_rawDescGZIP() []byte {
	file_msg_proto_rawDescOnce.Do(func() {
		file_msg_proto_rawDescData = protoimpl.X.CompressGZIP(file_msg_proto_rawDescData)
	})
	return file_msg_proto_rawDescData
}

var file_msg_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_msg_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_msg_proto_goTypes = []interface{}{
	(BlockType)(0),            // 0: myMsg.BlockType
	(Cmd)(0),                  // 1: myMsg.Cmd
	(*Block)(nil),             // 2: myMsg.Block
	(*Obj)(nil),               // 3: myMsg.Obj
	(*Msg)(nil),               // 4: myMsg.Msg
	(*MsgFromClient)(nil),     // 5: myMsg.MsgFromClient
	(*MsgMove)(nil),           // 6: myMsg.MsgMove
	(*MsgFromService)(nil),    // 7: myMsg.MsgFromService
	(*MsgAuthentication)(nil), // 8: myMsg.MsgAuthentication
	(*MsgScene)(nil),          // 9: myMsg.MsgScene
	(*CharInfo)(nil),          // 10: myMsg.CharInfo
	(*LocationInfo)(nil),      // 11: myMsg.LocationInfo
}
var file_msg_proto_depIdxs = []int32{
	0,  // 0: myMsg.Block.type:type_name -> myMsg.BlockType
	11, // 1: myMsg.Block.index:type_name -> myMsg.LocationInfo
	3,  // 2: myMsg.Block.list:type_name -> myMsg.Obj
	11, // 3: myMsg.Obj.index:type_name -> myMsg.LocationInfo
	8,  // 4: myMsg.Msg.Authentication:type_name -> myMsg.MsgAuthentication
	9,  // 5: myMsg.Msg.scene:type_name -> myMsg.MsgScene
	1,  // 6: myMsg.MsgFromClient.cmd:type_name -> myMsg.Cmd
	8,  // 7: myMsg.MsgFromClient.Authentication:type_name -> myMsg.MsgAuthentication
	6,  // 8: myMsg.MsgFromClient.move:type_name -> myMsg.MsgMove
	9,  // 9: myMsg.MsgFromService.scene:type_name -> myMsg.MsgScene
	2,  // 10: myMsg.MsgScene.blocks:type_name -> myMsg.Block
	10, // 11: myMsg.MsgScene.chars:type_name -> myMsg.CharInfo
	11, // 12: myMsg.CharInfo.index:type_name -> myMsg.LocationInfo
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_msg_proto_init() }
func file_msg_proto_init() {
	if File_msg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_msg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
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
		file_msg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Obj); i {
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
		file_msg_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Msg); i {
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
		file_msg_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgFromClient); i {
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
		file_msg_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgMove); i {
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
		file_msg_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgFromService); i {
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
		file_msg_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgAuthentication); i {
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
		file_msg_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgScene); i {
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
		file_msg_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CharInfo); i {
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
		file_msg_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LocationInfo); i {
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
			RawDescriptor: file_msg_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_msg_proto_goTypes,
		DependencyIndexes: file_msg_proto_depIdxs,
		EnumInfos:         file_msg_proto_enumTypes,
		MessageInfos:      file_msg_proto_msgTypes,
	}.Build()
	File_msg_proto = out.File
	file_msg_proto_rawDesc = nil
	file_msg_proto_goTypes = nil
	file_msg_proto_depIdxs = nil
}
