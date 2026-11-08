package main

import (
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// DAY 2: CHANNELS & CONCURRENCY COORDINATION IN GO
// =============================================================================
//
// Motto of Go Concurrency:
// "Do not communicate by sharing memory; instead, share memory by communicating."
//
// What is a Channel?
// A typed conduit through which you can send and receive values using the
// channel operator: `<-`
//
// Direction of Arrow (<-):
//   ch <- value   // SEND: Arrow points into the channel
//   value := <-ch // RECEIVE: Arrow points out of the channel
// =============================================================================

// -----------------------------------------------------------------------------
// 1. Basic Channel Worker
// -----------------------------------------------------------------------------
func basicWorker(ch chan int) {
	ch <- 100 // Send 100 into the channel
}

func demoBasicChannel() {
	fmt.Println("--- 1. Basic Channel ---")
	ch := make(chan int)

	go basicWorker(ch)

	result := <-ch // Blocks until worker sends value
	fmt.Println("Received from basicWorker:", result)
}

// -----------------------------------------------------------------------------
// 2. Unbuffered Channel (Synchronous / Blocking Rendezvous)
// -----------------------------------------------------------------------------
// - Created with: make(chan T)  (Capacity is 0)
// - Synchronous:
//   - Sender BLOCKS until Receiver is ready to read.
//   - Receiver BLOCKS until Sender provides a value.
func demoUnbufferedChannel() {
	fmt.Println("\n--- 2. Unbuffered Channel (Synchronous) ---")
	ch := make(chan string)

	go func() {
		fmt.Println("[Worker] Working on task...")
		time.Sleep(150 * time.Millisecond)

		fmt.Println("[Worker] Task done. Sending data into channel...")
		ch <- "Task Result: SUCCESS" // Blocks until receiver is ready
		fmt.Println("[Worker] Data was received!")
	}()

	fmt.Println("[Main] Waiting at channel receiver `<-ch`...")
	result := <-ch
	fmt.Println("[Main] Received from channel:", result)
}

// -----------------------------------------------------------------------------
// 3. Buffered Channel (Asynchronous / Non-blocking up to capacity)
// -----------------------------------------------------------------------------
// - Created with: make(chan T, capacity)
// - Asynchronous:
//   - Sender does NOT block until the buffer is FULL.
//   - Receiver does NOT block until the buffer is EMPTY.
func demoBufferedChannel() {
	fmt.Println("\n--- 3. Buffered Channel (Asynchronous) ---")

	bufCh := make(chan int, 3)

	// Send 3 items immediately without blocking
	bufCh <- 10
	bufCh <- 20
	bufCh <- 30

	fmt.Printf("Buffer length (items inside): %d | Capacity: %d\n", len(bufCh), cap(bufCh))

	// Read values out of the buffer (FIFO order)
	fmt.Println("Received item 1:", <-bufCh)
	fmt.Println("Received item 2:", <-bufCh)
	fmt.Println("Received item 3:", <-bufCh)
}

// -----------------------------------------------------------------------------
// 4. Closing Channels, Comma-OK Idiom & for-range
// -----------------------------------------------------------------------------
func demoCloseAndRange() {
	fmt.Println("\n--- 4. Closing Channels & for range Loop ---")

	numbers := make(chan int, 3)

	go func() {
		for i := 1; i <= 3; i++ {
			numbers <- i * 10
		}
		// ALWAYS close from the sender side, NEVER from receiver side
		close(numbers)
		fmt.Println("[Producer] Channel closed.")
	}()

	// Consumer reads until channel is closed
	for val := range numbers {
		fmt.Printf("   Received: %d\n", val)
	}

	// Comma-ok check on a closed channel
	val, ok := <-numbers
	if !ok {
		fmt.Println("[Consumer] Channel is closed. Default zero-value received:", val)
	}
}

// -----------------------------------------------------------------------------
// 5. Directional Channels (Send-Only vs Receive-Only)
// -----------------------------------------------------------------------------
func producer(out chan<- string) { // chan<- is send-only
	out <- "Data produced by producer"
	close(out)
}

func consumer(in <-chan string) { // <-chan is receive-only
	msg := <-in
	fmt.Println("[Directional Consumer] Received:", msg)
}

func demoDirectionalChannels() {
	fmt.Println("\n--- 5. Directional Channels (Type-Safe API) ---")
	pipe := make(chan string)

	go producer(pipe)
	consumer(pipe)
}

// -----------------------------------------------------------------------------
// 6. The `select` Statement (Multiplexing Channels)
// -----------------------------------------------------------------------------
// `select` lets a goroutine wait on MULTIPLE channel operations simultaneously.
// - It blocks until one of its cases is ready to execute.
// - If multiple cases are ready, it picks one uniformly at random.
// - A `default` case makes the select non-blocking.
// - Commonly used with `time.After()` for timeouts!
func demoSelect() {
	fmt.Println("\n--- 6. The select Statement ---")

	ch1 := make(chan string)
	ch2 := make(chan string)

	// Worker 1 finishes in 100ms
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "Response from Service A"
	}()

	// Worker 2 finishes in 200ms
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "Response from Service B"
	}()

	// Example A: Multiplexing between ch1, ch2, and a Timeout
	fmt.Println("Waiting on select with timeout...")
	select {
	case msg1 := <-ch1:
		fmt.Println("Received from ch1:", msg1)
	case msg2 := <-ch2:
		fmt.Println("Received from ch2:", msg2)
	case <-time.After(300 * time.Millisecond):
		fmt.Println("Timed out waiting for response!")
	}

	// Example B: Non-blocking channel read with default
	nonBlockCh := make(chan int)
	select {
	case val := <-nonBlockCh:
		fmt.Println("Received value:", val)
	default:
		fmt.Println("Non-blocking select: No value available immediately on nonBlockCh!")
	}
}

