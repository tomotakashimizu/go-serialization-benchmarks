package serializers

import (
	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/models"
)

// MsgpSerializer implements Serializer interface for tinylib/msgp
type MsgpSerializer struct{}

// NewMsgpSerializer creates a new MsgpSerializer
func NewMsgpSerializer() *MsgpSerializer {
	return &MsgpSerializer{}
}

// Name returns the name of the serializer
func (m *MsgpSerializer) Name() string {
	return "msgp"
}

// Marshal serializes a User to MessagePack bytes using tinylib/msgp
func (m *MsgpSerializer) Marshal(user models.User) ([]byte, error) {
	return user.MarshalMsg(nil)
}

// Unmarshal deserializes MessagePack bytes to a User using tinylib/msgp
func (m *MsgpSerializer) Unmarshal(data []byte) (models.User, error) {
	var user models.User
	_, err := user.UnmarshalMsg(data)
	return user, err
}

// MarshalUsers serializes a collection of Users to MessagePack bytes using tinylib/msgp
func (m *MsgpSerializer) MarshalUsers(users models.Users) ([]byte, error) {
	return users.MarshalMsg(nil)
}

// UnmarshalUsers deserializes MessagePack bytes to a collection of Users using tinylib/msgp
func (m *MsgpSerializer) UnmarshalUsers(data []byte) (models.Users, error) {
	var users models.Users
	_, err := users.UnmarshalMsg(data)
	return users, err
}
