package analyzer

import (
	"fmt"
	"strings"

	"github.com/suer/tf-arm/internal/parser"
)

type EC2Analyzer struct{}

func (a *EC2Analyzer) SupportedType() string {
	return "aws_instance"
}

func (a *EC2Analyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: false,
	}

	for _, instance := range resource.Instances {
		if instanceType, exists := instance.Attributes["instance_type"]; exists {
			instanceTypeStr, ok := instanceType.(string)
			if !ok {
				continue
			}
			analysis.CurrentArch = getArchFromInstanceType(instanceTypeStr)

			if isARM64InstanceType(instanceTypeStr) {
				analysis.ARM64Compatible = true
				analysis.AlreadyUsingARM64 = true
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "Already using ARM64 instance type"
			} else if hasARM64Alternative(instanceTypeStr) {
				analysis.ARM64Compatible = true
				analysis.RecommendedArch = getARM64Alternative(instanceTypeStr)
				analysis.Notes = fmt.Sprintf("Can migrate to ARM64 instance type %s", analysis.RecommendedArch)
			} else {
				analysis.Notes = "No ARM64 compatible instance type available"
			}
		}
	}
	return analysis
}

type LaunchTemplateAnalyzer struct{}

func (a *LaunchTemplateAnalyzer) SupportedType() string {
	return "aws_launch_template"
}

func (a *LaunchTemplateAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: false,
	}

	for _, instance := range resource.Instances {
		if instanceType, exists := instance.Attributes["instance_type"]; exists {
			instanceTypeStr, ok := instanceType.(string)
			if !ok {
				continue
			}
			analysis.CurrentArch = getArchFromInstanceType(instanceTypeStr)

			if isARM64InstanceType(instanceTypeStr) {
				analysis.ARM64Compatible = true
				analysis.AlreadyUsingARM64 = true
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "Already using ARM64 instance type"
			} else if hasARM64Alternative(instanceTypeStr) {
				analysis.ARM64Compatible = true
				analysis.RecommendedArch = getARM64Alternative(instanceTypeStr)
				analysis.Notes = fmt.Sprintf("Can migrate to ARM64 instance type %s", analysis.RecommendedArch)
			}
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
	_, exists := getX86ToArm64Map()[instanceType]
	return exists
}

func getARM64Alternative(instanceType string) string {
	return getX86ToArm64Map()[instanceType]
}

func getArchFromInstanceType(instanceType string) string {
	if isARM64InstanceType(instanceType) {
		return "ARM64"
	}
	return "X86_64"
}

func getX86ToArm64Map() map[string]string {
	return map[string]string{
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
}
