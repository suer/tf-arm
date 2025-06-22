package reporter

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/suer/tf-arm/internal/analyzer"
)

func TestNew(t *testing.T) {
	reporter := New()
	if reporter == nil {
		t.Error("New() returned nil")
	}
}

func TestReporter_PrintAnalysis(t *testing.T) {
	tests := []struct {
		name     string
		analysis analyzer.ARM64Analysis
		expected []string
	}{
		{
			name: "ARM64 compatible resource with recommendation",
			analysis: analyzer.ARM64Analysis{
				ResourceType:      "aws_instance",
				ResourceName:      "example",
				FullAddress:       "aws_instance.example",
				CurrentArch:       "x86_64",
				ARM64Compatible:   true,
				AlreadyUsingARM64: false,
				RecommendedArch:   "arm64",
				Notes:             "Can migrate to ARM64 for cost savings",
			},
			expected: []string{
				"Resource: aws_instance.example",
				"Current Architecture: x86_64",
				"ARM64 Compatible: true",
				"Recommended: arm64",
				"Notes: Can migrate to ARM64 for cost savings",
			},
		},
		{
			name: "ARM64 incompatible resource",
			analysis: analyzer.ARM64Analysis{
				ResourceType:      "aws_instance",
				ResourceName:      "legacy",
				FullAddress:       "aws_instance.legacy",
				CurrentArch:       "x86_64",
				ARM64Compatible:   false,
				AlreadyUsingARM64: false,
				RecommendedArch:   "",
				Notes:             "Instance type not available in ARM64",
			},
			expected: []string{
				"Resource: aws_instance.legacy",
				"Current Architecture: x86_64",
				"ARM64 Compatible: false",
				"Notes: Instance type not available in ARM64",
			},
		},
		{
			name: "Already using ARM64",
			analysis: analyzer.ARM64Analysis{
				ResourceType:      "aws_instance",
				ResourceName:      "arm_instance",
				FullAddress:       "aws_instance.arm_instance",
				CurrentArch:       "arm64",
				ARM64Compatible:   true,
				AlreadyUsingARM64: true,
				RecommendedArch:   "",
				Notes:             "Already optimized for ARM64",
			},
			expected: []string{
				"Resource: aws_instance.arm_instance",
				"Current Architecture: arm64",
				"ARM64 Compatible: true",
				"Notes: Already optimized for ARM64",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				reporter := New()
				reporter.PrintAnalysis(tt.analysis)
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("PrintAnalysis() output missing expected string %q\nGot: %s", expected, output)
				}
			}

			if tt.analysis.ARM64Compatible && tt.analysis.RecommendedArch != "" {
				if !strings.Contains(output, "Recommended:") {
					t.Error("PrintAnalysis() should include recommendation when ARM64Compatible is true and RecommendedArch is set")
				}
			}
		})
	}
}

func TestReporter_PrintSummary(t *testing.T) {
	tests := []struct {
		name               string
		totalAnalyzed      int
		arm64Compatible    int
		nonArm64Compatible int
		expectedStrings    []string
	}{
		{
			name:               "basic summary",
			totalAnalyzed:      10,
			arm64Compatible:    6,
			nonArm64Compatible: 4,
			expectedStrings: []string{
				"Analysis Summary:",
				"Total analyzed resources: 10",
				"ARM64 compatible resources: 6",
				"Resources already using ARM64: 2",
				"Resources that can migrate to ARM64: 4",
				"Percentage of ARM64-capable resources not using ARM64: 66.7%",
			},
		},
		{
			name:               "no ARM64 compatible resources",
			totalAnalyzed:      5,
			arm64Compatible:    0,
			nonArm64Compatible: 0,
			expectedStrings: []string{
				"Analysis Summary:",
				"Total analyzed resources: 5",
				"ARM64 compatible resources: 0",
				"Resources already using ARM64: 0",
				"Resources that can migrate to ARM64: 0",
			},
		},
		{
			name:               "all resources already using ARM64",
			totalAnalyzed:      3,
			arm64Compatible:    3,
			nonArm64Compatible: 0,
			expectedStrings: []string{
				"Analysis Summary:",
				"Total analyzed resources: 3",
				"ARM64 compatible resources: 3",
				"Resources already using ARM64: 3",
				"Resources that can migrate to ARM64: 0",
				"Percentage of ARM64-capable resources not using ARM64: 0.0%",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				reporter := New()
				reporter.PrintSummary(tt.totalAnalyzed, tt.arm64Compatible, tt.nonArm64Compatible)
			})

			for _, expected := range tt.expectedStrings {
				if !strings.Contains(output, expected) {
					t.Errorf("PrintSummary() output missing expected string %q\nGot: %s", expected, output)
				}
			}

			if strings.Count(output, "=") < 80 {
				t.Error("PrintSummary() should include separator line with 80 equals signs")
			}
		})
	}
}

func TestReporter_PrintHeader(t *testing.T) {
	tests := []struct {
		name          string
		resourceCount int
		expected      []string
	}{
		{
			name:          "single resource",
			resourceCount: 1,
			expected: []string{
				"Found 1 resources",
			},
		},
		{
			name:          "multiple resources",
			resourceCount: 42,
			expected: []string{
				"Found 42 resources",
			},
		},
		{
			name:          "zero resources",
			resourceCount: 0,
			expected: []string{
				"Found 0 resources",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				reporter := New()
				reporter.PrintHeader(tt.resourceCount)
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("PrintHeader() output missing expected string %q\nGot: %s", expected, output)
				}
			}

			if strings.Count(output, "=") < 80 {
				t.Error("PrintHeader() should include separator line with 80 equals signs")
			}
		})
	}
}

func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}