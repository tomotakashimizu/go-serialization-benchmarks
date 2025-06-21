package serializers

import (
	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/models"
)

// EasyJSONSerializer implements Serializer interface for EasyJSON
type EasyJSONSerializer struct{}

// NewEasyJSONSerializer creates a new EasyJSONSerializer
func NewEasyJSONSerializer() *EasyJSONSerializer {
	return &EasyJSONSerializer{}
}

// Name returns the name of the serializer
func (e *EasyJSONSerializer) Name() string {
	return "EasyJSON"
}

// Marshal serializes a User to JSON bytes using EasyJSON
func (e *EasyJSONSerializer) Marshal(user models.User) ([]byte, error) {
	return user.MarshalJSON()
}

// Unmarshal deserializes JSON bytes to a User using EasyJSON
func (e *EasyJSONSerializer) Unmarshal(data []byte) (models.User, error) {
	var user models.User
	err := user.UnmarshalJSON(data)
	return user, err
}

// MarshalUsers serializes a collection of Users to JSON bytes using EasyJSON
func (e *EasyJSONSerializer) MarshalUsers(users models.Users) ([]byte, error) {
	return users.MarshalJSON()
}

// UnmarshalUsers deserializes JSON bytes to a collection of Users using EasyJSON
func (e *EasyJSONSerializer) UnmarshalUsers(data []byte) (models.Users, error) {
	var users models.Users
	err := users.UnmarshalJSON(data)
	return users, err
}
