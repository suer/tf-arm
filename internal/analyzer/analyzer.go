package analyzer

import "github.com/suer/tf-arm/internal/parser"

type ARM64Analysis struct {
	ResourceType    string
	ResourceName    string
	CurrentArch     string
	ARM64Compatible bool
	RecommendedArch string
	Notes           string
}

type Analyzer interface {
	Analyze(resource parser.TerraformResource) ARM64Analysis
	SupportedType() string
}

func AnalyzeResource(resource parser.TerraformResource) ARM64Analysis {
	var analyzer Analyzer

	switch resource.Type {
	case "aws_instance":
		analyzer = &EC2Analyzer{}
	case "aws_launch_template":
		analyzer = &LaunchTemplateAnalyzer{}
	case "aws_ecs_task_definition":
		analyzer = &ECSAnalyzer{}
	case "aws_lambda_function":
		analyzer = &LambdaAnalyzer{}
	case "aws_codebuild_project":
		analyzer = &CodeBuildAnalyzer{}
	default:
		return ARM64Analysis{
			ResourceType:    resource.Type,
			ResourceName:    resource.Name,
			ARM64Compatible: false,
			Notes:           "Resource type not supported for ARM64 compatibility check",
		}
	}

	return analyzer.Analyze(resource)
}
