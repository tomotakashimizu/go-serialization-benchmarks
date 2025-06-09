package serializers

import (
	"bytes"
	"encoding/gob"

	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/models"
)

// GobSerializer implements Serializer interface for Gob
type GobSerializer struct{}

// NewGobSerializer creates a new GobSerializer
func NewGobSerializer() *GobSerializer {
	return &GobSerializer{}
}

// Name returns the name of the serializer
func (g *GobSerializer) Name() string {
	return "Gob"
}

// Marshal serializes a User to Gob bytes
func (g *GobSerializer) Marshal(user models.User) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(user)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal deserializes Gob bytes to a User
func (g *GobSerializer) Unmarshal(data []byte) (models.User, error) {
	var user models.User
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&user)
	return user, err
}

// MarshalUsers serializes a slice of Users to Gob bytes
func (g *GobSerializer) MarshalUsers(users []models.User) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(users)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalUsers deserializes Gob bytes to a slice of Users
func (g *GobSerializer) UnmarshalUsers(data []byte) ([]models.User, error) {
	var users []models.User
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&users)
	return users, err
}
