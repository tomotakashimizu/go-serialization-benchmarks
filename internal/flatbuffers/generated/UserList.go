// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package generated

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type UserList struct {
	_tab flatbuffers.Table
}

func GetRootAsUserList(buf []byte, offset flatbuffers.UOffsetT) *UserList {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &UserList{}
	x.Init(buf, n+offset)
	return x
}

func FinishUserListBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsUserList(buf []byte, offset flatbuffers.UOffsetT) *UserList {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &UserList{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedUserListBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *UserList) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *UserList) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *UserList) Users(obj *User, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *UserList) UsersLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func UserListStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func UserListAddUsers(builder *flatbuffers.Builder, users flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(users), 0)
}
func UserListStartUsersVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func UserListEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
