package main

import "fmt"

// -------------------------------------------------------------
// 1. Primitive Pointers (Pass by Value vs Pass by Pointer)
// -------------------------------------------------------------
func modifyByValue(x int) {
	x = 100
}

func modifyByPointer(x *int) {
	*x = 100 // Dereferencing: write to the address x points to
}

// -------------------------------------------------------------
// 2. Struct Pointers (Extremely common in backend development)
// -------------------------------------------------------------
type User struct {
	Name string
	Age  int
}

func updateUser(user *User) {
	// Go automatically dereferences struct pointers when accessing fields,
	// so user.Age works directly without needing (*user).Age = 25.
	user.Age = 25
}

func main() {
	// ==========================================
	// Section 1: Pointer Basics (& vs *)
	// ==========================================
	fmt.Println("=== 1. Pointer Basics ===")
	val := 10
	var ptr *int = &val // & gets the memory address of val

	fmt.Println("val value:         ", val)
	fmt.Println("val address (ptr): ", ptr)
	fmt.Println("dereference (*ptr):", *ptr)

	// ==========================================
	// Section 2: Pass by Value vs Pointer
	// ==========================================
	fmt.Println("\n=== 2. Value vs Pointer ===")
	num1 := 10
	modifyByValue(num1)
	fmt.Println("num1 after modifyByValue:  ", num1) // 10 (unchanged)

	num2 := 10
	modifyByPointer(&num2)
	fmt.Println("num2 after modifyByPointer:", num2) // 100 (modified!)

	// ==========================================
	// Section 3: Struct Pointers
	// ==========================================
	fmt.Println("\n=== 3. Struct Pointers ===")
	user := User{
		Name: "Alex",
		Age:  20,
	}

	fmt.Println("Before updateUser:", user.Age) // 20
	updateUser(&user)                           // Pass address of struct
	fmt.Println("After updateUser: ", user.Age) // 25 (modified!)

	// ==========================================
	// Section 4: The new() Built-in Function
	// ==========================================
	fmt.Println("\n=== 4. The new() Function ===")
	// new(int) allocates a zero-valued integer and returns its pointer (*int)
	newPtr := new(int)
	fmt.Println("Default value of new(int):", *newPtr) // 0

	*newPtr = 100
	fmt.Println("Value after assignment:   ", *newPtr) // 100
}
