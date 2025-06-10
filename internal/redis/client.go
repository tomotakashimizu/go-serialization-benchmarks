package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/models"
	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/serializers"
	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/utils"
)

// Client wraps Redis client with benchmark functionality
type Client struct {
	rdb *redis.Client
	ctx context.Context
}

// RedisResult contains Redis SET/GET performance results
type RedisResult struct {
	SerializerName string
	SetTimes       []int64 // nanoseconds
	GetTimes       []int64 // nanoseconds
	SetAvgNs       int64
	SetMedianNs    int64
	GetAvgNs       int64
	GetMedianNs    int64
}

// NewClient creates a new Redis client
func NewClient(addr, password string, db int) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Client{
		rdb: rdb,
		ctx: context.Background(),
	}
}

// Ping tests the connection to Redis
func (c *Client) Ping() error {
	_, err := c.rdb.Ping(c.ctx).Result()
	return err
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.rdb.Close()
}

// BenchmarkRedisOperations benchmarks Redis SET/GET operations for all serializers
func (c *Client) BenchmarkRedisOperations(serializers []serializers.Serializer, users []models.User, iterations int) ([]RedisResult, error) {
	results := make([]RedisResult, 0, len(serializers))

	for _, ser := range serializers {
		fmt.Printf("Running Redis benchmark for %s...\n", ser.Name())
		result, err := c.benchmarkSerializer(ser, users, iterations)
		if err != nil {
			return nil, fmt.Errorf("error benchmarking %s with Redis: %w", ser.Name(), err)
		}
		results = append(results, result)
	}

	return results, nil
}

// benchmarkSerializer benchmarks Redis operations for a single serializer
func (c *Client) benchmarkSerializer(ser serializers.Serializer, users []models.User, iterations int) (RedisResult, error) {
	result := RedisResult{
		SerializerName: ser.Name(),
		SetTimes:       make([]int64, iterations),
		GetTimes:       make([]int64, iterations),
	}

	// Serialize the users slice once
	data, err := ser.MarshalUsers(users)
	if err != nil {
		return result, fmt.Errorf("failed to marshal users: %w", err)
	}

	keyPrefix := fmt.Sprintf("benchmark:%s:users", ser.Name())

	// Run iterations
	for i := 0; i < iterations; i++ {
		key := fmt.Sprintf("%s:%d", keyPrefix, i)

		// Measure SET operation
		start := time.Now()
		err := c.rdb.Set(c.ctx, key, data, 0).Err()
		setTime := time.Since(start).Nanoseconds()
		if err != nil {
			return result, fmt.Errorf("SET operation failed: %w", err)
		}
		result.SetTimes[i] = setTime

		// Measure GET operation
		start = time.Now()
		retrievedData, err := c.rdb.Get(c.ctx, key).Bytes()
		getTime := time.Since(start).Nanoseconds()
		if err != nil {
			return result, fmt.Errorf("GET operation failed: %w", err)
		}
		result.GetTimes[i] = getTime

		// Verify data integrity
		_, err = ser.UnmarshalUsers(retrievedData)
		if err != nil {
			return result, fmt.Errorf("failed to unmarshal retrieved users data: %w", err)
		}

		// Clean up the key
		c.rdb.Del(c.ctx, key)
	}

	// Calculate statistics
	result.SetAvgNs = utils.CalculateAverage(result.SetTimes)
	result.SetMedianNs = utils.CalculateMedian(result.SetTimes)
	result.GetAvgNs = utils.CalculateAverage(result.GetTimes)
	result.GetMedianNs = utils.CalculateMedian(result.GetTimes)

	return result, nil
}

// CleanupTestKeys removes all test keys from Redis
func (c *Client) CleanupTestKeys() error {
	keys, err := c.rdb.Keys(c.ctx, "benchmark:*").Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.rdb.Del(c.ctx, keys...).Err()
	}

	return nil
}
