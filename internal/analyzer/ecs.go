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
				analysis.AlreadyUsingARM64 = true
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "Already using ARM64 architecture"
			} else {
				analysis.CurrentArch = "X86_64"
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "Can change cpu_architecture to ARM64"
			}
		} else {
			analysis.CurrentArch = "X86_64 (default)"
			analysis.RecommendedArch = "ARM64"
			analysis.Notes = "Can add cpu_architecture = \"ARM64\""
		}
	}
	return analysis
}
