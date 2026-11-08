package main

import "fmt"

func add(a int , b int) int{
	return a+b
}
//In Go, when consecutive function parameters 
//share the same type, you can omit the type for all but the last parameter.
func divide(a, b float64) (float64, error){
	if b == 0{
		return 0, fmt.Errorf("division by zero")
	}

	return a/b, nil
}
func main(){
	result := add(10,30)
	result1, err := divide(10,2)
	fmt.Println(result)
	fmt.Println(result1,err)
}