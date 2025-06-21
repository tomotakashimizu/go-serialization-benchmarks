package serializers

import (
	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/models"
	"github.com/vmihailenco/msgpack/v5"
)

// MsgPackSerializer implements Serializer interface for MessagePack
type MsgPackSerializer struct{}

// NewMsgPackSerializer creates a new MsgPackSerializer
func NewMsgPackSerializer() *MsgPackSerializer {
	return &MsgPackSerializer{}
}

// Name returns the name of the serializer
func (m *MsgPackSerializer) Name() string {
	return "msgpack"
}

// Marshal serializes a User to MessagePack bytes
func (m *MsgPackSerializer) Marshal(user models.User) ([]byte, error) {
	return msgpack.Marshal(user)
}

// Unmarshal deserializes MessagePack bytes to a User
func (m *MsgPackSerializer) Unmarshal(data []byte) (models.User, error) {
	var user models.User
	err := msgpack.Unmarshal(data, &user)
	return user, err
}

// MarshalUsers serializes a collection of Users to MessagePack bytes
func (m *MsgPackSerializer) MarshalUsers(users models.Users) ([]byte, error) {
	return msgpack.Marshal(users)
}

// UnmarshalUsers deserializes MessagePack bytes to a collection of Users
func (m *MsgPackSerializer) UnmarshalUsers(data []byte) (models.Users, error) {
	var users models.Users
	err := msgpack.Unmarshal(data, &users)
	return users, err
}
