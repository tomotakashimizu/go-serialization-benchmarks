# Go Serialization Benchmarks

[日本語版](./README.ja.md) | English

A performance comparison tool for various serialization formats in Go.

## Overview

This project is a benchmark tool for measuring and comparing the performance of different serialization formats. It is particularly useful for:

- Choosing the optimal serialization format for your application
- Performance optimization in data persistence and API communication
- Evaluating serialization performance for Redis caching

## Supported Serializers

- **JSON** - Go standard library ([`encoding/json`](https://pkg.go.dev/encoding/json))
- **CBOR** - [`github.com/fxamacker/cbor/v2`](https://github.com/fxamacker/cbor)
- **EasyJSON** - High-performance JSON with code generation ([`github.com/mailru/easyjson`](https://github.com/mailru/easyjson)) - Code generation based high-performance JSON serializer
- **FlatBuffers** - Zero-copy serialization ([`github.com/google/flatbuffers`](https://github.com/google/flatbuffers)) - Memory-efficient cross-platform serialization format
- **Gob** - Go standard library ([`encoding/gob`](https://pkg.go.dev/encoding/gob))
- **GoJSON** - High-performance JSON ([`github.com/goccy/go-json`](https://github.com/goccy/go-json)) - 100% compatible high-performance version of the standard library
- **JSONiter** - High-performance JSON ([`github.com/json-iterator/go`](https://github.com/json-iterator/go)) - 100% compatible high-performance version of the standard library
- **Msgp** - High-performance MessagePack with code generation ([`github.com/tinylib/msgp`](https://github.com/tinylib/msgp)) - Code generation based high-performance MessagePack serializer
- **MsgPack** - [`github.com/vmihailenco/msgpack/v5`](https://github.com/vmihailenco/msgpack)
- **Protobuf** - Google Protocol Buffers ([`google.golang.org/protobuf`](https://pkg.go.dev/google.golang.org/protobuf#section-readme)) - Efficient, language-neutral serialization format

## Measurements

### 1. Serialization Performance

- Marshal/Unmarshal speed (average and median)
- Serialized data size

### 2. Marshal/Unmarshal Symmetry Tests

- Empty slice/map Marshal→Unmarshal symmetry
- Nil slice/map Marshal→Unmarshal symmetry

### 3. Redis Performance Measurements (Optional)

- Redis SET/GET operation performance
- Evaluation in actual cache usage scenarios

## Project Structure

```plaintext
go-serialization-benchmarks/
├── cmd/
│   └── benchmark/
│       └── main.go                 # Execution entry point
├── internal/
│   ├── benchmark/
│   │   └── runner.go              # Benchmark execution logic
│   ├── models/
│   │   └── test_data.go           # Test data structures
│   ├── flatbuffers/
│   │   ├── user.fbs               # FlatBuffers schema definition
│   │   └── generated/             # FlatBuffers generated code
│   ├── proto/
│   │   ├── user.proto             # Protocol Buffers schema definition
│   │   └── user.pb.go             # Generated Protocol Buffers code
│   ├── redis/
│   │   └── client.go              # Redis performance measurement
│   ├── reporter/
│   │   └── reporter.go            # Result output and saving
│   └── serializers/
│       ├── serializer.go          # Common interface
│       ├── json.go                # JSON implementation
│       ├── cbor.go                # CBOR implementation
│       ├── easyjson.go            # EasyJSON implementation
│       ├── flatbuffers.go         # FlatBuffers implementation
│       ├── gob.go                 # Gob implementation
│       ├── gojson.go              # GoJSON implementation
│       ├── jsoniter.go            # JSONiter implementation
│       ├── msgp.go                # Msgp implementation
│       ├── msgpack.go             # MsgPack implementation
│       └── protobuf.go            # Protobuf implementation
├── results/                        # Result output directory
├── go.mod                          # Go module configuration
└── README.md                       # This file
```

## Installation

```bash
git clone https://github.com/tomotakashimizu/go-serialization-benchmarks.git
cd go-serialization-benchmarks
go mod tidy
```

## Usage

### Basic Execution

```bash
# Run with default settings (100,000 records, 5 iterations)
go run cmd/benchmark/main.go

# Run with specified number of records
go run cmd/benchmark/main.go -count=10000

# Skip Redis measurements
go run cmd/benchmark/main.go -skip-redis

# Show help
go run cmd/benchmark/main.go -help
```

### Command Line Arguments

| Argument          | Default        | Description                 |
| ----------------- | -------------- | --------------------------- |
| `-count`          | 100000         | Number of test records      |
| `-iterations`     | 5              | Number of benchmark runs    |
| `-redis-addr`     | localhost:6379 | Redis server address        |
| `-redis-password` | ""             | Redis password              |
| `-redis-db`       | 0              | Redis database number       |
| `-output`         | ./results      | Result output directory     |
| `-skip-redis`     | false          | Skip Redis measurements     |
| `-help`           | false          | Show help                   |

### Execution Examples

```bash
# Small test (10,000 records, no Redis)
go run cmd/benchmark/main.go -count=10000 -skip-redis

# Run with custom Redis settings
go run cmd/benchmark/main.go -redis-addr=192.168.1.100:6379 -redis-password=secret

# Run with 10 iterations
go run cmd/benchmark/main.go -iterations=10
```

## Test Data

Uses a User model with 4-layer nested structure:

```go
type User struct {
    ID        int64
    Name      string
    Email     string
    Age       int
    IsActive  bool
    Profile   Profile                // Layer 2
    Settings  Settings               // Layer 2
    Tags      []string
    Metadata  map[string]interface{}
    CreatedAt time.Time
}

type Profile struct {
    FirstName   string
    LastName    string
    Bio         string
    Avatar      string
    SocialLinks []Link             // Layer 3
    Preferences Preferences        // Layer 3
}

type Preferences struct {
    Theme         string
    Language      string
    Notifications map[string]bool
    Privacy       PrivacySettings  // Layer 4
}
```

The complex nested structure enables benchmarking under conditions close to actual application data.

## Result Output

### Console Output

The following information is displayed in table format during execution:

#### Command Execution Example

```bash
go run cmd/benchmark/main.go -count=10000 -skip-redis
```

```bash
Serializer Performance Benchmark
=================================
Test data count: 10000
Benchmark iterations: 5
Output directory: ./results
Redis: localhost:6379 (skip: true)

========================================================================================================================
SERIALIZATION BENCHMARK RESULTS
========================================================================================================================
Serializer   | Data Size    | Marshal Avg  | Marshal Med  | Unmarshal Avg | Unmarshal Med
             | (MB)         | (ms)         | (ms)         | (ms)         | (ms)        
------------------------------------------------------------------------------------------------------------------------
JSON         | 9.06         | 139.28       | 137.27       | 218.02       | 215.54      
CBOR         | 6.87         | 56.14        | 56.50        | 131.92       | 132.44      
EasyJSON     | 9.06         | 37.16        | 35.80        | 92.34        | 87.53       
FlatBuffers  | 9.27         | 94.89        | 93.25        | 49.40        | 47.49       
Gob          | 4.72         | 62.01        | 61.74        | 88.61        | 88.28       
GoJSON       | 9.06         | 116.68       | 121.31       | 122.82       | 130.12      
JSONiter     | 9.06         | 42.17        | 36.15        | 147.17       | 139.42      
Msgp         | 6.94         | 25.53        | 25.08        | 38.18        | 38.39       
MsgPack      | 6.94         | 77.75        | 74.31        | 165.28       | 147.20      
Protobuf     | 4.97         | 114.01       | 111.78       | 134.50       | 136.57      
========================================================================================================================

====================================================================================================
STRICT TYPE PRESERVATION TEST RESULTS
====================================================================================================
Serializer   | Empty→Empty  | Empty{}→{}   | Nil→Nil      | Nil→Nil     
             | (Slices)     | (Maps)       | (Slices)     | (Maps)      
----------------------------------------------------------------------------------------------------
JSON         | ✓            | ✓            | ✓            | ✓           
CBOR         | ✓            | ✓            | ✓            | ✓           
EasyJSON     | ✓            | ✓            | ✓            | ✓           
FlatBuffers  | ✓            | ✓            | ✗            | ✗           
Gob          | ✗            | ✓            | ✓            | ✓           
GoJSON       | ✓            | ✓            | ✓            | ✓           
JSONiter     | ✓            | ✓            | ✓            | ✓           
Msgp         | ✗            | ✓            | ✓            | ✗           
MsgPack      | ✓            | ✓            | ✓            | ✓           
Protobuf     | ✗            | ✗            | ✓            | ✗           
====================================================================================================

Benchmark completed successfully!
Results saved to: ./results
```

#### Output Contents

1. **Serialization Performance Results**
   - Data size (MB)
   - Marshal/Unmarshal speed (average and median)

2. **Marshal/Unmarshal Symmetry Test Results**
   - Type preservation for empty/nil slices and maps
   - ✓: Strict type preservation, ✗: Type conversion occurred

3. **Redis Performance Results** (if Redis measurements were performed)
   - SET/GET operation speed

### File Output

The following CSV files are saved in the `results/` directory:

- `serialization_results_YYYYMMDD_HHMMSS.csv` - Serialization performance
- `symmetry_results_YYYYMMDD_HHMMSS.csv` - Marshal/Unmarshal symmetry test results
- `redis_results_YYYYMMDD_HHMMSS.csv` - Redis performance (if executed)