// -----------------------------------------------------------------------------
// 7. sync.WaitGroup (Coordinating Multiple Goroutines)
// -----------------------------------------------------------------------------
// Remember: WaitGroup waits for COMPLETION. It does NOT transfer data.
// Use Channels to transfer data, and WaitGroup to know when workers are done.
//
// Worker function receives a POINTER (*sync.WaitGroup) to avoid copying!
func waitGroupWorker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	// Rule: Usually use defer wg.Done() inside the worker so it always runs
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("   [Worker %d] Processing job %d...\n", id, job)
		time.Sleep(50 * time.Millisecond) // Simulate work
		results <- job * 2                // Send result through channel
	}
}

func demoWaitGroupWithChannels() {
	fmt.Println("\n--- 7. sync.WaitGroup with Channels ---")

	const numWorkers = 3
	const numJobs = 5

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	var wg sync.WaitGroup

	// Rule: Call Add() BEFORE starting the goroutine!
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go waitGroupWorker(w, jobs, results, &wg) // Pass *sync.WaitGroup as pointer!
	}

	// Send jobs into the jobs channel
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // No more jobs to send

	// Separate monitor goroutine: Wait for all workers to finish, then close results
	go func() {
		wg.Wait()      // Rule: Blocks until all workers call wg.Done()
		close(results) // Close results channel so consumer loop can exit
	}()

	// Collect and print all results
	fmt.Println("Collecting results:")
	for res := range results {
		fmt.Println("   Result received:", res)
	}
}

// =============================================================================
// Main Entry Point
// =============================================================================
func main() {
	fmt.Println("==================================================")
	fmt.Println("    GO CHANNELS, SELECT & WAITGROUP COMPLETE      ")
	fmt.Println("==================================================")

	// 1. Basic channel
	demoBasicChannel()

	// 2. Unbuffered / Synchronous
	demoUnbufferedChannel()

	// 3. Buffered / Asynchronous
	demoBufferedChannel()

	// 4. Closing & for-range / comma-ok
	demoCloseAndRange()

	// 5. Directional channels
	demoDirectionalChannels()

	// 6. Select statement (multiplexing & timeouts)
	demoSelect()

	// 7. sync.WaitGroup coordination
	demoWaitGroupWithChannels()

	fmt.Println("==================================================")
}

/*
================================================================================
Summary: Unbuffered vs Buffered Channels
================================================================================
| Feature         | Unbuffered (Sync)               | Buffered (Async)                  |
|-----------------|---------------------------------|-----------------------------------|
| Declaration     | make(chan T)                    | make(chan T, capacity)            |
| Capacity        | 0                               | > 0                                |
| Synchronization | Guaranteed immediate rendezvous | Decoupled up to buffer capacity   |
| Sender blocks?  | Yes, until receiver arrives     | Only when buffer is FULL          |
| Receiver blocks?| Yes, until sender arrives       | Only when buffer is EMPTY         |
| Primary Use     | Handshakes, signaling, locks    | Throughput, work queues, batches  |

Golden Rules of Go Channels:
1. Sending to a nil channel blocks forever.
2. Receiving from a nil channel blocks forever.
3. Sending to a CLOSED channel causes a PANIC!
4. Closing an already CLOSED channel causes a PANIC!
5. Receiving from a closed channel returns the zero value and ok = false immediately:
      value, ok := <-ch
      if !ok {
          fmt.Println("Channel closed")
      }

================================================================================
Important Production Rules for sync.WaitGroup
================================================================================
1. Call Add() BEFORE starting the goroutine:
   - If you call wg.Add(1) inside the goroutine, wg.Wait() in main may execute
     and exit before the goroutine even starts!
2. Usually use defer wg.Done() inside the worker:
   - Ensures wg.Done() is called even if the worker panics or returns early on error.
3. Pass *sync.WaitGroup (pointer), NOT a copied WaitGroup:
   - If passed by value, the function gets a copy with its own internal counter.
     Calling Done() on the copy will NEVER decrement the original!
4. Don't copy a WaitGroup after first use:
   - WaitGroup contains internal state and mutexes that cannot be copied.
5. Don't call Wait() before the relevant Add() operations have been established.
6. Key Distinction:
   - WaitGroup: WAITS FOR COMPLETION. Does not transfer data.
   - Channels:  TRANSFERS DATA and signals between goroutines.
================================================================================
*/
