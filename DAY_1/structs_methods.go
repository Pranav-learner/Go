package main

import "fmt"

// 1. Defining the struct (Go's equivalent of a Class data model)
type Employee struct {
	ID         int
	Name       string
	Department string
	Salary     float64
}

// 2. Defining a Method on Employee:
//
//	func (e Employee) GiveBonus() float64
//	      ^^^^^^^^^^
//	      This is the RECEIVER! It attaches this function to the Employee struct.
//	      'e' acts just like 'self' in Python or 'this' in Java/C++.
func (e Employee) GiveBonus() float64 {
	return e.Salary * 0.2
}

// 3. Pointer Receiver (if you want the method to MODIFY the struct)
func (e *Employee) RaiseSalary(amount float64) {
	e.Salary += amount
}

func main() {
	// Creating an instance ("object") of Employee inside main
	employee := Employee{
		ID:         1,
		Name:       "John",
		Department: "Engineering",
		Salary:     100000,
	}

	fmt.Println("Employee Name:  ", employee.Name)
	fmt.Println("Calculated Bonus:", employee.GiveBonus())

	// Using the pointer receiver method to modify salary
	employee.RaiseSalary(20000)
	fmt.Println("New Salary:     ", employee.Salary)
}
