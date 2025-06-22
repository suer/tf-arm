package analyzer

import (
	"strings"
	"testing"

	"github.com/suer/tf-arm/internal/parser"
)

func TestEC2Analyzer_SupportedType(t *testing.T) {
	analyzer := &EC2Analyzer{}
	expected := "aws_instance"
	
	if analyzer.SupportedType() != expected {
		t.Errorf("SupportedType() = %v, want %v", analyzer.SupportedType(), expected)
	}
}

func TestEC2Analyzer_Analyze(t *testing.T) {
	tests := []struct {
		name         string
		resource     parser.TerraformResource
		expectARM64  bool
		expectUsing  bool
		expectNotes  string
	}{
		{
			name: "x86 instance with ARM64 alternative",
			resource: parser.TerraformResource{
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
			expectARM64: true,
			expectUsing: false,
			expectNotes: "Can migrate to ARM64 instance type t4g.micro",
		},
		{
			name: "ARM64 instance already in use",
			resource: parser.TerraformResource{
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
			expectARM64: true,
			expectUsing: true,
			expectNotes: "Already using ARM64 instance type",
		},
		{
			name: "x86 instance without ARM64 alternative",
			resource: parser.TerraformResource{
				Type: "aws_instance",
				Name: "legacy",
				Instances: []parser.ResourceInstance{
					{
						Attributes: map[string]interface{}{
							"instance_type": "t2.micro",
						},
					},
				},
			},
			expectARM64: false,
			expectUsing: false,
			expectNotes: "No ARM64 compatible instance type available",
		},
		{
			name: "resource without instance_type",
			resource: parser.TerraformResource{
				Type: "aws_instance",
				Name: "no_type",
				Instances: []parser.ResourceInstance{
					{
						Attributes: map[string]interface{}{
							"ami": "ami-12345",
						},
					},
				},
			},
			expectARM64: false,
			expectUsing: false,
			expectNotes: "",
		},
		{
			name: "resource with non-string instance_type",
			resource: parser.TerraformResource{
				Type: "aws_instance",
				Name: "invalid_type",
				Instances: []parser.ResourceInstance{
					{
						Attributes: map[string]interface{}{
							"instance_type": 123,
						},
					},
				},
			},
			expectARM64: false,
			expectUsing: false,
			expectNotes: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := &EC2Analyzer{}
			analysis := analyzer.Analyze(tt.resource)

			if analysis.ResourceType != "aws_instance" {
				t.Errorf("Analyze() ResourceType = %v, want aws_instance", analysis.ResourceType)
			}

			if analysis.ResourceName != tt.resource.Name {
				t.Errorf("Analyze() ResourceName = %v, want %v", analysis.ResourceName, tt.resource.Name)
			}

			if analysis.ARM64Compatible != tt.expectARM64 {
				t.Errorf("Analyze() ARM64Compatible = %v, want %v", analysis.ARM64Compatible, tt.expectARM64)
			}

			if analysis.AlreadyUsingARM64 != tt.expectUsing {
				t.Errorf("Analyze() AlreadyUsingARM64 = %v, want %v", analysis.AlreadyUsingARM64, tt.expectUsing)
			}

			if tt.expectNotes != "" && !strings.Contains(analysis.Notes, tt.expectNotes) {
				t.Errorf("Analyze() Notes = %v, want to contain %v", analysis.Notes, tt.expectNotes)
			}
		})
	}
}

func TestLaunchTemplateAnalyzer_SupportedType(t *testing.T) {
	analyzer := &LaunchTemplateAnalyzer{}
	expected := "aws_launch_template"
	
	if analyzer.SupportedType() != expected {
		t.Errorf("SupportedType() = %v, want %v", analyzer.SupportedType(), expected)
	}
}

func TestLaunchTemplateAnalyzer_Analyze(t *testing.T) {
	tests := []struct {
		name         string
		resource     parser.TerraformResource
		expectARM64  bool
		expectUsing  bool
	}{
		{
			name: "launch template with x86 instance",
			resource: parser.TerraformResource{
				Type: "aws_launch_template",
				Name: "example",
				Instances: []parser.ResourceInstance{
					{
						Attributes: map[string]interface{}{
							"instance_type": "m5.large",
						},
					},
				},
			},
			expectARM64: true,
			expectUsing: false,
		},
		{
			name: "launch template with ARM64 instance",
			resource: parser.TerraformResource{
				Type: "aws_launch_template",
				Name: "arm_example",
				Instances: []parser.ResourceInstance{
					{
						Attributes: map[string]interface{}{
							"instance_type": "m6g.large",
						},
					},
				},
			},
			expectARM64: true,
			expectUsing: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := &LaunchTemplateAnalyzer{}
			analysis := analyzer.Analyze(tt.resource)

			if analysis.ResourceType != "aws_launch_template" {
				t.Errorf("Analyze() ResourceType = %v, want aws_launch_template", analysis.ResourceType)
			}

			if analysis.ARM64Compatible != tt.expectARM64 {
				t.Errorf("Analyze() ARM64Compatible = %v, want %v", analysis.ARM64Compatible, tt.expectARM64)
			}

			if analysis.AlreadyUsingARM64 != tt.expectUsing {
				t.Errorf("Analyze() AlreadyUsingARM64 = %v, want %v", analysis.AlreadyUsingARM64, tt.expectUsing)
			}
		})
	}
}

