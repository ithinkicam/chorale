// Package main calls fmt.Println, which the fixture database marks as vulnerable.
package main

import "fmt"

func main() {
	fmt.Println("fixture")
}
