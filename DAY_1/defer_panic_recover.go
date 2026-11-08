package main

import (
	"fmt"
	"os"
)

// =============================================================================
// 1. Basic Defer & LIFO Order
// =============================================================================
// defer schedules a function call to execute when the surrounding function returns.
// Multiple deferred calls execute in LIFO (Last In, First Out) order.
func demoDeferLIFO() {
	fmt.Println("--- Defer & LIFO Order ---")

	// Notice: Execution order will be Third -> Second -> First
	defer fmt.Println("Deferred: First")
	defer fmt.Println("Deferred: Second")
	defer fmt.Println("Deferred: Third")

	fmt.Println("Normal: Start of function")
	fmt.Println("Normal: Processing...")
	fmt.Println("Normal: End of function")
}

// =============================================================================
// 2. Practical Production Example: Resource Cleanup
// =============================================================================
// Common pattern: defer resource.Close(), defer mutex.Unlock()
func readFileExample(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	// Even if an error or early return occurs later, file.Close() is guaranteed to run
	defer file.Close()

	// Pretend to read file here
	fmt.Println("File opened successfully:", filename)
	return nil
}

// =============================================================================
// 3. Panic & 4. Recover
// =============================================================================
// divide panics on unexpected invalid states
func divide(a, b int) int {
	if b == 0 {
		panic("division by zero") // Interrupts normal execution and unwinds stack
	}
	return a / b
}

// recover() must be called inside a deferred function to regain control from panic
func safeFunction() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	fmt.Println("Attempting dangerous division (10 / 0)...")
	divide(10, 0)
	fmt.Println("This line will NEVER execute because of panic")
}

// =============================================================================
// Main Entry Point
// =============================================================================
func main() {
	// Section 1: Basic Defer and LIFO demonstration
	fmt.Println("=== 1. Defer & LIFO ===")
	demoDeferLIFO()

	// Section 2: Resource Cleanup (production pattern)
	fmt.Println("\n=== 2. Resource Cleanup with Defer ===")
	_ = readFileExample("non_existent_file.txt")

	// Section 3 & 4: Panic and Recover
	fmt.Println("\n=== 3 & 4. Panic & Recover ===")
	safeFunction()

	fmt.Println("\nApplication safely continues executing after recovered panic!")
}

/*
================================================================================
Production Engineering Perspective & Comparison
================================================================================
1. Separation of Concerns:
   - Use 'error' for expected failures (network error, invalid input, file missing).
   - Use 'panic' only for genuinely exceptional bugs or violated invariants.
   - Use 'recover' only at top-level boundaries (e.g. HTTP middleware, worker pools).

2. Common Defer Patterns:
   - defer file.Close()
   - defer rows.Close()          (database queries)
   - defer resp.Body.Close()     (HTTP requests)
   - defer mutex.Unlock()        (concurrency locks)

3. Goroutine Boundary Rule:
   - 'recover()' ONLY catches panics inside the SAME goroutine.
   - It cannot catch a panic originating from another goroutine!

4. Java Comparison:
   | Go Pattern                   | Java Equivalent                    |
   |------------------------------|------------------------------------|
   | defer cleanup()              | try-finally / try-with-resources   |
   | panic("boom")                | throw new RuntimeException("boom") |
   | recover() inside defer       | catch (Throwable t)                |
================================================================================
*/
