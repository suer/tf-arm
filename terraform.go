package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type TerraformState struct {
	Version   int                    `json:"version"`
	Resources []TerraformResource    `json:"resources"`
}

type TerraformResource struct {
	Mode      string                 `json:"mode"`
	Type      string                 `json:"type"`
	Name      string                 `json:"name"`
	Provider  string                 `json:"provider"`
	Instances []ResourceInstance     `json:"instances"`
}

type ResourceInstance struct {
	Attributes map[string]interface{} `json:"attributes"`
}

type ARM64Analysis struct {
	ResourceType     string
	ResourceName     string
	CurrentArch      string
	ARM64Compatible  bool
	RecommendedArch  string
	Notes           string
}

func parseStateFile(filename string) (*TerraformState, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state TerraformState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &state, nil
}

func analyzeARM64Compatibility(resource TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType: resource.Type,
		ResourceName: resource.Name,
		ARM64Compatible: false,
	}

	switch resource.Type {
	case "aws_instance":
		analysis = analyzeEC2Instance(resource, analysis)
	case "aws_launch_template":
		analysis = analyzeLaunchTemplate(resource, analysis)
	case "aws_autoscaling_group":
		analysis = analyzeAutoScalingGroup(resource, analysis)
	case "aws_ecs_task_definition":
		analysis = analyzeECSTaskDefinition(resource, analysis)
	case "aws_lambda_function":
		analysis = analyzeLambdaFunction(resource, analysis)
	default:
		analysis.Notes = "リソースタイプはARM64互換性チェック対象外"
	}

	return analysis
}

func analyzeEC2Instance(resource TerraformResource, analysis ARM64Analysis) ARM64Analysis {
	for _, instance := range resource.Instances {
		if instanceType, exists := instance.Attributes["instance_type"]; exists {
			instanceTypeStr := instanceType.(string)
			analysis.CurrentArch = getArchFromInstanceType(instanceTypeStr)
			
			if isARM64InstanceType(instanceTypeStr) {
				analysis.ARM64Compatible = true
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "既にARM64インスタンスタイプを使用"
			} else if hasARM64Alternative(instanceTypeStr) {
				analysis.ARM64Compatible = true
				analysis.RecommendedArch = getARM64Alternative(instanceTypeStr)
				analysis.Notes = fmt.Sprintf("ARM64インスタンスタイプ %s に変更可能", analysis.RecommendedArch)
			} else {
				analysis.Notes = "ARM64対応インスタンスタイプなし"
			}
		}
	}
	return analysis
}

func analyzeLaunchTemplate(resource TerraformResource, analysis ARM64Analysis) ARM64Analysis {
	for _, instance := range resource.Instances {
		if instanceType, exists := instance.Attributes["instance_type"]; exists {
			instanceTypeStr := instanceType.(string)
			analysis.CurrentArch = getArchFromInstanceType(instanceTypeStr)
			
			if isARM64InstanceType(instanceTypeStr) {
				analysis.ARM64Compatible = true
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "既にARM64インスタンスタイプを使用"
			} else if hasARM64Alternative(instanceTypeStr) {
				analysis.ARM64Compatible = true
				analysis.RecommendedArch = getARM64Alternative(instanceTypeStr)
				analysis.Notes = fmt.Sprintf("ARM64インスタンスタイプ %s に変更可能", analysis.RecommendedArch)
			}
		}
	}
	return analysis
}

func analyzeAutoScalingGroup(resource TerraformResource, analysis ARM64Analysis) ARM64Analysis {
	analysis.Notes = "Launch TemplateまたはLaunch Configurationを確認してください"
	return analysis
}

func analyzeECSTaskDefinition(resource TerraformResource, analysis ARM64Analysis) ARM64Analysis {
	for _, instance := range resource.Instances {
		if cpuArch, exists := instance.Attributes["cpu_architecture"]; exists {
			if cpuArch == "ARM64" {
				analysis.ARM64Compatible = true
				analysis.CurrentArch = "ARM64"
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "既にARM64アーキテクチャを使用"
			} else {
				analysis.ARM64Compatible = true
				analysis.CurrentArch = "X86_64"
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "cpu_architectureをARM64に変更可能"
			}
		} else {
			analysis.ARM64Compatible = true
			analysis.CurrentArch = "X86_64 (デフォルト)"
			analysis.RecommendedArch = "ARM64"
			analysis.Notes = "cpu_architecture = \"ARM64\" を追加可能"
		}
	}
	return analysis
}

