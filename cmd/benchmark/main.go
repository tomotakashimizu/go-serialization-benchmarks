package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/benchmark"
	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/models"
	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/redis"
	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/reporter"
	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/serializers"
)

func main() {
	// Command line flags
	var (
		dataCount     = flag.Int("count", 100000, "Number of test records to generate")
		iterations    = flag.Int("iterations", 5, "Number of benchmark iterations")
		redisAddr     = flag.String("redis-addr", "localhost:6379", "Redis server address")
		redisPassword = flag.String("redis-password", "", "Redis password")
		redisDB       = flag.Int("redis-db", 0, "Redis database number")
		outputDir     = flag.String("output", "./results", "Output directory for results")
		skipRedis     = flag.Bool("skip-redis", false, "Skip Redis benchmarks")
		help          = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	fmt.Printf("Serializer Performance Benchmark\n")
	fmt.Printf("=================================\n")
	fmt.Printf("Test data count: %d\n", *dataCount)
	fmt.Printf("Benchmark iterations: %d\n", *iterations)
	fmt.Printf("Output directory: %s\n", *outputDir)
	fmt.Printf("Redis: %s (skip: %t)\n\n", *redisAddr, *skipRedis)

	// Initialize reporter
	rep := reporter.NewReporter(*outputDir)
	if err := rep.EnsureOutputDir(); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Generate test data
	fmt.Printf("Generating %d test records...\n", *dataCount)
	users := models.GenerateTestUsers(*dataCount)
	fmt.Printf("Test data generated successfully.\n\n")

	// Initialize benchmark runner
	runner := benchmark.NewRunner()
	runner.SetTestData(users)

	// Add all serializers (JSON first, then alphabetical order)
	runner.AddSerializer(serializers.NewJSONSerializer()) // Most common format first
	runner.AddSerializer(serializers.NewCBORSerializer())
	runner.AddSerializer(serializers.NewEasyJSONSerializer())
	runner.AddSerializer(serializers.NewGobSerializer())
	runner.AddSerializer(serializers.NewGoJSONSerializer())
	runner.AddSerializer(serializers.NewJSONiterSerializer())
	runner.AddSerializer(serializers.NewMsgPackSerializer())
	runner.AddSerializer(serializers.NewProtobufSerializer())

	// Run serialization benchmarks
	fmt.Println("Running serialization benchmarks...")
	serializationResults, err := runner.RunBenchmarks(*iterations)
	if err != nil {
		log.Fatalf("Serialization benchmark failed: %v", err)
	}

	// Print and save serialization results
	rep.PrintSerializationResults(serializationResults)
	if err := rep.SaveSerializationResults(serializationResults); err != nil {
		log.Printf("Failed to save serialization results: %v", err)
	}

	// Run symmetry tests
	fmt.Println("\nRunning symmetry tests...")
	symmetryResults, err := runner.RunSymmetryTests()
	if err != nil {
		log.Fatalf("Symmetry test failed: %v", err)
	}

	// Print and save symmetry results
	rep.PrintSymmetryResults(symmetryResults)
	if err := rep.SaveSymmetryResults(symmetryResults); err != nil {
		log.Printf("Failed to save symmetry results: %v", err)
	}

	// Run Redis benchmarks if not skipped
	if !*skipRedis {
		fmt.Println("\nRunning Redis benchmarks...")
		redisClient := redis.NewClient(*redisAddr, *redisPassword, *redisDB)
		defer redisClient.Close()

		// Test Redis connection
		if err := redisClient.Ping(); err != nil {
			log.Printf("Warning: Redis connection failed (%v). Skipping Redis benchmarks.", err)
		} else {
			// Cleanup any existing test keys
			if err := redisClient.CleanupTestKeys(); err != nil {
				log.Printf("Warning: Failed to cleanup Redis test keys: %v", err)
			}

			// Use all users for Redis benchmarks
			// Create serializers for Redis test (JSON first, then alphabetical order)
			redisSerializers := []serializers.Serializer{
				serializers.NewJSONSerializer(), // Most common format first
				serializers.NewCBORSerializer(),
				serializers.NewEasyJSONSerializer(),
				serializers.NewGobSerializer(),
				serializers.NewGoJSONSerializer(),
				serializers.NewJSONiterSerializer(),
				serializers.NewMsgPackSerializer(),
				serializers.NewProtobufSerializer(),
			}

			redisResults, err := redisClient.BenchmarkRedisOperations(redisSerializers, users, *iterations)
			if err != nil {
				log.Printf("Redis benchmark failed: %v", err)
			} else {
				// Print and save Redis results
				rep.PrintRedisResults(redisResults)
				if err := rep.SaveRedisResults(redisResults); err != nil {
					log.Printf("Failed to save Redis results: %v", err)
				}
			}

			// Cleanup test keys
			if err := redisClient.CleanupTestKeys(); err != nil {
				log.Printf("Warning: Failed to cleanup Redis test keys: %v", err)
			}
		}
	}

	fmt.Printf("\nBenchmark completed successfully!\n")
	fmt.Printf("Results saved to: %s\n", *outputDir)
}

func showHelp() {
	fmt.Printf("Serializer Performance Benchmark Tool\n")
	fmt.Printf("=====================================\n\n")
	fmt.Printf("This tool compares the performance of different serialization formats:\n")
	fmt.Printf("- JSON (standard library)\n")
	fmt.Printf("- CBOR (github.com/fxamacker/cbor/v2)\n")
	fmt.Printf("- EasyJSON (github.com/mailru/easyjson - high-performance JSON with code generation)\n")
	fmt.Printf("- Gob (standard library)\n")
	fmt.Printf("- GoJSON (github.com/goccy/go-json - high-performance JSON)\n")
	fmt.Printf("- JSONiter (github.com/json-iterator/go - high-performance JSON)\n")
	fmt.Printf("- MessagePack (github.com/vmihailenco/msgpack/v5)\n")
	fmt.Printf("- Protocol Buffers (google.golang.org/protobuf)\n\n")

	fmt.Printf("The benchmark measures:\n")
	fmt.Printf("1. Serialization/deserialization speed (average & median)\n")
	fmt.Printf("2. Data size in bytes\n")
	fmt.Printf("3. Marshal/Unmarshal symmetry for empty/nil slices and maps\n")
	fmt.Printf("4. Redis SET/GET performance (optional)\n\n")

	fmt.Printf("Usage:\n")
	fmt.Printf("  %s [options]\n\n", os.Args[0])

	fmt.Printf("Options:\n")
	flag.PrintDefaults()

	fmt.Printf("\nExamples:\n")
	fmt.Printf("  # Run with default settings (100k records, 5 iterations)\n")
	fmt.Printf("  %s\n\n", os.Args[0])

	fmt.Printf("  # Run with 10k records and skip Redis tests\n")
	fmt.Printf("  %s -count=10000 -skip-redis\n\n", os.Args[0])

	fmt.Printf("  # Run with custom Redis settings\n")
	fmt.Printf("  %s -redis-addr=192.168.1.100:6379 -redis-password=secret\n\n", os.Args[0])
}
