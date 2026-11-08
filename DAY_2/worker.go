package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// DAY 2: THE WORKER POOL & FAN-OUT / FAN-IN PATTERNS (PRODUCTION CONCURRENCY)
// =============================================================================
//
// Part A: Why Worker Pools?
// Imagine you have 10,000 tasks (e.g. sending emails, resizing images, API calls).
// If you spawn 10,000 goroutines all at once:
//   - You could exhaust OS file descriptors or database connection pools.
//   - You could get rate-limited by external APIs.
//   - Memory usage spikes and CPU cache thrashing occurs.
//
// A WORKER POOL bounds concurrency to a fixed, manageable number of workers.
//
// Architecture:
//
//   [ JOB PRODUCER ]
//          |
//          v
//     Jobs Channel (buffered)
//    +-----+-----+-----+
//    |     |     |     |
//    v     v     v     v
// [Worker 1] [Worker 2] [Worker 3] ... [Worker N]
//    |     |     |     |
//    +-----+-----+-----+
//          v
//    Results Channel (buffered)
//          |
//          v
//     [ COLLECTOR ]
// =============================================================================

// Job represents an incoming unit of work.
type Job struct {
	ID    int
	Data  string
	Delay time.Duration // Simulated processing time
}

// Result represents the outcome of a processed job, including error propagation.
type Result struct {
	JobID  int
	Output string
	Err    error // Production rule: Always propagate errors!
}

// -----------------------------------------------------------------------------
// The Worker Function
// -----------------------------------------------------------------------------
func worker(
	id int,
	jobs <-chan Job,
	results chan<- Result,
	wg *sync.WaitGroup,
	ctx context.Context,
) {
	// Rule 1: Always defer wg.Done() inside the worker
	defer wg.Done()

	for {
		select {
		// Production Feature: Graceful shutdown support via Context
		case <-ctx.Done():
			fmt.Printf("   [Worker %d] Received shutdown signal. Exiting...\n", id)
			return

		case job, ok := <-jobs:
			if !ok {
				// Jobs channel is closed and drained -> worker exits cleanly
				return
			}

			// Process the job
			fmt.Printf("   [Worker %d] Started job #%d (%s)\n", id, job.ID, job.Data)
			time.Sleep(job.Delay) // Simulate workload

			// Simulate error condition for demonstration
			var err error
			if job.ID == 7 {
				err = fmt.Errorf("simulated network failure on job #%d", job.ID)
			}

			// Send outcome to results channel
			results <- Result{
				JobID:  job.ID,
				Output: fmt.Sprintf("PROCESSED: %s", job.Data),
				Err:    err,
			}

			fmt.Printf("   [Worker %d] Finished job #%d\n", id, job.ID)
		}
	}
}

// =============================================================================
// Part B: Fan-Out / Fan-In Pattern Implementation
// =============================================================================
// - Fan-Out: Multiple goroutines read from the SAME channel to divide and conquer work.
// - Fan-In:  A multiplexer merges MULTIPLE channels into a SINGLE output channel.
//
// Pipeline:
//                     +-> [Worker 1] (ch1) -+
//  [Producer] -> In --+                     +--> [Fan-In / Merge] -> Combined Out
//                     +-> [Worker 2] (ch2) -+
// -----------------------------------------------------------------------------

// Fan-Out Worker: Reads from input channel and returns its own output channel
func squareWorker(id int, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			time.Sleep(30 * time.Millisecond) // Simulate processing
			fmt.Printf("   [Fan-Out Worker %d] Squaring %d => %d\n", id, n, n*n)
			out <- n * n
		}
	}()
	return out
}

