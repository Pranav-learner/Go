package main

import (
	"errors"
	"fmt"
)

// =============================================================================
// 1. Basic Error Handling (Returning error as the last value)
// =============================================================================
func withdraw(balance, amount float64) (float64, error) {
	if amount <= 0 {
		return balance, errors.New("invalid amount")
	}

	if amount > balance {
		return balance, errors.New("insufficient balance")
	}

	return balance - amount, nil
}

// =============================================================================
// 2. Error Wrapping (fmt.Errorf with %w)
// =============================================================================
// Simulated database function
func findUserInDB(id int) error {
	return errors.New("connection timeout")
}

// getUser preserves the original error while adding contextual info using %w
func getUser(id int) error {
	err := findUserInDB(id)
	if err != nil {
		return fmt.Errorf("failed to get user %d: %w", id, err)
	}
	return nil
}

// =============================================================================
// 3. Sentinel Errors and errors.Is()
// =============================================================================
// Sentinel error: A known, predefined error variable
var ErrNotFound = errors.New("user not found")

func findUser(id int) error {
	if id == 10 {
		// Can even be wrapped, and errors.Is will still detect it!
		return fmt.Errorf("lookup failed: %w", ErrNotFound)
	}
	return nil
}

// =============================================================================
// 4. Custom Error Types and errors.As()
// =============================================================================
// By implementing Error() string, ValidationError satisfies the `error` interface
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

func validateInput(input string) error {
	if input == "" {
		return &ValidationError{
			Field:   "email",
			Message: "cannot be empty",
		}
	}
	return nil
}

// =============================================================================
// Main Execution
// =============================================================================
func main() {
	// --------------------------------------------------
	// 1. Basic error check
	// --------------------------------------------------
	fmt.Println("=== 1. Basic Error Check ===")
	balance, err := withdraw(1000, 500)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Remaining balance:", balance)
	}

	// --------------------------------------------------
	// 2. Error wrapping with %w
	// --------------------------------------------------
	fmt.Println("\n=== 2. Error Wrapping ===")
	wrapErr := getUser(42)
	if wrapErr != nil {
		fmt.Println("Wrapped Error:", wrapErr)
		// Unwrap to inspect the root cause
		fmt.Println("Root Cause:   ", errors.Unwrap(wrapErr))
	}

	// --------------------------------------------------
	// 3. errors.Is() - Checks if error chain contains a specific error
	// --------------------------------------------------
	fmt.Println("\n=== 3. errors.Is() ===")
	userErr := findUser(10)
	if errors.Is(userErr, ErrNotFound) {
		fmt.Println("Identified: User does not exist (matched ErrNotFound)!")
	}

	// --------------------------------------------------
	// 4. errors.As() - Checks if error matches a specific type & casts it
	// --------------------------------------------------
	fmt.Println("\n=== 4. errors.As() ===")
	valErr := validateInput("")
	var targetErr *ValidationError

	if errors.As(valErr, &targetErr) {
		fmt.Println("Validation error on field:", targetErr.Field)
		fmt.Println("Validation message:        ", targetErr.Message)
	}
}

/*
================================================================================
Production Error-Handling Principles & Java Comparison
================================================================================
1. Never silently ignore errors (don't do `_ = doSomething()`).
2. Add useful context when returning errors (use fmt.Errorf("...: %w", err)).
3. Don't expose internal database/system errors directly to API clients.
4. Use sentinel errors (e.g. ErrNotFound) for commonly identifiable conditions.
5. Keep errors separate from normal business results using multiple return values.
6. Avoid excessive error wrapping at every single line/layer.

Java Comparison:
- Java: Uses exceptions (`try-catch-finally`, `throw new CustomException()`)
  for flow control and error handling.
- Go: Has NO exceptions for regular flow control. Errors are ordinary return
  values (`val, err := ...`).
- Go's `errors.Is(err, Target)` is like checking `err.equals(Target)`.
- Go's `errors.As(err, &target)` is like Java's `catch (ValidationError e)`.
================================================================================
*/
