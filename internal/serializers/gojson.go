package serializers

import (
	gojson "github.com/goccy/go-json"

	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/models"
)

// GoJSONSerializer implements Serializer interface for goccy/go-json
type GoJSONSerializer struct{}

// NewGoJSONSerializer creates a new GoJSONSerializer
func NewGoJSONSerializer() *GoJSONSerializer {
	return &GoJSONSerializer{}
}

// Name returns the name of the serializer
func (g *GoJSONSerializer) Name() string {
	return "GoJSON"
}

// Marshal serializes a User to JSON bytes using goccy/go-json
func (g *GoJSONSerializer) Marshal(user models.User) ([]byte, error) {
	return gojson.Marshal(user)
}

// Unmarshal deserializes JSON bytes to a User using goccy/go-json
func (g *GoJSONSerializer) Unmarshal(data []byte) (models.User, error) {
	var user models.User
	err := gojson.Unmarshal(data, &user)
	return user, err
}

// MarshalUsers serializes a collection of Users to JSON bytes using goccy/go-json
func (g *GoJSONSerializer) MarshalUsers(users models.Users) ([]byte, error) {
	return gojson.Marshal(users)
}

// UnmarshalUsers deserializes JSON bytes to a collection of Users using goccy/go-json
func (g *GoJSONSerializer) UnmarshalUsers(data []byte) (models.Users, error) {
	var users models.Users
	err := gojson.Unmarshal(data, &users)
	return users, err
}
