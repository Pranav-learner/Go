package main

import (
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// DAY 2: MUTEX, RWMUTEX, RACE CONDITIONS & RACE DETECTOR
// =============================================================================
//
// 1. What is a Race Condition?
//    A race condition occurs when two or more goroutines access the same memory
//    location concurrently, and AT LEAST ONE of the accesses is a WRITE.
//    (e.g., `counter++` is NOT atomic — it involves: Read -> Modify -> Write).
//
// 2. The Go Race Detector:
//    Go has a built-in data race detector!
//    Run with:
//        go run -race filename.go
//        go test -race ./...
// =============================================================================

// -----------------------------------------------------------------------------
// 1. Unsynchronized Counter (Shows Race Condition Concept)
// -----------------------------------------------------------------------------
// Note: If you run unsynchronized concurrent increments, `go run -race`
// will immediately catch and flag a DATA RACE warning!
func demoRaceConditionConcept() {
	fmt.Println("--- 1. Why We Need Locks (The Counter Problem) ---")

	// Notice: 1000 goroutines trying to increment the same number without a lock
	// will lose updates due to concurrent read/write interleaving.
	fmt.Println("Without synchronization, `count++` loses updates and causes race conditions.")
	fmt.Println("Solution: Use sync.Mutex or sync.RWMutex!")
}

// -----------------------------------------------------------------------------
// 2. sync.Mutex (Mutual Exclusion Lock: 1 Writer OR 1 Reader)
// -----------------------------------------------------------------------------
// SafeCounter wraps an integer with a sync.Mutex to protect concurrent access.
type SafeCounter struct {
	mu    sync.Mutex // Mutex guards the value field below
	value int
}

// Inc increments the counter safely.
// Rule: Use pointer receiver so we don't copy the mutex!
func (c *SafeCounter) Inc() {
	c.mu.Lock()         // Acquire exclusive lock
	defer c.mu.Unlock() // Always defer unlock right after locking
	c.value++
}

// Value returns the current value safely.
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func demoMutex() {
	fmt.Println("\n--- 2. sync.Mutex (Exclusive Lock) ---")

	counter := SafeCounter{}
	var wg sync.WaitGroup

	const numGoroutines = 1000

	// Launch 1000 concurrent goroutines safely incrementing counter
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Inc()
		}()
	}

	wg.Wait()
	fmt.Printf("Expected: %d | Actual Counter Value: %d (100%% Thread-Safe!)\n", numGoroutines, counter.Value())
}

// -----------------------------------------------------------------------------
// 3. sync.RWMutex (Reader-Writer Lock: Many Readers OR One Writer)
// -----------------------------------------------------------------------------
// When to use RWMutex:
// - Ideal for READ-HEAVY workloads (e.g., caches, config lookups).
// - Multiple goroutines can hold a READ lock (RLock) simultaneously.
// - Only ONE goroutine can hold a WRITE lock (Lock) exclusively.
type SafeCache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewSafeCache() *SafeCache {
	return &SafeCache{
		data: make(map[string]string),
	}
}

// Set acquires an EXCLUSIVE WRITE lock (blocks all readers and other writers).
func (c *SafeCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

// Get acquires a SHARED READ lock (multiple readers can read at the same time!).
func (c *SafeCache) Get(key string) (string, bool) {
	c.mu.RLock()         // Read lock
	defer c.mu.RUnlock() // Read unlock

	val, exists := c.data[key]
	return val, exists
}

func demoRWMutex() {
	fmt.Println("\n--- 3. sync.RWMutex (Concurrent Readers & Single Writer) ---")

	cache := NewSafeCache()
	var wg sync.WaitGroup

	// Initial write
	cache.Set("status", "ACTIVE")

	// 5 Concurrent Readers (Can all read simultaneously without blocking each other!)
	for r := 1; r <= 5; r++ {
		wg.Add(1)
		go func(readerID int) {
			defer wg.Done()
			val, _ := cache.Get("status")
			fmt.Printf("   [Reader %d] Read cache: status=%s\n", readerID, val)
		}(r)
	}

	// 1 Concurrent Writer (Will safely wait for readers to release, update, then unlock)
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
		cache.Set("status", "MAINTENANCE")
		fmt.Println("   [Writer] Cache updated to status=MAINTENANCE")
	}()

	wg.Wait()
	finalVal, _ := cache.Get("status")
	fmt.Println("Final Cache Status:", finalVal)
}

// =============================================================================
// Main Entry Point
// =============================================================================
func main() {
	fmt.Println("==================================================")
	fmt.Println("    GO MUTEX, RWMUTEX & RACE DETECTOR GUIDE       ")
	fmt.Println("==================================================")

	// 1. Concept of Race Conditions
	demoRaceConditionConcept()

	// 2. sync.Mutex demo
	demoMutex()

	// 3. sync.RWMutex demo
	demoRWMutex()

	fmt.Println("==================================================")
}

/*
================================================================================
Production Rules & Best Practices for Mutexes
================================================================================

1. Mutex vs RWMutex Comparison:
   | Lock Type     | Readers Allowed | Writers Allowed | Best Use Case                   |
   |---------------|-----------------|-----------------|---------------------------------|
   | sync.Mutex    | 1               | 1 (Exclusive)   | Balanced or write-heavy loads   |
   | sync.RWMutex  | Many (Shared)   | 1 (Exclusive)   | Read-heavy loads (e.g. Caches)  |

2. Always Defer Unlock:
   - Call `c.mu.Lock()` followed immediately by `defer c.mu.Unlock()`.
   - Prevents deadlocks if a panic or early return occurs inside the method.

3. NEVER Copy a Mutex:
   - Always pass structs containing a Mutex BY POINTER (`*SafeCounter`).
   - If copied, the copy gets an independent lock state, destroying synchronization!

4. Keep Critical Sections Short:
   - Do NOT perform long-running operations (I/O, database queries, HTTP calls)
     while holding a lock. Do calculations/I/O outside the lock, and only hold
     the lock during the actual shared-memory read/write.

5. How to Run the Race Detector:
   In your terminal:
       go run -race DAY_2/locking_mutex.go
   In your automated CI/CD tests:
       go test -race ./...
   (The race detector adds ~2-10x CPU overhead, so it is used during testing and
   staging, not production binaries).

6. Mutex vs Channels (When to use which?):
   - Use Mutex: Protecting shared internal state / caching / counters / data structures.
   - Use Channels: Passing ownership of data, coordinating work, pipelines, signaling.
================================================================================

Common concurrency problems
Problem	Meaning
Deadlock	Tasks wait forever
Race condition	Unsynchronized access to shared data
Starvation	A task cannot get sufficient execution resources
Livelock	Tasks keep reacting but make no useful progress
Goroutine leak	A goroutine remains blocked or active unnecessarily
Lock contention	Multiple goroutines compete for the same lock
*/

