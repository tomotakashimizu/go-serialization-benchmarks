package serializers

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/models"
)

// JSONiterSerializer implements Serializer interface for json-iterator
type JSONiterSerializer struct {
	json jsoniter.API
}

// NewJSONiterSerializer creates a new JSONiterSerializer
func NewJSONiterSerializer() *JSONiterSerializer {
	// ConfigCompatibleWithStandardLibrary provides 100% compatibility with standard lib
	return &JSONiterSerializer{
		json: jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

// Name returns the name of the serializer
func (j *JSONiterSerializer) Name() string {
	return "JSONiter"
}

// Marshal serializes a User to JSON bytes using json-iterator
func (j *JSONiterSerializer) Marshal(user models.User) ([]byte, error) {
	return j.json.Marshal(user)
}

// Unmarshal deserializes JSON bytes to a User using json-iterator
func (j *JSONiterSerializer) Unmarshal(data []byte) (models.User, error) {
	var user models.User
	err := j.json.Unmarshal(data, &user)
	return user, err
}

// MarshalUsers serializes a slice of Users to JSON bytes using json-iterator
func (j *JSONiterSerializer) MarshalUsers(users []models.User) ([]byte, error) {
	return j.json.Marshal(users)
}

// UnmarshalUsers deserializes JSON bytes to a slice of Users using json-iterator
func (j *JSONiterSerializer) UnmarshalUsers(data []byte) ([]models.User, error) {
	var users []models.User
	err := j.json.Unmarshal(data, &users)
	return users, err
}
