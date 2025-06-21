package serializers

import (
	"encoding/json"

	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/models"
)

// JSONSerializer implements Serializer interface for JSON
type JSONSerializer struct{}

// NewJSONSerializer creates a new JSONSerializer
func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

// Name returns the name of the serializer
func (j *JSONSerializer) Name() string {
	return "JSON"
}

// Marshal serializes a User to JSON bytes
func (j *JSONSerializer) Marshal(user models.User) ([]byte, error) {
	return json.Marshal(user)
}

// Unmarshal deserializes JSON bytes to a User
func (j *JSONSerializer) Unmarshal(data []byte) (models.User, error) {
	var user models.User
	err := json.Unmarshal(data, &user)
	return user, err
}

// MarshalUsers serializes a collection of Users to JSON bytes
func (j *JSONSerializer) MarshalUsers(users models.Users) ([]byte, error) {
	return json.Marshal(users)
}

// UnmarshalUsers deserializes JSON bytes to a collection of Users
func (j *JSONSerializer) UnmarshalUsers(data []byte) (models.Users, error) {
	var users models.Users
	err := json.Unmarshal(data, &users)
	return users, err
}
