package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("tf-arm: Terraform State ARM64 Analyzer")

	if len(os.Args) < 2 {
		fmt.Println("Usage: tf-arm <terraform-state-file>")
		os.Exit(1)
	}

	stateFile := os.Args[1]
	fmt.Printf("Analyzing Terraform state file: %s\n", stateFile)
}
