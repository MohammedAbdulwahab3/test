package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Result struct {
	Duration time.Duration
	Status   int
	Error    error
}

func main() {
	url := flag.String("url", "http://100.104.105.207:8080/api/items", "Target URL")
	concurrency := flag.Int("c", 50, "Number of concurrent workers")
	totalRequests := flag.Int("n", 10000, "Total number of requests")
	flag.Parse()

	fmt.Printf("🚀 Starting load test on %s\n", *url)
	fmt.Printf("   Workers: %d\n", *concurrency)
	fmt.Printf("   Requests: %d\n\n", *totalRequests)

	results := make(chan Result, *totalRequests)
	var wg sync.WaitGroup

	startTime := time.Now()
	requestsPerWorker := *totalRequests / *concurrency

	// Start workers
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{
				Timeout: 10 * time.Second,
				Transport: &http.Transport{
					MaxIdleConns:        100,
					MaxIdleConnsPerHost: 100,
				},
			}

			for j := 0; j < requestsPerWorker; j++ {
				start := time.Now()
				resp, err := client.Get(*url)
				duration := time.Since(start)

				status := 0
				if err == nil {
					status = resp.StatusCode
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}

				results <- Result{
					Duration: duration,
					Status:   status,
					Error:    err,
				}
			}
		}()
	}

	// Wait for completion
	wg.Wait()
	close(results)
	totalDuration := time.Since(startTime)

	// Process results
	var latencies []time.Duration
	var successCount, errorCount int64
	statusCodes := make(map[int]int)

	for res := range results {
		if res.Error != nil || res.Status >= 400 {
			errorCount++
		} else {
			successCount++
			latencies = append(latencies, res.Duration)
		}
		if res.Status > 0 {
			statusCodes[res.Status]++
		}
	}

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	// Calculate metrics
	rps := float64(*totalRequests) / totalDuration.Seconds()
	avgLatency := time.Duration(0)
	if len(latencies) > 0 {
		var totalLatency time.Duration
		for _, l := range latencies {
			totalLatency += l
		}
		avgLatency = totalLatency / time.Duration(len(latencies))
	}

	// Print Report
	fmt.Println("📊 Load Test Results")
	fmt.Println("====================")
	fmt.Printf("Total Duration:   %v\n", totalDuration)
	fmt.Printf("Total Requests:   %d\n", *totalRequests)
	fmt.Printf("Success Rate:     %.2f%%\n", float64(successCount)/float64(*totalRequests)*100)
	fmt.Printf("Requests/Sec:     %.2f req/s\n", rps)
	fmt.Println("\nLatency Distribution:")
	fmt.Printf("  Average:   %v\n", avgLatency)
	if len(latencies) > 0 {
		fmt.Printf("  P50 (Med): %v\n", latencies[len(latencies)*50/100])
		fmt.Printf("  P95:       %v\n", latencies[len(latencies)*95/100])
		fmt.Printf("  P99:       %v\n", latencies[len(latencies)*99/100])
		fmt.Printf("  Max:       %v\n", latencies[len(latencies)-1])
	}

	fmt.Println("\nStatus Codes:")
	for code, count := range statusCodes {
		fmt.Printf("  [%d]: %d\n", code, count)
	}

	if errorCount > 0 {
		fmt.Printf("\n⚠️  Errors: %d\n", errorCount)
	}
}
