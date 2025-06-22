package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/suer/tf-arm/internal/analyzer"
	"github.com/suer/tf-arm/internal/parser"
)

func TestCanMigrateToARM64(t *testing.T) {
	tests := []struct {
		name     string
		analysis analyzer.ARM64Analysis
		expected bool
	}{
		{
			name: "can migrate - ARM64 compatible but not using ARM64",
			analysis: analyzer.ARM64Analysis{
				ARM64Compatible:   true,
				AlreadyUsingARM64: false,
			},
			expected: true,
		},
		{
			name: "cannot migrate - already using ARM64",
			analysis: analyzer.ARM64Analysis{
				ARM64Compatible:   true,
				AlreadyUsingARM64: true,
			},
			expected: false,
		},
		{
			name: "cannot migrate - not ARM64 compatible",
			analysis: analyzer.ARM64Analysis{
				ARM64Compatible:   false,
				AlreadyUsingARM64: false,
			},
			expected: false,
		},
		{
			name: "cannot migrate - not compatible and already using ARM64",
			analysis: analyzer.ARM64Analysis{
				ARM64Compatible:   false,
				AlreadyUsingARM64: true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := canMigrateToARM64(tt.analysis)
			if result != tt.expected {
				t.Errorf("canMigrateToARM64() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateMigrateablePercent(t *testing.T) {
	tests := []struct {
		name                 string
		migrateableCount     int
		arm64CompatibleCount int
		expected             float64
	}{
		{
			name:                 "normal calculation",
			migrateableCount:     4,
			arm64CompatibleCount: 6,
			expected:             66.0,
		},
		{
			name:                 "all migrateable",
			migrateableCount:     5,
			arm64CompatibleCount: 5,
			expected:             100.0,
		},
		{
			name:                 "none migrateable",
			migrateableCount:     0,
			arm64CompatibleCount: 3,
			expected:             0.0,
		},
		{
			name:                 "zero ARM64 compatible",
			migrateableCount:     0,
			arm64CompatibleCount: 0,
			expected:             0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateMigrateablePercent(tt.migrateableCount, tt.arm64CompatibleCount)
			if int(result) != int(tt.expected) {
				t.Errorf("calculateMigrateablePercent() = %v (int: %d), want %v (int: %d)", result, int(result), tt.expected, int(tt.expected))
			}
		})
	}
}

func TestAnalyzeStateFile_NonExistentFile(t *testing.T) {
	// This test assumes that analyzeStateFile will call os.Exit(1) for non-existent files
	// In a real test environment, we would need to mock os.Exit or test the function differently
	t.Skip("Skipping test that requires mocking os.Exit")
}

func TestAnalyzeStateFile_JSONOutput(t *testing.T) {
	// Create a temporary state file
	tempDir := t.TempDir()
	stateFile := filepath.Join(tempDir, "test.tfstate")
	
	state := parser.TerraformState{
		Version: 4,
		Resources: []parser.TerraformResource{
			{
				Mode: "managed",
				Type: "aws_instance",
				Name: "example",
				Instances: []parser.ResourceInstance{
					{
						Attributes: map[string]interface{}{
							"instance_type": "t3.micro",
						},
					},
				},
			},
			{
				Mode: "managed",
				Type: "aws_instance",
				Name: "arm_example",
				Instances: []parser.ResourceInstance{
					{
						Attributes: map[string]interface{}{
							"instance_type": "t4g.micro",
						},
					},
				},
			},
		},
	}
	
	stateData, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("Failed to marshal state: %v", err)
	}
	
	if err := os.WriteFile(stateFile, stateData, 0644); err != nil {
		t.Fatalf("Failed to write state file: %v", err)
	}

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	analyzeStateFile(stateFile, "json", 0)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Parse JSON output
	var jsonOutput JSONOutput
	if err := json.Unmarshal([]byte(output), &jsonOutput); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Verify summary
	if jsonOutput.Summary.TotalAnalyzed != 2 {
		t.Errorf("Expected TotalAnalyzed = 2, got %d", jsonOutput.Summary.TotalAnalyzed)
	}

	if jsonOutput.Summary.ARM64Compatible != 2 {
		t.Errorf("Expected ARM64Compatible = 2, got %d", jsonOutput.Summary.ARM64Compatible)
	}

	if jsonOutput.Summary.Migrateable != 1 {
		t.Errorf("Expected Migrateable = 1, got %d", jsonOutput.Summary.Migrateable)
	}

	// Verify resources
	if len(jsonOutput.Resources) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(jsonOutput.Resources))
	}
}

func TestAnalyzeStateFile_TextOutput(t *testing.T) {
	// Create a temporary state file
	tempDir := t.TempDir()
	stateFile := filepath.Join(tempDir, "test.tfstate")
	
	state := parser.TerraformState{
		Version: 4,
		Resources: []parser.TerraformResource{
			{
				Mode: "managed",
				Type: "aws_instance",
				Name: "example",
				Instances: []parser.ResourceInstance{
					{
						Attributes: map[string]interface{}{
							"instance_type": "t3.micro",
						},
					},
				},
			},
		},
	}
	
	stateData, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("Failed to marshal state: %v", err)
	}
	
	if err := os.WriteFile(stateFile, stateData, 0644); err != nil {
		t.Fatalf("Failed to write state file: %v", err)
	}

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	analyzeStateFile(stateFile, "text", 0)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expectedStrings := []string{
		"tf-arm: Terraform State ARM64 Analyzer",
		"Analyzing Terraform state file:",
		"Found 1 resources",
		"Resource: aws_instance.example",
		"Analysis Summary:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q, but it didn't. Output: %s", expected, output)
		}
	}
}

func TestAnalyzeStateFile_WithExitCode(t *testing.T) {
	// Create a temporary state file with migrateable resources
	tempDir := t.TempDir()
	stateFile := filepath.Join(tempDir, "test.tfstate")
	
	state := parser.TerraformState{
		Version: 4,
		Resources: []parser.TerraformResource{
			{
				Mode: "managed",
				Type: "aws_instance",
				Name: "example",
				Instances: []parser.ResourceInstance{
					{
						Attributes: map[string]interface{}{
							"instance_type": "t3.micro",
						},
					},
				},
			},
		},
	}
	
	stateData, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("Failed to marshal state: %v", err)
	}
	
	if err := os.WriteFile(stateFile, stateData, 0644); err != nil {
		t.Fatalf("Failed to write state file: %v", err)
	}

	// Skip this test as it requires mocking os.Exit
	t.Skip("Skipping test that requires mocking os.Exit")
}

func TestAnalyzeStateFile_InvalidStateFile(t *testing.T) {
	// Create a temporary invalid state file
	tempDir := t.TempDir()
	stateFile := filepath.Join(tempDir, "invalid.tfstate")
	
	if err := os.WriteFile(stateFile, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to write invalid state file: %v", err)
	}

	// Skip this test as it requires mocking os.Exit
	t.Skip("Skipping test that requires mocking os.Exit")
}