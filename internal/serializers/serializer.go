package serializers

import "github.com/tomotakashimizu/performance-testing-toolkit/benchmarks/serializer/internal/models"

// Serializer defines the interface for serialization operations
type Serializer interface {
	Name() string
	Marshal(user models.User) ([]byte, error)
	Unmarshal(data []byte) (models.User, error)
	MarshalUsers(users []models.User) ([]byte, error)
	UnmarshalUsers(data []byte) ([]models.User, error)
}

// SerializationResult contains the results of serialization benchmarks
type SerializationResult struct {
	SerializerName    string
	MarshalTimes      []int64 // nanoseconds
	UnmarshalTimes    []int64 // nanoseconds
	DataSize          int     // bytes
	MarshalAvgNs      int64
	MarshalMedianNs   int64
	UnmarshalAvgNs    int64
	UnmarshalMedianNs int64
}

// SymmetryResult contains the results of strict type preservation tests
type SymmetryResult struct {
	SerializerName      string
	StrictEmptySlicesOK bool // Strict type preservation ([] stays [])
	StrictEmptyMapsOK   bool // Strict type preservation ({} stays {})
	StrictNilSlicesOK   bool // Strict nil preservation (nil stays nil)
	StrictNilMapsOK     bool // Strict nil preservation (nil stays nil)
	Details             string
}
