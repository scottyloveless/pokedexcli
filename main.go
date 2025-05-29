package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func cleanInput(text string) []string {
	strings.ToLower(text)
	parts := strings.Fields(text)
	return parts
}