func TestIsARM64InstanceType(t *testing.T) {
	tests := []struct {
		instanceType string
		expected     bool
	}{
		{"a1.medium", true},
		{"t4g.micro", true},
		{"m6g.large", true},
		{"m6gd.xlarge", true},
		{"c6g.2xlarge", true},
		{"c6gd.4xlarge", true},
		{"c6gn.large", true},
		{"r6g.xlarge", true},
		{"r6gd.2xlarge", true},
		{"x2gd.medium", true},
		{"t3.micro", false},
		{"m5.large", false},
		{"c5.xlarge", false},
		{"r5.2xlarge", false},
		{"t2.micro", false},
	}

	for _, tt := range tests {
		t.Run(tt.instanceType, func(t *testing.T) {
			result := isARM64InstanceType(tt.instanceType)
			if result != tt.expected {
				t.Errorf("isARM64InstanceType(%v) = %v, want %v", tt.instanceType, result, tt.expected)
			}
		})
	}
}

func TestHasARM64Alternative(t *testing.T) {
	tests := []struct {
		instanceType string
		expected     bool
	}{
		{"t3.micro", true},
		{"t3.small", true},
		{"m5.large", true},
		{"c5.xlarge", true},
		{"r5.2xlarge", true},
		{"t2.micro", false},
		{"i3.large", false},
		{"unknown.type", false},
	}

	for _, tt := range tests {
		t.Run(tt.instanceType, func(t *testing.T) {
			result := hasARM64Alternative(tt.instanceType)
			if result != tt.expected {
				t.Errorf("hasARM64Alternative(%v) = %v, want %v", tt.instanceType, result, tt.expected)
			}
		})
	}
}

func TestGetARM64Alternative(t *testing.T) {
	tests := []struct {
		instanceType string
		expected     string
	}{
		{"t3.micro", "t4g.micro"},
		{"t3.small", "t4g.small"},
		{"m5.large", "m6g.large"},
		{"c5.xlarge", "c6g.xlarge"},
		{"r5.2xlarge", "r6g.2xlarge"},
	}

	for _, tt := range tests {
		t.Run(tt.instanceType, func(t *testing.T) {
			result := getARM64Alternative(tt.instanceType)
			if result != tt.expected {
				t.Errorf("getARM64Alternative(%v) = %v, want %v", tt.instanceType, result, tt.expected)
			}
		})
	}
}

func TestGetArchFromInstanceType(t *testing.T) {
	tests := []struct {
		instanceType string
		expected     string
	}{
		{"t4g.micro", "ARM64"},
		{"m6g.large", "ARM64"},
		{"c6g.xlarge", "ARM64"},
		{"t3.micro", "X86_64"},
		{"m5.large", "X86_64"},
		{"c5.xlarge", "X86_64"},
	}

	for _, tt := range tests {
		t.Run(tt.instanceType, func(t *testing.T) {
			result := getArchFromInstanceType(tt.instanceType)
			if result != tt.expected {
				t.Errorf("getArchFromInstanceType(%v) = %v, want %v", tt.instanceType, result, tt.expected)
			}
		})
	}
}

func TestGetX86ToArm64Map(t *testing.T) {
	mapping := getX86ToArm64Map()
	
	if len(mapping) == 0 {
		t.Error("getX86ToArm64Map() should return non-empty map")
	}

	expectedMappings := map[string]string{
		"t3.micro":    "t4g.micro",
		"m5.large":    "m6g.large",
		"c5.xlarge":   "c6g.xlarge",
		"r5.2xlarge":  "r6g.2xlarge",
	}

	for x86Type, expectedArm64Type := range expectedMappings {
		if arm64Type, exists := mapping[x86Type]; !exists {
			t.Errorf("getX86ToArm64Map() missing mapping for %v", x86Type)
		} else if arm64Type != expectedArm64Type {
			t.Errorf("getX86ToArm64Map()[%v] = %v, want %v", x86Type, arm64Type, expectedArm64Type)
		}
	}
}