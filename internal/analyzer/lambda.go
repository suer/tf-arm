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
				analysis.Notes = "既にARM64アーキテクチャを使用"
			} else {
				analysis.CurrentArch = "X86_64"
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "architectures = [\"arm64\"] に変更可能"
			}
		} else {
			analysis.CurrentArch = "X86_64 (デフォルト)"
			analysis.RecommendedArch = "ARM64"
			analysis.Notes = "architectures = [\"arm64\"] を追加可能"
		}
	}
	return analysis
}
