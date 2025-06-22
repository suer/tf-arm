package analyzer

import (
	"slices"

	"github.com/suer/tf-arm/internal/parser"
)

type CodeBuildAnalyzer struct{}

func (a *CodeBuildAnalyzer) SupportedType() string {
	return "aws_codebuild_project"
}

func (a *CodeBuildAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: true,
	}

	for _, instance := range resource.Instances {
		if environment, exists := instance.Attributes["environment"]; exists {
			envList, ok := environment.([]any)
			if !ok {
				continue
			}
			if len(envList) > 0 {
				env, ok := envList[0].(map[string]any)
				if !ok {
					continue
				}
				if computeType, exists := env["compute_type"]; exists {
					computeTypeStr, ok := computeType.(string)
					if !ok {
						continue
					}
					if isARM64ComputeType(computeTypeStr) {
						analysis.CurrentArch = "ARM64"
						analysis.AlreadyUsingARM64 = true
						analysis.RecommendedArch = "ARM64"
						analysis.Notes = "Already using ARM64 compute type"
					} else if hasARM64ComputeTypeAlternative(computeTypeStr) {
						analysis.CurrentArch = "X86_64"
						analysis.RecommendedArch = getARM64ComputeTypeAlternative(computeTypeStr)
						analysis.Notes = "Can migrate to ARM64 compute type: " + analysis.RecommendedArch
					} else {
						analysis.CurrentArch = "X86_64"
						analysis.ARM64Compatible = false
						analysis.Notes = "No ARM64 compatible compute type available"
					}
				}
			}
		}
	}
	return analysis
}

func isARM64ComputeType(computeType string) bool {
	arm64Types := []string{
		"BUILD_GENERAL1_SMALL_ARM",
		"BUILD_GENERAL1_MEDIUM_ARM",
		"BUILD_GENERAL1_LARGE_ARM",
		"BUILD_GENERAL1_2XLARGE_ARM",
	}

	return slices.Contains(arm64Types, computeType)
}

func hasARM64ComputeTypeAlternative(computeType string) bool {
	_, exists := getX86ToArm64ComputeTypeMap()[computeType]
	return exists
}

func getARM64ComputeTypeAlternative(computeType string) string {
	return getX86ToArm64ComputeTypeMap()[computeType]
}

func getX86ToArm64ComputeTypeMap() map[string]string {
	return map[string]string{
		"BUILD_GENERAL1_SMALL":   "BUILD_GENERAL1_SMALL_ARM",
		"BUILD_GENERAL1_MEDIUM":  "BUILD_GENERAL1_MEDIUM_ARM",
		"BUILD_GENERAL1_LARGE":   "BUILD_GENERAL1_LARGE_ARM",
		"BUILD_GENERAL1_2XLARGE": "BUILD_GENERAL1_2XLARGE_ARM",
	}
}
