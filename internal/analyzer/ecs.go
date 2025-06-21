package analyzer

import "github.com/suer/tf-arm/internal/parser"

type ECSAnalyzer struct{}

func (a *ECSAnalyzer) SupportedType() string {
	return "aws_ecs_task_definition"
}

func (a *ECSAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: true,
	}

	for _, instance := range resource.Instances {
		if cpuArch, exists := instance.Attributes["cpu_architecture"]; exists {
			if cpuArch == "ARM64" {
				analysis.CurrentArch = "ARM64"
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "既にARM64アーキテクチャを使用"
			} else {
				analysis.CurrentArch = "X86_64"
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "cpu_architectureをARM64に変更可能"
			}
		} else {
			analysis.CurrentArch = "X86_64 (デフォルト)"
			analysis.RecommendedArch = "ARM64"
			analysis.Notes = "cpu_architecture = \"ARM64\" を追加可能"
		}
	}
	return analysis
}
