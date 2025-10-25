package main

import (
	"fmt"

	"github.com/yourusername/monorepo-lib/libs/greetings"
	"github.com/yourusername/monorepo-lib/libs/math"
)

func main() {
	// Using greetings library
	fmt.Println("=== Greetings Library Demo ===")
	fmt.Println(greetings.Hello("World"))
	fmt.Println(greetings.Hello("Gopher"))
	fmt.Println(greetings.Goodbye("Friend"))
	fmt.Println(greetings.Welcome("Alice", "Bob", "Charlie"))

	fmt.Println()

	// Using math library
	fmt.Println("=== Math Library Demo ===")
	a, b := 10, 5

	fmt.Printf("Add(%d, %d) = %d\n", a, b, math.Add(a, b))
	fmt.Printf("Subtract(%d, %d) = %d\n", a, b, math.Subtract(a, b))
	fmt.Printf("Multiply(%d, %d) = %d\n", a, b, math.Multiply(a, b))
	fmt.Printf("Divide(%d, %d) = %d\n", a, b, math.Divide(a, b))
	fmt.Printf("Max(%d, %d) = %d\n", a, b, math.Max(a, b))
	fmt.Printf("Min(%d, %d) = %d\n", a, b, math.Min(a, b))

	fmt.Println()
	fmt.Println("=== Monorepo Workspace Example ===")
	fmt.Println("This example demonstrates a Go workspace with multiple libraries!")
}
