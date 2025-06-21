package main

import (
	"fmt"
	"os"
	"strings"
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

	state, err := parseStateFile(stateFile)
	if err != nil {
		fmt.Printf("Error parsing state file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d resources\n", len(state.Resources))
	fmt.Println(strings.Repeat("=", 80))

	var arm64CompatibleCount int
	var totalAnalyzedCount int

	for _, resource := range state.Resources {
		if resource.Mode != "managed" {
			continue
		}

		analysis := analyzeARM64Compatibility(resource)
		
		if analysis.Notes != "リソースタイプはARM64互換性チェック対象外" {
			totalAnalyzedCount++
			printAnalysis(analysis)
			
			if analysis.ARM64Compatible {
				arm64CompatibleCount++
			}
		}
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Analysis Summary:\n")
	fmt.Printf("  Total analyzed resources: %d\n", totalAnalyzedCount)
	fmt.Printf("  ARM64 compatible resources: %d\n", arm64CompatibleCount)
	if totalAnalyzedCount > 0 {
		fmt.Printf("  Compatibility rate: %.1f%%\n", float64(arm64CompatibleCount)/float64(totalAnalyzedCount)*100)
	}
}

func printAnalysis(analysis ARM64Analysis) {
	fmt.Printf("Resource: %s.%s\n", analysis.ResourceType, analysis.ResourceName)
	fmt.Printf("  Current Architecture: %s\n", analysis.CurrentArch)
	fmt.Printf("  ARM64 Compatible: %v\n", analysis.ARM64Compatible)
	if analysis.ARM64Compatible && analysis.RecommendedArch != "" {
		fmt.Printf("  Recommended: %s\n", analysis.RecommendedArch)
	}
	fmt.Printf("  Notes: %s\n", analysis.Notes)
	fmt.Println()
}
