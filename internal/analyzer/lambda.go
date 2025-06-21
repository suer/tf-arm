package analyzer

import "github.com/suer/tf-arm/internal/parser"

type LambdaAnalyzer struct{}

func (a *LambdaAnalyzer) SupportedType() string {
	return "aws_lambda_function"
}

func (a *LambdaAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: true,
	}

	for _, instance := range resource.Instances {
		if architectures, exists := instance.Attributes["architectures"]; exists {
			archList := architectures.([]any)
			if len(archList) > 0 && archList[0] == "arm64" {
				analysis.CurrentArch = "ARM64"
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "Already using ARM64 architecture"
			} else {
				analysis.CurrentArch = "X86_64"
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "Can change architectures to [\"arm64\"]"
			}
		} else {
			analysis.CurrentArch = "X86_64 (default)"
			analysis.RecommendedArch = "ARM64"
			analysis.Notes = "Can add architectures = [\"arm64\"]"
		}
	}
	return analysis
}
