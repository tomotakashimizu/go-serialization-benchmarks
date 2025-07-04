package reporter

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/redis"
	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/serializers"
)

// Reporter handles reporting of benchmark results
type Reporter struct {
	outputDir string
}

// NewReporter creates a new reporter
func NewReporter(outputDir string) *Reporter {
	return &Reporter{
		outputDir: outputDir,
	}
}

// PrintSerializationResults prints serialization benchmark results to console
func (r *Reporter) PrintSerializationResults(results []serializers.SerializationResult) {
	fmt.Println("\n" + strings.Repeat("=", 120))
	fmt.Println("SERIALIZATION BENCHMARK RESULTS")
	fmt.Println(strings.Repeat("=", 120))

	// Header
	fmt.Printf("%-12s | %-12s | %-12s | %-12s | %-12s | %-12s\n",
		"Serializer", "Data Size", "Marshal Avg", "Marshal Med", "Unmarshal Avg", "Unmarshal Med")
	fmt.Printf("%-12s | %-12s | %-12s | %-12s | %-12s | %-12s\n",
		"", "(MB)", "(ms)", "(ms)", "(ms)", "(ms)")
	fmt.Println(strings.Repeat("-", 120))

	for _, result := range results {
		fmt.Printf("%-12s | %-12.2f | %-12.2f | %-12.2f | %-12.2f | %-12.2f\n",
			result.SerializerName,
			float64(result.DataSize)/1000000.0,
			float64(result.MarshalAvgNs)/1000000.0,
			float64(result.MarshalMedianNs)/1000000.0,
			float64(result.UnmarshalAvgNs)/1000000.0,
			float64(result.UnmarshalMedianNs)/1000000.0)
	}
	fmt.Println(strings.Repeat("=", 120))
}

// PrintSymmetryResults prints symmetry test results to console
func (r *Reporter) PrintSymmetryResults(results []serializers.SymmetryResult) {
	fmt.Println("\n" + strings.Repeat("=", 100))
	fmt.Println("STRICT TYPE PRESERVATION TEST RESULTS")
	fmt.Println(strings.Repeat("=", 100))

	fmt.Printf("%-12s | %-12s | %-12s | %-12s | %-12s\n",
		"Serializer", "Empty→Empty", "Empty{}→{}", "Nil→Nil", "Nil→Nil")
	fmt.Printf("%-12s | %-12s | %-12s | %-12s | %-12s\n",
		"", "(Slices)", "(Maps)", "(Slices)", "(Maps)")
	fmt.Println(strings.Repeat("-", 100))

	for _, result := range results {
		fmt.Printf("%-12s | %-12s | %-12s | %-12s | %-12s\n",
			result.SerializerName,
			boolToString(result.StrictEmptySlicesOK),
			boolToString(result.StrictEmptyMapsOK),
			boolToString(result.StrictNilSlicesOK),
			boolToString(result.StrictNilMapsOK))
	}

	fmt.Println(strings.Repeat("=", 100))

	// Print details
	fmt.Println("\nDetails:")
	for _, result := range results {
		fmt.Printf("%-12s: %s\n", result.SerializerName, result.Details)
	}
}

// PrintRedisResults prints Redis benchmark results to console
func (r *Reporter) PrintRedisResults(results []redis.RedisResult) {
	fmt.Println("\n" + strings.Repeat("=", 100))
	fmt.Println("REDIS PERFORMANCE RESULTS")
	fmt.Println(strings.Repeat("=", 100))

	// First table: Total time (including serialization)
	fmt.Println("Total Time (including serialization):")
	fmt.Printf("%-12s | %-12s | %-12s | %-12s | %-12s\n",
		"Serializer", "SET Avg", "SET Med", "GET Avg", "GET Med")
	fmt.Printf("%-12s | %-12s | %-12s | %-12s | %-12s\n",
		"", "(ms)", "(ms)", "(ms)", "(ms)")
	fmt.Println(strings.Repeat("-", 100))

	for _, result := range results {
		fmt.Printf("%-12s | %-12.2f | %-12.2f | %-12.2f | %-12.2f\n",
			result.SerializerName,
			float64(result.TotalSetAvgNs)/1000000.0,
			float64(result.TotalSetMedianNs)/1000000.0,
			float64(result.TotalGetAvgNs)/1000000.0,
			float64(result.TotalGetMedianNs)/1000000.0)
	}

	fmt.Println()
	fmt.Println("Pure I/O Time (Redis operations only):")
	fmt.Printf("%-12s | %-12s | %-12s | %-12s | %-12s\n",
		"Serializer", "SET Avg", "SET Med", "GET Avg", "GET Med")
	fmt.Printf("%-12s | %-12s | %-12s | %-12s | %-12s\n",
		"", "(ms)", "(ms)", "(ms)", "(ms)")
	fmt.Println(strings.Repeat("-", 100))

	for _, result := range results {
		fmt.Printf("%-12s | %-12.2f | %-12.2f | %-12.2f | %-12.2f\n",
			result.SerializerName,
			float64(result.SetAvgNs)/1000000.0,
			float64(result.SetMedianNs)/1000000.0,
			float64(result.GetAvgNs)/1000000.0,
			float64(result.GetMedianNs)/1000000.0)
	}
	fmt.Println(strings.Repeat("=", 100))
}

