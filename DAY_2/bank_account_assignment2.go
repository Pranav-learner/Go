package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// DAY 2: ASSIGNMENT 2 — SYNCHRONIZATION
// CONCURRENT BANK ACCOUNT SIMULATOR
// =============================================================================
// Requirements:
//   - Initial balance: 10,000
//   - Start 100 goroutines: Each deposits 100
//   - Start 50 goroutines: Each withdraws 50
//   - Protect balance using sync.Mutex
//   - Verify final balance: 10,000 + 10,000 - 2,500 = 17,500
// =============================================================================

var ErrInsufficientFunds = errors.New("insufficient funds for withdrawal")

// BankAccount encapsulates balance and a Mutex for thread safety
type BankAccount struct {
	mu      sync.Mutex // Guards balance against concurrent access
	balance float64
}

// NewBankAccount creates an account with an initial balance
func NewBankAccount(initialBalance float64) *BankAccount {
	return &BankAccount{
		balance: initialBalance,
	}
}

// Deposit safely adds funds to the account
func (a *BankAccount) Deposit(amount float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.balance += amount
}

// Withdraw safely removes funds if balance is sufficient
func (a *BankAccount) Withdraw(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.balance < amount {
		return ErrInsufficientFunds
	}

	a.balance -= amount
	return nil
}

// Balance safely reads the current account balance
func (a *BankAccount) Balance() float64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.balance
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("   ASSIGNMENT 2: CONCURRENT BANK ACCOUNT          ")
	fmt.Println("==================================================")

	const (
		initialBalance = 10000.0
		numDeposits    = 100
		depositAmount  = 100.0
		numWithdrawals = 50
		withdrawAmount = 50.0
	)

	// Expected calculation:
	// 10,000 + (100 * 100) - (50 * 50) = 10,000 + 10,000 - 2,500 = 17,500
	expectedBalance := initialBalance + (float64(numDeposits) * depositAmount) - (float64(numWithdrawals) * withdrawAmount)

	account := NewBankAccount(initialBalance)
	fmt.Printf("[Account Initialized] Starting Balance: $%.2f\n\n", account.Balance())

	var wg sync.WaitGroup

	startTime := time.Now()

	// -------------------------------------------------------------
	// 1. Launch 100 Deposit Goroutines concurrently
	// -------------------------------------------------------------
	fmt.Printf("--> Launching %d Deposit goroutines (+$%.2f each)...\n", numDeposits, depositAmount)
	for i := 1; i <= numDeposits; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			account.Deposit(depositAmount)
		}(i)
	}

	// -------------------------------------------------------------
	// 2. Launch 50 Withdrawal Goroutines concurrently
	// -------------------------------------------------------------
	fmt.Printf("--> Launching %d Withdrawal goroutines (-$%.2f each)...\n", numWithdrawals, withdrawAmount)
	for i := 1; i <= numWithdrawals; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := account.Withdraw(withdrawAmount); err != nil {
				fmt.Printf("[Warning] Withdrawal #%d failed: %v\n", id, err)
			}
		}(i)
	}

	// -------------------------------------------------------------
	// 3. Wait for all 150 concurrent transactions to complete
	// -------------------------------------------------------------
	wg.Wait()
	elapsed := time.Since(startTime)

	// -------------------------------------------------------------
	// 4. Verify Final Balance
	// -------------------------------------------------------------
	finalBalance := account.Balance()

	fmt.Println("\n==================================================")
	fmt.Println("                 AUDIT REPORT                     ")
	fmt.Println("==================================================")
	fmt.Printf("Initial Balance:    $%.2f\n", initialBalance)
	fmt.Printf("Total Deposits:     +%d x $%.2f = +$%.2f\n", numDeposits, depositAmount, float64(numDeposits)*depositAmount)
	fmt.Printf("Total Withdrawals:  -%d x $%.2f = -$%.2f\n", numWithdrawals, withdrawAmount, float64(numWithdrawals)*withdrawAmount)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Expected Balance:   $%.2f\n", expectedBalance)
	fmt.Printf("Actual Balance:     $%.2f\n", finalBalance)
	fmt.Printf("Execution Time:     %v\n", elapsed)
	fmt.Println("--------------------------------------------------")

	if finalBalance == expectedBalance {
		fmt.Println("STATUS: [PASSED] Balance matches expected 17,500! Zero race conditions.")
	} else {
		fmt.Printf("STATUS: [FAILED] Discrepancy detected! Difference: $%.2f\n", finalBalance-expectedBalance)
	}
	fmt.Println("==================================================")
}

