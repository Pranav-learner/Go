package main

import (
	"context"
	"fmt"
	"time"
)

// =============================================================================
// DAY 2: THE CONTEXT PACKAGE (CANCELLATION, TIMEOUTS & VALUES)
// =============================================================================
//
// What is Context (`context.Context`)?
// In Go, `context` is used to carry:
//   1. Cancellation signals
//   2. Deadlines / Timeouts
//   3. Request-scoped values
// across API boundaries and between goroutines.
//
// Why is it critical?
// Without Context, background goroutines could keep running forever even after
// a client disconnects or an operation is aborted — this causes GOROUTINE LEAKS!
// =============================================================================

// -----------------------------------------------------------------------------
// 1. Root Contexts: Background() vs TODO()
// -----------------------------------------------------------------------------
//   - `context.Background()`: The top-level root context. Used in main(), tests,
//     and top-level request entrypoints.
//   - `context.TODO()`:       A placeholder when you aren't sure which context
//     to use yet, or code is in transition.
func demoRootContexts() {
	fmt.Println("--- 1. Root Contexts ---")
	ctxBg := context.Background()
	ctxTodo := context.TODO()

	fmt.Println("Background Context:", ctxBg)
	fmt.Println("TODO Context:      ", ctxTodo)
}

// -----------------------------------------------------------------------------
// 2. context.WithCancel (Manual Cancellation)
// -----------------------------------------------------------------------------
// Used when an operation should be cancelled on-demand (e.g. user aborts).
func longRunningTask(ctx context.Context, name string) {
	for {
		select {
		// ctx.Done() returns a channel that is CLOSED when the context is cancelled
		case <-ctx.Done():
			fmt.Printf("[%s] Aborting! Reason: %v\n", name, ctx.Err())
			return
		default:
			fmt.Printf("[%s] Working...\n", name)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func demoWithCancel() {
	fmt.Println("\n--- 2. context.WithCancel (Manual Cancellation) ---")

	// Create child context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// Rule: ALWAYS defer cancel() to prevent resource/goroutine leaks!
	defer cancel()

	// Launch worker
	go longRunningTask(ctx, "Worker-1")

	// Let the worker run for 250ms, then manually cancel it
	time.Sleep(250 * time.Millisecond)
	fmt.Println("[Main] Calling cancel() now...")
	cancel() // Signals ctx.Done() channel to close!

	// Give worker time to print abort message
	time.Sleep(50 * time.Millisecond)
}

// -----------------------------------------------------------------------------
// 3. context.WithTimeout & context.WithDeadline
// -----------------------------------------------------------------------------
// Used when an operation MUST complete within a timeframe (e.g. HTTP call, DB query).
// If time expires before completion, Go automatically cancels the context.
func fetchRemoteData(ctx context.Context) (string, error) {
	// Simulate a slow database query or external API call taking 300ms
	select {
	case <-time.After(300 * time.Millisecond):
		return "Database query result", nil
	case <-ctx.Done():
		// Context timed out before work completed
		return "", ctx.Err() // returns context.DeadlineExceeded
	}
}

func demoWithTimeout() {
	fmt.Println("\n--- 3. context.WithTimeout (Automatic Cancellation) ---")

	// Set a timeout of only 150ms (the operation takes 300ms, so it WILL time out!)
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	fmt.Println("[Main] Requesting data with 150ms timeout...")
	result, err := fetchRemoteData(ctx)
	if err != nil {
		fmt.Println("[Main] Request Failed:", err) // Prints: context deadline exceeded
	} else {
		fmt.Println("[Main] Success:", result)
	}
}

// -----------------------------------------------------------------------------
// 4. context.WithValue (Request-Scoped Metadata)
// -----------------------------------------------------------------------------
// Used for passing request-scoped metadata (e.g. Request ID, Trace ID, User Auth).
//
// Best Practice: Use a custom unexported type for the key to avoid key
// collisions between different packages!
type contextKey string

const requestIDKey contextKey = "requestID"

func processUserRequest(ctx context.Context) {
	// Retrieve value from context
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		fmt.Println("[Service] Processing request with ID:", reqID)
	} else {
		fmt.Println("[Service] No request ID found in context.")
	}
}

func demoWithValue() {
	fmt.Println("\n--- 4. context.WithValue (Request-Scoped Values) ---")

	// Attach a request ID to the context
	ctx := context.WithValue(context.Background(), requestIDKey, "req-xyz-98765")

	processUserRequest(ctx)
}

// =============================================================================
// Main Entry Point
// =============================================================================
func main() {
	fmt.Println("==================================================")
	fmt.Println("            GO CONTEXT PACKAGE GUIDE              ")
	fmt.Println("==================================================")

	// 1. Root contexts
	demoRootContexts()

	// 2. Manual cancellation
	demoWithCancel()

	// 3. Timeouts & Deadlines
	demoWithTimeout()

	// 4. Request-scoped values
	demoWithValue()

	fmt.Println("==================================================")
}

/*
================================================================================
Golden Rules of Context in Production Go
================================================================================

1. Rule #1: First parameter rule
   - Always pass `ctx context.Context` as the FIRST parameter of functions:
     `func QueryDB(ctx context.Context, query string) (*Row, error)`
   - Do NOT pass a nil context! If unsure, use `context.TODO()`.

2. Rule #2: Never store Context inside a struct
   - Pass context explicitly along the call chain. Context is meant to be
     ephemeral and scoped to an operation/request.

3. Rule #3: Always defer cancel()
   - Whenever you call `WithCancel`, `WithTimeout`, or `WithDeadline`, ALWAYS
     immediately write `defer cancel()`. This releases timers and resources.

4. Rule #4: Context forms a tree
   - Cancelling a parent context automatically cancels ALL its child contexts!
   - Children cannot cancel their parent.

5. Rule #5: WithValue is ONLY for request-scoped metadata
   - Use `WithValue` for: Request IDs, trace headers, authenticated user tokens.
   - Do NOT use `WithValue` for optional function parameters or business logic!
     Use standard function parameters or configuration structs instead.

6. How to Listen to Cancellation:
   - Check `<-ctx.Done()` inside a `select` statement.
   - Check `ctx.Err()` to inspect whether it was `context.Canceled` or
     `context.DeadlineExceeded`.
================================================================================
Production rules for context
1. Pass context as the first parameter: func GetUser(ctx context.Context, id int).
2. Don't store request contexts inside long-lived structs.
3. Call cancellation functions, usually with defer cancel().
4. Don't pass nil as a context.
5. Use request contexts for downstream database and HTTP calls.
6. Don't use context values as a replacement for ordinary function parameters.
7. Use context cancellation to prevent unnecessary background work.

Goroutine leaks
Consider:
func worker(ch chan int) {

    value := <-ch

    fmt.Println(value)
}


If nobody ever sends to ch, the goroutine remains blocked.

If your application repeatedly creates such goroutines, resources accumulate.

This is a goroutine leak.

Prevent it through:
- Context cancellation.
- Correct channel lifecycle management.
- Timeouts where appropriate.
- Ensuring every goroutine has a termination path.
6. Debugging deadlocks
Useful commands:
go test -timeout 30s ./...


go test -race ./...


You can also inspect goroutine stacks using runtime diagnostics and pprof.

Production takeaway: Whenever you create a goroutine, ask yourself:
1. What makes it stop?
2. What happens if its input never arrives?
3. What happens if its output cannot be consumed?
4. Who owns its channels?
5. Who handles its errors?
*/
