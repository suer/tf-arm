package analyzer

import (
	"testing"

	"github.com/suer/tf-arm/internal/parser"
)

func TestAnalyzeResource(t *testing.T) {
	tests := []struct {
		name         string
		resource     parser.TerraformResource
		expectedType string
		expected     func(ARM64Analysis) bool
	}{
		{
			name: "aws_instance resource",
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
			expectedType: "aws_instance",
			expected: func(analysis ARM64Analysis) bool {
				return analysis.ResourceType == "aws_instance" &&
					analysis.ResourceName == "example" &&
					analysis.Supported == true &&
					analysis.ARM64Compatible == true
			},
		},
		{
			name: "aws_lambda_function resource",
			resource: parser.TerraformResource{
				Type: "aws_lambda_function",
				Name: "test_function",
				Instances: []parser.ResourceInstance{
					{
						Attributes: map[string]interface{}{
							"runtime": "python3.9",
						},
					},
				},
			},
			expectedType: "aws_lambda_function",
			expected: func(analysis ARM64Analysis) bool {
				return analysis.ResourceType == "aws_lambda_function" &&
					analysis.ResourceName == "test_function" &&
					analysis.Supported == true
			},
		},
		{
			name: "unsupported resource type",
			resource: parser.TerraformResource{
				Type: "aws_unsupported_resource",
				Name: "test",
				Instances: []parser.ResourceInstance{
					{
						Attributes: map[string]interface{}{
							"some_attribute": "value",
						},
					},
				},
			},
			expectedType: "aws_unsupported_resource",
			expected: func(analysis ARM64Analysis) bool {
				return analysis.ResourceType == "aws_unsupported_resource" &&
					analysis.ResourceName == "test" &&
					analysis.Supported == false &&
					analysis.ARM64Compatible == false &&
					analysis.Notes == "Resource type not supported for ARM64 compatibility check"
			},
		},
		{
			name: "aws_db_instance resource",
			resource: parser.TerraformResource{
				Type: "aws_db_instance",
				Name: "database",
				Instances: []parser.ResourceInstance{
					{
						Attributes: map[string]interface{}{
							"instance_class": "db.t3.micro",
						},
					},
				},
			},
			expectedType: "aws_db_instance",
			expected: func(analysis ARM64Analysis) bool {
				return analysis.ResourceType == "aws_db_instance" &&
					analysis.ResourceName == "database" &&
					analysis.Supported == true
			},
		},
		{
			name: "aws_ecs_task_definition resource",
			resource: parser.TerraformResource{
				Type: "aws_ecs_task_definition",
				Name: "task",
				Instances: []parser.ResourceInstance{
					{
						Attributes: map[string]interface{}{
							"cpu_architecture": "X86_64",
						},
					},
				},
			},
			expectedType: "aws_ecs_task_definition",
			expected: func(analysis ARM64Analysis) bool {
				return analysis.ResourceType == "aws_ecs_task_definition" &&
					analysis.ResourceName == "task" &&
					analysis.Supported == true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := AnalyzeResource(tt.resource)

			if analysis.ResourceType != tt.expectedType {
				t.Errorf("AnalyzeResource() ResourceType = %v, want %v", analysis.ResourceType, tt.expectedType)
			}

			if analysis.FullAddress == "" {
				t.Error("AnalyzeResource() should set FullAddress")
			}

			if !tt.expected(analysis) {
				t.Errorf("AnalyzeResource() failed validation for %s", tt.name)
			}
		})
	}
}

func TestAnalyzeResource_AllSupportedTypes(t *testing.T) {
	supportedTypes := []string{
		"aws_instance",
		"aws_launch_template",
		"aws_ecs_task_definition",
		"aws_ecs_service",
		"aws_lambda_function",
		"aws_codebuild_project",
		"aws_db_instance",
		"aws_rds_cluster",
		"aws_elasticache_cluster",
		"aws_memorydb_cluster",
		"aws_eks_node_group",
		"aws_emr_cluster",
		"aws_emrserverless_application",
		"aws_opensearch_domain",
		"aws_msk_cluster",
		"aws_sagemaker_endpoint_configuration",
		"aws_gamelift_fleet",
	}

	for _, resourceType := range supportedTypes {
		t.Run(resourceType, func(t *testing.T) {
			resource := parser.TerraformResource{
				Type: resourceType,
				Name: "test",
				Instances: []parser.ResourceInstance{
					{
						Attributes: map[string]interface{}{
							"test_attribute": "value",
						},
					},
				},
			}

			analysis := AnalyzeResource(resource)

			if !analysis.Supported {
				t.Errorf("AnalyzeResource() should support resource type %s", resourceType)
			}

			if analysis.ResourceType != resourceType {
				t.Errorf("AnalyzeResource() ResourceType = %v, want %v", analysis.ResourceType, resourceType)
			}

			if analysis.ResourceName != "test" {
				t.Errorf("AnalyzeResource() ResourceName = %v, want %v", analysis.ResourceName, "test")
			}

			if analysis.FullAddress == "" {
				t.Error("AnalyzeResource() should set FullAddress")
			}
		})
	}
}

func TestAnalyzeResource_WithModule(t *testing.T) {
	resource := parser.TerraformResource{
		Type:   "aws_instance",
		Name:   "example",
		Module: "module.web",
		Instances: []parser.ResourceInstance{
			{
				Attributes: map[string]interface{}{
					"instance_type": "t3.micro",
				},
			},
		},
	}

	analysis := AnalyzeResource(resource)

	expectedAddress := "module.web.aws_instance.example"
	if analysis.FullAddress != expectedAddress {
		t.Errorf("AnalyzeResource() FullAddress = %v, want %v", analysis.FullAddress, expectedAddress)
	}
}