package main

import (
	"fmt"
)

func add(x, y float64) float64 {
	return x + y
}

func subtract(x, y float64) float64 {
	return x - y
}

func multiply(x, y float64) float64 {
	return x * y
}

func divide(x, y float64) (float64, error) {
	if y == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return x / y, nil
}

func main() {
	var num1, num2 float64
	var operator string

	fmt.Print("Enter first number: ")
	fmt.Scan(&num1)

	fmt.Print("Enter operator (+, -, *, /): ")
	fmt.Scan(&operator)

	fmt.Print("Enter second number: ")
	fmt.Scan(&num2)

	var result float64
	var err error

	switch operator {
	case "+":
		result = add(num1, num2)
	case "-":
		result = subtract(num1, num2)
	case "*":
		result = multiply(num1, num2)
	case "/":
		result, err = divide(num1, num2)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	default:
		fmt.Println("Invalid operator")
		return
	}

	fmt.Printf("Result: %v %s %v = %v\n", num1, operator, num2, result)
}
