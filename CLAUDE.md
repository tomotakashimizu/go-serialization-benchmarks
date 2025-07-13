# Go Serialization Benchmarks

## Tech Stack

- Language: Go 1.24.2
- Testing: Built-in testing package + custom benchmarking
- Dependencies: Multiple serialization libraries (see go.mod)
- Optional: Redis for cache performance testing

## Project Structure

- `cmd/benchmark/`: Application entry point (main.go)
- `internal/benchmark/`: Benchmark execution logic
- `internal/models/`: Test data structures and generation
- `internal/serializers/`: Serializer implementations for each format
- `internal/flatbuffers/`: FlatBuffers schema and generated code
- `internal/proto/`: Protocol Buffers schema and generated code
- `internal/redis/`: Redis performance testing client
- `internal/reporter/`: Result output and CSV generation
- `internal/utils/`: Statistical calculation utilities
- `results/`: Benchmark output files (CSV format)

## Commands

- `go run cmd/benchmark/main.go`: Run all benchmarks with default settings
- `go run cmd/benchmark/main.go -help`: Show available command-line options
- `go run cmd/benchmark/main.go -count=10000 -skip-redis`: Run with 10k records, skip Redis tests
- `go mod tidy`: Update and clean dependencies
- `go generate ./...`: Regenerate EasyJSON and Msgp code (if needed)

## Code Generation

- EasyJSON: `//go:generate easyjson -all test_data.go` in models package
- Msgp: `//go:generate msgp` in models package
- FlatBuffers: Manual schema compilation required for schema changes
- Protobuf: Manual compilation required for schema changes

## Code Style

- Follow standard Go formatting with `gofmt`
- Use descriptive variable names, especially for serializer instances
- Prefer composition over inheritance for serializer implementations
- Keep error handling explicit and informative
- Use consistent struct tags across all models: `json:"field" msgpack:"field" cbor:"field" msg:"field"`

## Benchmark Design Principles

- Measure both marshal and unmarshal operations separately
- Calculate both average and median for statistical accuracy
- Use realistic nested data structures (4-layer nesting)
- Test with large datasets (default: 100,000 records)
- Include symmetry tests for type preservation (empty/nil slices and maps)

## Adding New Serializers

1. Create new file in `internal/serializers/` (e.g., `newformat.go`)
2. Implement the `Serializer` interface:
   - `Name() string`
   - `Marshal(user models.User) ([]byte, error)`
   - `Unmarshal(data []byte) (models.User, error)`
   - `MarshalUsers(users models.Users) ([]byte, error)`
   - `UnmarshalUsers(data []byte) (models.Users, error)`
3. Add appropriate struct tags to models if needed
4. Register serializer in `cmd/benchmark/main.go` (alphabetical order after JSON)
5. Test with symmetry tests for type preservation behavior

## Testing Guidelines

- Run benchmarks with consistent hardware conditions
- Use multiple iterations (default: 5) for reliable averages
- Test with both small and large datasets for scalability insights
- Verify symmetry test results for production use cases
- Include Redis tests only when Redis server is available

## Performance Considerations

- Some serializers require code generation (EasyJSON, Msgp, FlatBuffers, Protobuf)
- Code generation serializers typically offer better performance
- FlatBuffers optimizes for zero-copy deserialization
- Binary formats (CBOR, MsgPack, Protobuf) generally produce smaller output
- JSON variants focus on compatibility vs. performance trade-offs

## Do Not

- Modify generated code files directly (regenerate instead)
- Change test data structure without updating all serializer implementations
- Commit Redis credentials or connection details
- Run benchmarks on shared/loaded systems for accurate measurements
- Compare results across different hardware without noting differences
