package logicServer

type RoomToClien struct {
	IsFind   bool   `protobuf:"varint,1,opt,name=IsFind,proto3" json:"IsFind,omitempty"`
	RoomId   string `protobuf:"bytes,2,opt,name=RoomId,proto3" json:"RoomId,omitempty"`
	RoomAddr string `protobuf:"bytes,3,opt,name=RoomAddr,proto3" json:"RoomAddr,omitempty"`
	Token    string `protobuf:"bytes,4,opt,name=Token,proto3" json:"Token,omitempty"`
}
