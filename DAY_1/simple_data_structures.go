package main

import "fmt"

func main() {
	// ==========================================
	// 1. ARRAYS (Fixed length)
	// ==========================================
	var numbers [5]int
	numbers[0] = 10
	numbers[1] = 20
	// numbers := [3]int{10, 20, 30}
	fmt.Println("Array:", numbers)

	// ==========================================
	// 2. SLICES (Dynamic arrays)
	// ==========================================
	// Notice: Must specify the type after [], e.g., []int
	sliceNumbers := []int{1, 2, 3, 4, 5}
	sliceNumbers = append(sliceNumbers, 6, 7)
	fmt.Println("Slice:", sliceNumbers)

	// ==========================================
	// 3. MAPS (Key-value pairs / HashMaps)
	// ==========================================
	user := map[string]string{
		"name": "Alex",
		"role": "Backend Engineer",
	}

	// Access a value
	fmt.Println("User Name:", user["name"])

	// Add or update a key
	user["city"] = "Berlin"

	// Check if a key exists ("comma ok" idiom)
	role, exists := user["role"]
	fmt.Println("Role:", role, "| Exists:", exists)

	missing, exists := user["salary"]
	fmt.Println("Salary:", missing, "| Exists:", exists)

	// Delete a key
	delete(user, "city")
	fmt.Println("Updated Map:", user)
}
