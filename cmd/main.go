package main

import (
	"fmt"
	"os"

	"github.com/suer/tf-arm/internal/analyzer"
	"github.com/suer/tf-arm/internal/parser"
	"github.com/suer/tf-arm/internal/reporter"
)

func main() {
	fmt.Println("tf-arm: Terraform State ARM64 Analyzer")

	if len(os.Args) < 2 {
		fmt.Println("Usage: tf-arm <terraform-state-file>")
		fmt.Println("")
		fmt.Println("This tool analyzes Terraform state files to identify resources")
		fmt.Println("that can be migrated to ARM64 architecture for cost optimization.")
		os.Exit(1)
	}

	stateFile := os.Args[1]
	fmt.Printf("Analyzing Terraform state file: %s\n", stateFile)
	fmt.Println("")

	state, err := parser.ParseStateFile(stateFile)
	if err != nil {
		fmt.Printf("Error parsing state file: %v\n", err)
		os.Exit(1)
	}

	rep := reporter.New()
	rep.PrintHeader(len(state.Resources))

	var arm64CompatibleCount int
	var totalAnalyzedCount int

	for _, resource := range state.Resources {
		if resource.Mode != "managed" {
			continue
		}

		analysis := analyzer.AnalyzeResource(resource)

		if analysis.Notes != "リソースタイプはARM64互換性チェック対象外" {
			totalAnalyzedCount++
			rep.PrintAnalysis(analysis)

			if analysis.ARM64Compatible {
				arm64CompatibleCount++
			}
		}
	}

	rep.PrintSummary(totalAnalyzedCount, arm64CompatibleCount)
}