// SaveSerializationResults saves serialization results to CSV
func (r *Reporter) SaveSerializationResults(results []serializers.SerializationResult) error {
	filename := fmt.Sprintf("serialization_results_%s.csv", time.Now().Format("20060102_150405"))
	filepath := filepath.Join(r.outputDir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Serializer", "DataSize_Bytes", "DataSize_MB", "MarshalAvg_ns", "MarshalMedian_ns",
		"UnmarshalAvg_ns", "UnmarshalMedian_ns", "MarshalAvg_ms", "MarshalMedian_ms",
		"UnmarshalAvg_ms", "UnmarshalMedian_ms",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data
	for _, result := range results {
		record := []string{
			result.SerializerName,
			strconv.Itoa(result.DataSize),
			fmt.Sprintf("%.2f", float64(result.DataSize)/1000000.0),
			strconv.FormatInt(result.MarshalAvgNs, 10),
			strconv.FormatInt(result.MarshalMedianNs, 10),
			strconv.FormatInt(result.UnmarshalAvgNs, 10),
			strconv.FormatInt(result.UnmarshalMedianNs, 10),
			fmt.Sprintf("%.2f", float64(result.MarshalAvgNs)/1000000.0),
			fmt.Sprintf("%.2f", float64(result.MarshalMedianNs)/1000000.0),
			fmt.Sprintf("%.2f", float64(result.UnmarshalAvgNs)/1000000.0),
			fmt.Sprintf("%.2f", float64(result.UnmarshalMedianNs)/1000000.0),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	fmt.Printf("Serialization results saved to: %s\n", filepath)
	return nil
}

// SaveSymmetryResults saves symmetry results to CSV
func (r *Reporter) SaveSymmetryResults(results []serializers.SymmetryResult) error {
	filename := fmt.Sprintf("symmetry_results_%s.csv", time.Now().Format("20060102_150405"))
	filepath := filepath.Join(r.outputDir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Serializer", "StrictEmptySlicesOK", "StrictEmptyMapsOK", "StrictNilSlicesOK", "StrictNilMapsOK", "Details",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data
	for _, result := range results {
		record := []string{
			result.SerializerName,
			boolToString(result.StrictEmptySlicesOK),
			boolToString(result.StrictEmptyMapsOK),
			boolToString(result.StrictNilSlicesOK),
			boolToString(result.StrictNilMapsOK),
			result.Details,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	fmt.Printf("Symmetry results saved to: %s\n", filepath)
	return nil
}

// SaveRedisResults saves Redis results to CSV
func (r *Reporter) SaveRedisResults(results []redis.RedisResult) error {
	filename := fmt.Sprintf("redis_results_%s.csv", time.Now().Format("20060102_150405"))
	filepath := filepath.Join(r.outputDir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Serializer",
		"TotalSetAvg_ns", "TotalSetMedian_ns", "TotalGetAvg_ns", "TotalGetMedian_ns",
		"TotalSetAvg_ms", "TotalSetMedian_ms", "TotalGetAvg_ms", "TotalGetMedian_ms",
		"IOSetAvg_ns", "IOSetMedian_ns", "IOGetAvg_ns", "IOGetMedian_ns",
		"IOSetAvg_ms", "IOSetMedian_ms", "IOGetAvg_ms", "IOGetMedian_ms",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data
	for _, result := range results {
		record := []string{
			result.SerializerName,
			// Total times (including serialization)
			strconv.FormatInt(result.TotalSetAvgNs, 10),
			strconv.FormatInt(result.TotalSetMedianNs, 10),
			strconv.FormatInt(result.TotalGetAvgNs, 10),
			strconv.FormatInt(result.TotalGetMedianNs, 10),
			fmt.Sprintf("%.2f", float64(result.TotalSetAvgNs)/1000000.0),
			fmt.Sprintf("%.2f", float64(result.TotalSetMedianNs)/1000000.0),
			fmt.Sprintf("%.2f", float64(result.TotalGetAvgNs)/1000000.0),
			fmt.Sprintf("%.2f", float64(result.TotalGetMedianNs)/1000000.0),
			// Pure I/O times (Redis operations only)
			strconv.FormatInt(result.SetAvgNs, 10),
			strconv.FormatInt(result.SetMedianNs, 10),
			strconv.FormatInt(result.GetAvgNs, 10),
			strconv.FormatInt(result.GetMedianNs, 10),
			fmt.Sprintf("%.2f", float64(result.SetAvgNs)/1000000.0),
			fmt.Sprintf("%.2f", float64(result.SetMedianNs)/1000000.0),
			fmt.Sprintf("%.2f", float64(result.GetAvgNs)/1000000.0),
			fmt.Sprintf("%.2f", float64(result.GetMedianNs)/1000000.0),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	fmt.Printf("Redis results saved to: %s\n", filepath)
	return nil
}

// EnsureOutputDir creates the output directory if it doesn't exist
func (r *Reporter) EnsureOutputDir() error {
	return os.MkdirAll(r.outputDir, 0755)
}

// boolToString converts boolean to string representation
func boolToString(b bool) string {
	if b {
		return "✓"
	}
	return "✗"
}
