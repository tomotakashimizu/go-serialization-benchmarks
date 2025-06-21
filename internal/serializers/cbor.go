package serializers

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/models"
)

// CBORSerializer implements Serializer interface for CBOR
type CBORSerializer struct{}

// NewCBORSerializer creates a new CBORSerializer
func NewCBORSerializer() *CBORSerializer {
	return &CBORSerializer{}
}

// Name returns the name of the serializer
func (c *CBORSerializer) Name() string {
	return "CBOR"
}

// Marshal serializes a User to CBOR bytes
func (c *CBORSerializer) Marshal(user models.User) ([]byte, error) {
	return cbor.Marshal(user)
}

// Unmarshal deserializes CBOR bytes to a User
func (c *CBORSerializer) Unmarshal(data []byte) (models.User, error) {
	var user models.User
	err := cbor.Unmarshal(data, &user)
	return user, err
}

// MarshalUsers serializes a collection of Users to CBOR bytes
func (c *CBORSerializer) MarshalUsers(users models.Users) ([]byte, error) {
	return cbor.Marshal(users)
}

// UnmarshalUsers deserializes CBOR bytes to a collection of Users
func (c *CBORSerializer) UnmarshalUsers(data []byte) (models.Users, error) {
	var users models.Users
	err := cbor.Unmarshal(data, &users)
	return users, err
}
