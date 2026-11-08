package main

import "fmt"

func main() {
	// ==========================================
	// 1. IF / ELSE
	// ==========================================
	age := 20
	if age >= 18 {
		fmt.Println("Adult")
	} else {
		fmt.Println("Minor")
	}

	// Go special: If with a short statement (scoped variable)
	if score := 85; score >= 80 {
		fmt.Println("Passed with distinction:", score)
	}

	// ==========================================
	// 2. FOR LOOP (Standard 3-component loop)
	// ==========================================
	fmt.Println("\nStandard for loop:")
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	// ==========================================
	// 3. "WHILE" LOOP (Go uses 'for' with only a condition)
	// ==========================================
	fmt.Println("\nWhile-style loop:")
	count := 0
	for count < 3 {
		fmt.Println(count)
		count++
	}

	// ==========================================
	// 4. FOR ... RANGE (Iterating over collections)
	// ==========================================
	fmt.Println("\nRange loop over slice:")
	numbers := []int{10, 20, 30}
	for index, value := range numbers {
		fmt.Println("Index:", index, "Value:", value)
	}

	// If you only need the value (use blank identifier _ for index)
	for _, val := range numbers {
		fmt.Println("Value only:", val)
	}

	// Infinite loop:
	for {
		// break when done
	}

	// ==========================================
	// 5. SWITCH STATEMENT (No 'break' needed in Go)
	// ==========================================
	fmt.Println("\nSwitch statement:")
	day := "Monday"
	switch day {
	case "Monday":
		fmt.Println("Start of the work week")
	case "Friday":
		fmt.Println("Weekend is near!")
	default:
		fmt.Println("Midweek day")
	}
}
