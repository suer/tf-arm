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
	arm64Prefixes := []string{
		// Graviton1
		"a1.",
		// Graviton2
		"t4g.", "m6g.", "m6gd.", "c6g.", "c6gd.", "c6gn.", "r6g.", "r6gd.", "x2gd.",
		// Graviton3
		"c7g.", "c7gd.", "c7gn.", "m7g.", "m7gd.", "r7g.", "r7gd.", "hpc7g.",
		// Graviton4
		"c8g.", "m8g.", "r8g.", "x8g.", "i8g.",
	}
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
		// T3 -> T4g (Graviton2)
		"t3.nano":     "t4g.nano",
		"t3.micro":    "t4g.micro",
		"t3.small":    "t4g.small",
		"t3.medium":   "t4g.medium",
		"t3.large":    "t4g.large",
		"t3.xlarge":   "t4g.xlarge",
		"t3.2xlarge":  "t4g.2xlarge",
		// M5 -> M7g (Graviton3 - better performance than M6g)
		"m5.large":    "m7g.large",
		"m5.xlarge":   "m7g.xlarge",
		"m5.2xlarge":  "m7g.2xlarge",
		"m5.4xlarge":  "m7g.4xlarge",
		"m5.8xlarge":  "m7g.8xlarge",
		"m5.12xlarge": "m7g.12xlarge",
		"m5.16xlarge": "m7g.16xlarge",
		// M6i -> M8g (Graviton4 - latest generation)
		"m6i.large":    "m8g.large",
		"m6i.xlarge":   "m8g.xlarge",
		"m6i.2xlarge":  "m8g.2xlarge",
		"m6i.4xlarge":  "m8g.4xlarge",
		"m6i.8xlarge":  "m8g.8xlarge",
		"m6i.12xlarge": "m8g.12xlarge",
		"m6i.16xlarge": "m8g.16xlarge",
		"m6i.24xlarge": "m8g.24xlarge",
		"m6i.32xlarge": "m8g.32xlarge",
		"m6i.48xlarge": "m8g.48xlarge",
		// C5 -> C7g (Graviton3 - better performance than C6g)
		"c5.large":    "c7g.large",
		"c5.xlarge":   "c7g.xlarge",
		"c5.2xlarge":  "c7g.2xlarge",
		"c5.4xlarge":  "c7g.4xlarge",
		"c5.9xlarge":  "c7g.9xlarge",
		"c5.12xlarge": "c7g.12xlarge",
		"c5.18xlarge": "c7g.16xlarge",
		// C6i -> C8g (Graviton4 - latest generation)
		"c6i.large":    "c8g.large",
		"c6i.xlarge":   "c8g.xlarge",
		"c6i.2xlarge":  "c8g.2xlarge",
		"c6i.4xlarge":  "c8g.4xlarge",
		"c6i.8xlarge":  "c8g.8xlarge",
		"c6i.12xlarge": "c8g.12xlarge",
		"c6i.16xlarge": "c8g.16xlarge",
		"c6i.24xlarge": "c8g.24xlarge",
		"c6i.32xlarge": "c8g.32xlarge",
		"c6i.48xlarge": "c8g.48xlarge",
		// R5 -> R7g (Graviton3 - better performance than R6g)
		"r5.large":    "r7g.large",
		"r5.xlarge":   "r7g.xlarge",
		"r5.2xlarge":  "r7g.2xlarge",
		"r5.4xlarge":  "r7g.4xlarge",
		"r5.8xlarge":  "r7g.8xlarge",
		"r5.12xlarge": "r7g.12xlarge",
		"r5.16xlarge": "r7g.16xlarge",
		// R6i -> R8g (Graviton4 - latest generation)
		"r6i.large":    "r8g.large",
		"r6i.xlarge":   "r8g.xlarge",
		"r6i.2xlarge":  "r8g.2xlarge",
		"r6i.4xlarge":  "r8g.4xlarge",
		"r6i.8xlarge":  "r8g.8xlarge",
		"r6i.12xlarge": "r8g.12xlarge",
		"r6i.16xlarge": "r8g.16xlarge",
		"r6i.24xlarge": "r8g.24xlarge",
		"r6i.32xlarge": "r8g.32xlarge",
		"r6i.48xlarge": "r8g.48xlarge",
	}
}
