package reporter

import (
	"fmt"
	"strings"

	"github.com/suer/tf-arm/internal/analyzer"
)

type Reporter struct{}

func New() *Reporter {
	return &Reporter{}
}

func (r *Reporter) PrintAnalysis(analysis analyzer.ARM64Analysis) {
	fmt.Printf("Resource: %s\n", analysis.FullAddress)
	fmt.Printf("  Current Architecture: %s\n", analysis.CurrentArch)
	fmt.Printf("  ARM64 Compatible: %v\n", analysis.ARM64Compatible)
	if analysis.ARM64Compatible && analysis.RecommendedArch != "" {
		fmt.Printf("  Recommended: %s\n", analysis.RecommendedArch)
	}
	fmt.Printf("  Notes: %s\n", analysis.Notes)
	fmt.Println()
}

func (r *Reporter) PrintSummary(totalAnalyzed, arm64Compatible, nonArm64Compatible int) {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Analysis Summary:\n")
	fmt.Printf("  Total analyzed resources: %d\n", totalAnalyzed)
	fmt.Printf("  ARM64 compatible resources: %d\n", arm64Compatible)
	fmt.Printf("  Resources already using ARM64: %d\n", arm64Compatible-nonArm64Compatible)
	fmt.Printf("  Resources that can migrate to ARM64: %d\n", nonArm64Compatible)
	if arm64Compatible > 0 {
		fmt.Printf("  Percentage of ARM64-capable resources not using ARM64: %.1f%%\n", float64(nonArm64Compatible)/float64(arm64Compatible)*100)
	}
}

func (r *Reporter) PrintHeader(resourceCount int) {
	fmt.Printf("Found %d resources\n", resourceCount)
	fmt.Println(strings.Repeat("=", 80))
}