func analyzeLambdaFunction(resource TerraformResource, analysis ARM64Analysis) ARM64Analysis {
	for _, instance := range resource.Instances {
		if architectures, exists := instance.Attributes["architectures"]; exists {
			archList := architectures.([]interface{})
			if len(archList) > 0 && archList[0] == "arm64" {
				analysis.ARM64Compatible = true
				analysis.CurrentArch = "ARM64"
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "既にARM64アーキテクチャを使用"
			} else {
				analysis.ARM64Compatible = true
				analysis.CurrentArch = "X86_64"
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "architectures = [\"arm64\"] に変更可能"
			}
		} else {
			analysis.ARM64Compatible = true
			analysis.CurrentArch = "X86_64 (デフォルト)"
			analysis.RecommendedArch = "ARM64"
			analysis.Notes = "architectures = [\"arm64\"] を追加可能"
		}
	}
	return analysis
}

func isARM64InstanceType(instanceType string) bool {
	arm64Prefixes := []string{"a1.", "t4g.", "m6g.", "m6gd.", "c6g.", "c6gd.", "c6gn.", "r6g.", "r6gd.", "x2gd."}
	for _, prefix := range arm64Prefixes {
		if strings.HasPrefix(instanceType, prefix) {
			return true
		}
	}
	return false
}

func hasARM64Alternative(instanceType string) bool {
	x86ToArm64Map := map[string]string{
		"t3.nano":     "t4g.nano",
		"t3.micro":    "t4g.micro",
		"t3.small":    "t4g.small",
		"t3.medium":   "t4g.medium",
		"t3.large":    "t4g.large",
		"t3.xlarge":   "t4g.xlarge",
		"t3.2xlarge":  "t4g.2xlarge",
		"m5.large":    "m6g.large",
		"m5.xlarge":   "m6g.xlarge",
		"m5.2xlarge":  "m6g.2xlarge",
		"m5.4xlarge":  "m6g.4xlarge",
		"m5.8xlarge":  "m6g.8xlarge",
		"m5.12xlarge": "m6g.12xlarge",
		"m5.16xlarge": "m6g.16xlarge",
		"c5.large":    "c6g.large",
		"c5.xlarge":   "c6g.xlarge",
		"c5.2xlarge":  "c6g.2xlarge",
		"c5.4xlarge":  "c6g.4xlarge",
		"c5.9xlarge":  "c6g.9xlarge",
		"c5.12xlarge": "c6g.12xlarge",
		"c5.18xlarge": "c6g.16xlarge",
		"r5.large":    "r6g.large",
		"r5.xlarge":   "r6g.xlarge",
		"r5.2xlarge":  "r6g.2xlarge",
		"r5.4xlarge":  "r6g.4xlarge",
		"r5.8xlarge":  "r6g.8xlarge",
		"r5.12xlarge": "r6g.12xlarge",
		"r5.16xlarge": "r6g.16xlarge",
	}
	
	_, exists := x86ToArm64Map[instanceType]
	return exists
}

func getARM64Alternative(instanceType string) string {
	x86ToArm64Map := map[string]string{
		"t3.nano":     "t4g.nano",
		"t3.micro":    "t4g.micro",
		"t3.small":    "t4g.small",
		"t3.medium":   "t4g.medium",
		"t3.large":    "t4g.large",
		"t3.xlarge":   "t4g.xlarge",
		"t3.2xlarge":  "t4g.2xlarge",
		"m5.large":    "m6g.large",
		"m5.xlarge":   "m6g.xlarge",
		"m5.2xlarge":  "m6g.2xlarge",
		"m5.4xlarge":  "m6g.4xlarge",
		"m5.8xlarge":  "m6g.8xlarge",
		"m5.12xlarge": "m6g.12xlarge",
		"m5.16xlarge": "m6g.16xlarge",
		"c5.large":    "c6g.large",
		"c5.xlarge":   "c6g.xlarge",
		"c5.2xlarge":  "c6g.2xlarge",
		"c5.4xlarge":  "c6g.4xlarge",
		"c5.9xlarge":  "c6g.9xlarge",
		"c5.12xlarge": "c6g.12xlarge",
		"c5.18xlarge": "c6g.16xlarge",
		"r5.large":    "r6g.large",
		"r5.xlarge":   "r6g.xlarge",
		"r5.2xlarge":  "r6g.2xlarge",
		"r5.4xlarge":  "r6g.4xlarge",
		"r5.8xlarge":  "r6g.8xlarge",
		"r5.12xlarge": "r6g.12xlarge",
		"r5.16xlarge": "r6g.16xlarge",
	}
	
	return x86ToArm64Map[instanceType]
}

func getArchFromInstanceType(instanceType string) string {
	if isARM64InstanceType(instanceType) {
		return "ARM64"
	}
	return "X86_64"
}