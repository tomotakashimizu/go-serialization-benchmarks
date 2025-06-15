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

	// Pure I/O times (Redis operations only)
	SetTimes    []int64 // nanoseconds
	GetTimes    []int64 // nanoseconds
	SetAvgNs    int64
	SetMedianNs int64
	GetAvgNs    int64
	GetMedianNs int64

	// Total times (including serialization)
	TotalSetTimes    []int64 // nanoseconds (marshal + SET)
	TotalGetTimes    []int64 // nanoseconds (GET + unmarshal)
	TotalSetAvgNs    int64
	TotalSetMedianNs int64
	TotalGetAvgNs    int64
	TotalGetMedianNs int64
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
		TotalSetTimes:  make([]int64, iterations),
		TotalGetTimes:  make([]int64, iterations),
	}

	keyPrefix := fmt.Sprintf("benchmark:%s:users", ser.Name())

	// Run iterations
	for i := 0; i < iterations; i++ {
		key := fmt.Sprintf("%s:%d", keyPrefix, i)

		// Measure total SET operation (marshal + SET)
		totalSetStart := time.Now()
		data, err := ser.MarshalUsers(users)
		if err != nil {
			return result, fmt.Errorf("failed to marshal users: %w", err)
		}

		// Measure pure SET operation
		setStart := time.Now()
		err = c.rdb.Set(c.ctx, key, data, 0).Err()
		setTime := time.Since(setStart).Nanoseconds()
		if err != nil {
			return result, fmt.Errorf("SET operation failed: %w", err)
		}
		totalSetTime := time.Since(totalSetStart).Nanoseconds()

		result.SetTimes[i] = setTime
		result.TotalSetTimes[i] = totalSetTime

		// Measure total GET operation (GET + unmarshal)
		totalGetStart := time.Now()

		// Measure pure GET operation
		getStart := time.Now()
		retrievedData, err := c.rdb.Get(c.ctx, key).Bytes()
		getTime := time.Since(getStart).Nanoseconds()
		if err != nil {
			return result, fmt.Errorf("GET operation failed: %w", err)
		}

		// Complete unmarshal for total time
		_, err = ser.UnmarshalUsers(retrievedData)
		if err != nil {
			return result, fmt.Errorf("failed to unmarshal retrieved users data: %w", err)
		}
		totalGetTime := time.Since(totalGetStart).Nanoseconds()

		result.GetTimes[i] = getTime
		result.TotalGetTimes[i] = totalGetTime

		// Clean up the key
		c.rdb.Del(c.ctx, key)
	}

	// Calculate statistics for pure I/O times
	result.SetAvgNs = utils.CalculateAverage(result.SetTimes)
	result.SetMedianNs = utils.CalculateMedian(result.SetTimes)
	result.GetAvgNs = utils.CalculateAverage(result.GetTimes)
	result.GetMedianNs = utils.CalculateMedian(result.GetTimes)

	// Calculate statistics for total times (including serialization)
	result.TotalSetAvgNs = utils.CalculateAverage(result.TotalSetTimes)
	result.TotalSetMedianNs = utils.CalculateMedian(result.TotalSetTimes)
	result.TotalGetAvgNs = utils.CalculateAverage(result.TotalGetTimes)
	result.TotalGetMedianNs = utils.CalculateMedian(result.TotalGetTimes)

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
