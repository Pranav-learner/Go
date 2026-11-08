package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

// =============================================================================
// DAY 2: ASSIGNMENT 3 — COMMUNICATION & WORKER POOL
// CONCURRENT URL CHECKER WITH CONTEXT TIMEOUTS
// =============================================================================
// Requirements:
//   - Accept 20 URLs.
//   - Distribute them to 5 workers (Worker Pool).
//   - Fetch URLs concurrently.
//   - Collect HTTP status codes.
//   - Handle failures (non-200, connection errors, timeouts).
//   - Apply a timeout using context.Context.
// Concepts:
//   - Worker pools, channels, context, error handling, WaitGroups.
// =============================================================================

// CheckResult represents the outcome of checking a single URL
type CheckResult struct {
	WorkerID   int
	URL        string
	StatusCode int
	Duration   time.Duration
	Err        error
}

// urlWorker pulls URLs from the channel and performs HTTP GET requests with context timeouts
func urlWorker(
	id int,
	urls <-chan string,
	results chan<- CheckResult,
	wg *sync.WaitGroup,
	requestTimeout time.Duration,
) {
	defer wg.Done()

	client := &http.Client{}

	for targetURL := range urls {
		start := time.Now()

		// 1. Apply a timeout per URL using context.WithTimeout
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

		// Create HTTP Request attached to context
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
		if err != nil {
			cancel()
			results <- CheckResult{
				WorkerID: id,
				URL:      targetURL,
				Duration: time.Since(start),
				Err:      err,
			}
			continue
		}

		// Perform HTTP Request
		resp, err := client.Do(req)
		elapsed := time.Since(start)
		cancel() // Release timer resources immediately after request returns

		if err != nil {
			results <- CheckResult{
				WorkerID: id,
				URL:      targetURL,
				Duration: elapsed,
				Err:      err, // Captures context.DeadlineExceeded or network errors
			}
			continue
		}

		// Close response body to prevent resource leaks
		_ = resp.Body.Close()

		results <- CheckResult{
			WorkerID:   id,
			URL:        targetURL,
			StatusCode: resp.StatusCode,
			Duration:   elapsed,
			Err:        nil,
		}
	}
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("     ASSIGNMENT 3: CONCURRENT URL CHECKER         ")
	fmt.Println("==================================================")

	const (
		numWorkers     = 5
		requestTimeout = 300 * time.Millisecond // 300ms timeout per request
	)

	// -------------------------------------------------------------
	// Setup Local Mock Test Server
	// (Simulates various HTTP responses: 200, 404, 500, and Slow Timeout)
	// -------------------------------------------------------------
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/slow":
			// Simulates a slow server taking 600ms (will exceed 300ms timeout!)
			time.Sleep(600 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		case "/notfound":
			w.WriteHeader(http.StatusNotFound)
		case "/server-error":
			w.WriteHeader(http.StatusInternalServerError)
		case "/unauthorized":
			w.WriteHeader(http.StatusUnauthorized)
		default:
			// Normal 200 OK endpoint with a fast 10-30ms response
			time.Sleep(15 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	// -------------------------------------------------------------
	// Generate 20 test URLs
	// -------------------------------------------------------------
	testURLs := []string{
		server.URL + "/api/v1/users",
		server.URL + "/api/v1/products",
		server.URL + "/slow", // Will trigger context timeout!
		server.URL + "/api/v1/orders",
		server.URL + "/notfound", // 404 Not Found
		server.URL + "/api/v1/health",
		server.URL + "/server-error", // 500 Server Error
		server.URL + "/api/v1/settings",
		server.URL + "/unauthorized", // 401 Unauthorized
		server.URL + "/slow",         // Will trigger context timeout!
		server.URL + "/api/v1/cart",
		server.URL + "/api/v1/checkout",
		server.URL + "/notfound", // 404 Not Found
		server.URL + "/api/v1/catalog",
		server.URL + "/api/v1/inventory",
		server.URL + "/api/v1/notifications",
		server.URL + "/api/v1/search",
		"http://invalid-unreachable-domain-xyz.local", // Connection failure
		server.URL + "/api/v1/reviews",
		server.URL + "/api/v1/billing",
	}

	numURLs := len(testURLs)
	fmt.Printf("[Config] URLs to check: %d | Worker pool size: %d | Timeout: %v\n\n", numURLs, numWorkers, requestTimeout)

	// Buffered channels
	urlChan := make(chan string, numURLs)
	resultChan := make(chan CheckResult, numURLs)

	var wg sync.WaitGroup

	startTime := time.Now()

	// -------------------------------------------------------------
	// 1. Launch 5 Workers in the Worker Pool
	// -------------------------------------------------------------
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go urlWorker(w, urlChan, resultChan, &wg, requestTimeout)
	}

	// -------------------------------------------------------------
	// 2. Feed 20 URLs into the urlChan
	// -------------------------------------------------------------
	for _, u := range testURLs {
		urlChan <- u
	}
	close(urlChan) // No more URLs to send; allows workers to finish their loops

	// -------------------------------------------------------------
	// 3. Monitor goroutine: Wait for all workers, then close resultChan
	// -------------------------------------------------------------
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// -------------------------------------------------------------
	// 4. Collect and Print Results
	// -------------------------------------------------------------
	successCount := 0
	clientErrorCount := 0
	serverErrorCount := 0
	failureCount := 0

	fmt.Printf("%-8s %-4s %-8s %-45s %s\n", "WORKER", "CODE", "TIME", "URL", "STATUS / ERROR")
	fmt.Println("----------------------------------------------------------------------------------------")

	for res := range resultChan {
		if res.Err != nil {
			failureCount++
			fmt.Printf("[W#%d]   %-4s %-8v %-45s FAILED: %v\n",
				res.WorkerID, "-", res.Duration.Round(time.Millisecond), res.URL, res.Err)
		} else {
			switch {
			case res.StatusCode >= 200 && res.StatusCode < 300:
				successCount++
				fmt.Printf("[W#%d]   %-4d %-8v %-45s OK\n",
					res.WorkerID, res.StatusCode, res.Duration.Round(time.Millisecond), res.URL)
			case res.StatusCode >= 400 && res.StatusCode < 500:
				clientErrorCount++
				fmt.Printf("[W#%d]   %-4d %-8v %-45s CLIENT ERROR\n",
					res.WorkerID, res.StatusCode, res.Duration.Round(time.Millisecond), res.URL)
			case res.StatusCode >= 500:
				serverErrorCount++
				fmt.Printf("[W#%d]   %-4d %-8v %-45s SERVER ERROR\n",
					res.WorkerID, res.StatusCode, res.Duration.Round(time.Millisecond), res.URL)
			}
		}
	}

	totalTime := time.Since(startTime)

	// -------------------------------------------------------------
	// 5. Summary Report
	// -------------------------------------------------------------
	fmt.Println("========================================================================================")
	fmt.Printf("Checked %d URLs across %d workers in %v\n", numURLs, numWorkers, totalTime.Round(time.Millisecond))
	fmt.Printf("Summary: 2xx Success: %d | 4xx Client Error: %d | 5xx Server Error: %d | Failures/Timeouts: %d\n",
		successCount, clientErrorCount, serverErrorCount, failureCount)
	fmt.Println("========================================================================================")
}
