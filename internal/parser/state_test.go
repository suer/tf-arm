package parser

import (
	"encoding/json"
	"os"
	"testing"
)

func TestTerraformResource_GetFullAddress(t *testing.T) {
	tests := []struct {
		name     string
		resource TerraformResource
		expected string
	}{
		{
			name: "resource without module",
			resource: TerraformResource{
				Type: "aws_instance",
				Name: "example",
			},
			expected: "aws_instance.example",
		},
		{
			name: "resource with module",
			resource: TerraformResource{
				Type:   "aws_instance",
				Name:   "example",
				Module: "module.web",
			},
			expected: "module.web.aws_instance.example",
		},
		{
			name: "resource with empty module",
			resource: TerraformResource{
				Type:   "aws_s3_bucket",
				Name:   "bucket",
				Module: "",
			},
			expected: "aws_s3_bucket.bucket",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.resource.GetFullAddress()
			if result != tt.expected {
				t.Errorf("GetFullAddress() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseStateFile(t *testing.T) {
	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid state file",
			setupFile: func(t *testing.T) string {
				state := TerraformState{
					Version: 4,
					Resources: []TerraformResource{
						{
							Mode:     "managed",
							Type:     "aws_instance",
							Name:     "example",
							Provider: "provider[\"registry.terraform.io/hashicorp/aws\"]",
							Instances: []ResourceInstance{
								{
									Attributes: map[string]interface{}{
										"id":            "i-1234567890abcdef0",
										"instance_type": "t2.micro",
									},
								},
							},
						},
					},
				}
				return createTempStateFile(t, state)
			},
			expectError: false,
		},
		{
			name: "empty filename",
			setupFile: func(t *testing.T) string {
				return ""
			},
			expectError: true,
			errorMsg:    "filename cannot be empty",
		},
		{
			name: "non-existent file",
			setupFile: func(t *testing.T) string {
				return "/non/existent/file.json"
			},
			expectError: true,
			errorMsg:    "failed to stat state file",
		},
		{
			name: "empty file",
			setupFile: func(t *testing.T) string {
				tmpFile, err := os.CreateTemp("", "terraform-state-*.json")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				tmpFile.Close()
				return tmpFile.Name()
			},
			expectError: true,
			errorMsg:    "state file is empty",
		},
		{
			name: "invalid JSON",
			setupFile: func(t *testing.T) string {
				tmpFile, err := os.CreateTemp("", "terraform-state-*.json")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				tmpFile.WriteString("invalid json content")
				tmpFile.Close()
				return tmpFile.Name()
			},
			expectError: true,
			errorMsg:    "failed to parse JSON",
		},
		{
			name: "state with version zero",
			setupFile: func(t *testing.T) string {
				state := TerraformState{
					Version:   0,
					Resources: []TerraformResource{},
				}
				return createTempStateFile(t, state)
			},
			expectError: true,
			errorMsg:    "invalid state file: version is missing or zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := tt.setupFile(t)
			
			// Clean up temp file if it exists
			if filename != "" && filename != "/non/existent/file.json" {
				defer os.Remove(filename)
			}

			result, err := ParseStateFile(filename)

			if tt.expectError {
				if err == nil {
					t.Errorf("ParseStateFile() expected error but got none")
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("ParseStateFile() error = %v, want error containing %v", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ParseStateFile() unexpected error = %v", err)
					return
				}
				if result == nil {
					t.Errorf("ParseStateFile() returned nil result")
					return
				}
				if result.Version == 0 {
					t.Errorf("ParseStateFile() returned state with version 0")
				}
			}
		})
	}
}

func createTempStateFile(t *testing.T, state TerraformState) string {
	tmpFile, err := os.CreateTemp("", "terraform-state-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	
	encoder := json.NewEncoder(tmpFile)
	if err := encoder.Encode(state); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		t.Fatalf("Failed to write state to temp file: %v", err)
	}
	
	tmpFile.Close()
	return tmpFile.Name()
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && containsAt(s, substr, 0)))
}

func containsAt(s, substr string, start int) bool {
	if start+len(substr) > len(s) {
		return false
	}
	for i := 0; i < len(substr); i++ {
		if s[start+i] != substr[i] {
			if start+1 < len(s) {
				return containsAt(s, substr, start+1)
			}
			return false
		}
	}
	return true
}