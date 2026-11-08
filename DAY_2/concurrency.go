package main

import (
	"fmt"
	"runtime"
	"time"
)

// =============================================================================
// DAY 2: CONCURRENCY & GOROUTINES IN GO
// =============================================================================
//
// 1. What is Concurrency?
//    - Concurrency is the composition of independently executing computations.
//    - Rob Pike's quote: "Concurrency is about dealing with lots of things at once.
//      Parallelism is about doing lots of things at once."
//
// 2. What is a Goroutine?
//    - A lightweight thread of execution managed entirely by the Go runtime (not OS).
//    - OS Thread: Heavy (~1-2 MB stack, high context switch cost).
//    - Goroutine: Ultra-light (~2 KB initial stack, grows/shrinks dynamically).
//    - You can easily run hundreds of thousands of goroutines simultaneously.
//
// 3. How to launch a Goroutine?
//    - Just add the `go` keyword before any function call!
// =============================================================================

// A sample worker function
func printNumbers(prefix string, count int) {
	for i := 1; i <= count; i++ {
		fmt.Printf("[%s] Number: %d\n", prefix, i)
		time.Sleep(100 * time.Millisecond) // Simulate some work
	}
}

// -----------------------------------------------------------------------------
// Concept 1: The Main Goroutine Exits Early Problem
// -----------------------------------------------------------------------------
// When main() finishes, the entire program terminates immediately, killing
// any running goroutines without waiting for them!
func demoGoroutineWithSleep() {
	fmt.Println("--- 1. Launching Goroutines with time.Sleep ---")

	// Launch background task using `go`
	go printNumbers("Goroutine-A", 3)
	go printNumbers("Goroutine-B", 3)

	// Direct call on main goroutine
	fmt.Println("[Main] Goroutines started. Waiting with time.Sleep...")

	// Why time.Sleep here?
	// If main does not wait, the program will terminate before Goroutine-A and B finish!
	time.Sleep(400 * time.Millisecond)

	fmt.Println("[Main] Finished waiting.")
}

// -----------------------------------------------------------------------------
// Concept 2: Anonymous Function as a Goroutine
// -----------------------------------------------------------------------------
func demoAnonymousGoroutine() {
	fmt.Println("\n--- 2. Anonymous Goroutine with Closure ---")

	msg := "Hello from background goroutine!"

	// Launch inline anonymous function
	go func(text string) {
		fmt.Println("[Worker]", text)
	}(msg) // Pass variable as argument to avoid closure race conditions

	time.Sleep(100 * time.Millisecond)
}

// -----------------------------------------------------------------------------
// Concept 3: Concurrency vs Parallelism in Go
// -----------------------------------------------------------------------------
//   - Concurrency:  Structuring a program by breaking it into independent tasks.
//     Can run on a SINGLE CPU core (swapping via time-slicing).
//   - Parallelism:  Actually executing multiple tasks at the exact same physical instant.
//     REQUIRES multiple CPU cores!
//
// Go's M:N Scheduler:
//
//	Maps M goroutines onto N OS threads across P logical processors (runtime.GOMAXPROCS).
func demoParallelismInfo() {
	fmt.Println("\n--- 3. Parallelism & Go Runtime Info ---")
	fmt.Println("Total CPU Cores Available:    ", runtime.NumCPU())
	fmt.Println("GOMAXPROCS (Active CPU Cores):", runtime.GOMAXPROCS(0))
	fmt.Println("Active Goroutines count:      ", runtime.NumGoroutine())
}

// =============================================================================
// Main Entry Point
// =============================================================================
func main() {
	fmt.Println("==================================================")
	fmt.Println("       GO CONCURRENCY & GOROUTINES GUIDE          ")
	fmt.Println("==================================================")

	// 1. Basic Goroutine execution with time.Sleep
	demoGoroutineWithSleep()

	// 2. Anonymous Goroutine execution
	demoAnonymousGoroutine()

	// 3. Inspecting Go runtime parallelism
	demoParallelismInfo()

	fmt.Println("==================================================")
}

/*
================================================================================
Production Note: Why time.Sleep() is an Anti-Pattern
================================================================================
In real production code, NEVER use time.Sleep() to wait for goroutines!
- If the task takes longer than the sleep duration: The program exits before
  the work completes (data loss, corrupted state).
- If the task takes less time: You waste idle CPU time waiting unnecessarily.

Production Solutions:
1. Channels: Used for communication and synchronization between goroutines.
2. sync.WaitGroup: Used to wait for a collection of goroutines to finish.
================================================================================
*/