// Fan-In Multiplexer: Merges multiple input channels into a single output channel
func merge(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	// Helper function that drains one channel into the merged output
	multiplex := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	// Close out when all incoming channels are fully drained
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func demoFanOutFanIn() {
	fmt.Println("\n==================================================")
	fmt.Println("             FAN-OUT / FAN-IN DEMO                ")
	fmt.Println("==================================================")

	// Step 1: Input source channel
	in := make(chan int)
	go func() {
		for i := 1; i <= 6; i++ {
			in <- i
		}
		close(in)
	}()

	// Step 2: FAN-OUT
	// Launch 2 worker goroutines consuming from the SAME `in` channel simultaneously
	worker1 := squareWorker(1, in)
	worker2 := squareWorker(2, in)

	// Step 3: FAN-IN
	// Merge both worker output channels into a single unified stream
	merged := merge(worker1, worker2)

	// Step 4: Consume merged stream
	fmt.Println("[Fan-In] Reading merged stream:")
	for val := range merged {
		fmt.Printf("   Merged output: %d\n", val)
	}
}

// =============================================================================
// Main Entry Point
// =============================================================================
func main() {
	fmt.Println("==================================================")
	fmt.Println("       PRODUCTION WORKER POOL PATTERN             ")
	fmt.Println("==================================================")

	const numWorkers = 3
	const numJobs = 10

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Buffered channels: Prevents producer and collector from blocking on every item
	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup

	// --------------------------------------------------
	// Step 1: Start the fixed pool of Workers
	// --------------------------------------------------
	fmt.Printf("[Init] Starting %d workers...\n", numWorkers)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1) // Rule 2: Call Add() BEFORE starting goroutine!
		go worker(w, jobs, results, &wg, ctx)
	}

	// --------------------------------------------------
	// Step 2: Producer - Enqueue Jobs
	// --------------------------------------------------
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- Job{
				ID:    j,
				Data:  fmt.Sprintf("Task-Payload-%d", j),
				Delay: 40 * time.Millisecond,
			}
		}
		// Rule 3: Close jobs channel when all work has been submitted!
		close(jobs)
		fmt.Println("[Producer] All jobs submitted. Jobs channel closed.")
	}()

	// --------------------------------------------------
	// Step 3: Monitor - Close Results Channel when Workers Finish
	// --------------------------------------------------
	go func() {
		wg.Wait()
		close(results)
		fmt.Println("[Monitor] All workers completed. Results channel closed.")
	}()

	// --------------------------------------------------
	// Step 4: Collector - Process & Handle Results
	// --------------------------------------------------
	fmt.Println("\n[Collector] Reading results:")
	successCount := 0
	failureCount := 0

	for res := range results {
		if res.Err != nil {
			fmt.Printf("   [FAIL] Job #%d failed: %v\n", res.JobID, res.Err)
			failureCount++
		} else {
			fmt.Printf("   [OK]   Job #%d => %s\n", res.JobID, res.Output)
			successCount++
		}
	}

	fmt.Println("--------------------------------------------------")
	fmt.Printf("Summary: Total: %d | Success: %d | Failed: %d\n", numJobs, successCount, failureCount)

	// --------------------------------------------------
	// Part B Demo: Run Fan-Out / Fan-In
	// --------------------------------------------------
	demoFanOutFanIn()
}

/*
================================================================================
Real-World Applications of Worker Pools
================================================================================
| Application          | Worker Task                                            |
|----------------------|--------------------------------------------------------|
| Email Service        | Send queued marketing/transactional emails             |
| Image Processing     | Resize, compress, and upload images to S3              |
| Data Pipelines       | Ingest, transform, and parse Kafka/SQS event records   |
| Web Scraping         | Fetch and parse multiple web pages with rate limits    |
| Notification System  | Deliver push notifications / SMS to mobile users       |
| File Processing      | Parse uploaded large CSV/Excel files in chunks         |
| Background Jobs      | Execute periodic maintenance and scheduled tasks       |

================================================================================
Fan-Out / Fan-In Pattern & Real Production Use Cases
================================================================================
1. Definitions:
   - Fan-Out: Multiple goroutines read from the SAME channel to divide work
     concurrently across multiple CPU cores.
   - Fan-In:  A multiplexer merges MULTIPLE channels into a SINGLE output
     channel so a single consumer can process all incoming data in order.

2. Production Use Cases:
   | Production System         | How Fan-Out / Fan-In is Applied                |
   |---------------------------|------------------------------------------------|
   | Microservice Aggregator   | Fan-Out to fetch user profile, order history,  |
   | (API Gateway / BFF)       | and recommendations concurrently; Fan-In to    |
   |                           | merge all 3 into one response payload.         |
   |---------------------------|------------------------------------------------|
   | Log Ingestion / Analytics | Fan-Out to multiple parser workers for incoming|
   | (ELK / ClickHouse)        | log streams; Fan-In to a single batch-writer   |
   |                           | to persist records efficiently.                |
   |---------------------------|------------------------------------------------|
   | Web Crawler / Search      | Fan-Out to fetch URLs in parallel; Fan-In      |
   | Indexer                   | discovered links into a single deduplication   |
   |                           | queue.                                         |
   |---------------------------|------------------------------------------------|
   | Image / Video Transcoding | Fan-Out to transcode individual video segments;|
   |                           | Fan-In transcoded chunks to stitch final file. |

================================================================================
Production Engineering Checklist
================================================================================
1. Worker Pool Sizing Rule:
   - CPU-bound tasks (encryption, compression): numWorkers = runtime.NumCPU()
   - I/O-bound tasks (DB, network requests):   numWorkers = 10 to 100+

2. Backpressure via Channel Buffers:
   - Bounded channels prevent memory exhaustion by applying backpressure when
     workers cannot keep up with producers.

3. Always Propagate Errors:
   - Include `Err error` in your result struct so the consumer can retry or log.

4. Graceful Shutdown & Context:
   - Always wire `context.Context` to stop workers cleanly on SIGINT / SIGTERM.
================================================================================
*/
