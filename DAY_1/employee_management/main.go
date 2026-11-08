package main

import (
	"errors"
	"fmt"
	"employee-management/repository"
	"employee-management/service"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("        EMPLOYEE MANAGEMENT SYSTEM (CLI)          ")
	fmt.Println("==================================================")
	// 1. Dependency Injection: Wire up Repository and Service
	repo := repository.NewInMemoryRepository()
	svc := service.NewEmployeeService(repo)
	// ---------------------------------------------------------------
	// 2. CREATE (Add Employees)
	// ---------------------------------------------------------------
	fmt.Println("\n--- 1. Adding Employees ---")
	err := svc.CreateEmployee(1, "Alice Smith", "Backend", 95000)
	printResult("Add Alice (ID: 1)", err)
	err = svc.CreateEmployee(2, "Bob Jones", "Frontend", 85000)
	printResult("Add Bob (ID: 2)", err)
	err = svc.CreateEmployee(3, "Charlie Day", "DevOps", 90000)
	printResult("Add Charlie (ID: 3)", err)
	// Test Duplicate ID Error
	err = svc.CreateEmployee(1, "Duplicate Alice", "QA", 70000)
	printResult("Add Duplicate ID 1", err)
	if errors.Is(err, repository.ErrAlreadyExists) {
		fmt.Println("   --> Handled expected error: ID already exists")
	}
	// Test Validation Error
	err = svc.CreateEmployee(4, "", "QA", -500)
	printResult("Add with Empty Name & Negative Salary", err)
	// ---------------------------------------------------------------
	// 3. READ (Get by ID & Not Found Handling)
	// ---------------------------------------------------------------
	fmt.Println("\n--- 2. Fetching Employees ---")
	emp, err := svc.GetEmployee(1)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Found Employee #%d: %s | %s | $%.2f\n", emp.ID, emp.Name, emp.Department, emp.Salary)
	}
	// Lookup non-existent employee
	_, err = svc.GetEmployee(999)
	printResult("Lookup Non-Existent ID 999", err)
	if errors.Is(err, repository.ErrNotFound) {
		fmt.Println("   --> Handled expected error: Employee not found")
	}
	// ---------------------------------------------------------------
	// 4. UPDATE (Salary)
	// ---------------------------------------------------------------
	fmt.Println("\n--- 3. Updating Salary ---")
	err = svc.UpdateSalary(1, 105000)
	printResult("Promote Alice (ID: 1) to $105,000", err)
	updatedEmp, _ := svc.GetEmployee(1)
	fmt.Printf("Verified New Salary: $%.2f\n", updatedEmp.Salary)
	// ---------------------------------------------------------------
	// 5. DELETE
	// ---------------------------------------------------------------
	fmt.Println("\n--- 4. Deleting Employee ---")
	err = svc.DeleteEmployee(2)
	printResult("Delete Bob (ID: 2)", err)
	// Verify Bob is deleted
	_, err = svc.GetEmployee(2)
	printResult("Verify Bob is gone", err)
	// ---------------------------------------------------------------
	// 6. LIST ALL REMAINING
	// ---------------------------------------------------------------
	fmt.Println("\n--- 5. All Active Employees ---")
	all := svc.ListEmployees()
	fmt.Printf("%-5s %-15s %-15s %-10s\n", "ID", "Name", "Department", "Salary")
	fmt.Println("--------------------------------------------------")
	for _, e := range all {
		fmt.Printf("%-5d %-15s %-15s $%-9.2f\n", e.ID, e.Name, e.Department, e.Salary)
	}
	fmt.Println("==================================================")
}
func printResult(action string, err error) {
	if err != nil {
		fmt.Printf("[FAIL] %s: %v\n", action, err)
	} else {
		fmt.Printf("[OK]   %s\n", action)
	}
}