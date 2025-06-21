package analyzer

import "github.com/suer/tf-arm/internal/parser"

type ARM64Analysis struct {
	ResourceType      string
	ResourceName      string
	FullAddress       string
	CurrentArch       string
	ARM64Compatible   bool
	AlreadyUsingARM64 bool
	RecommendedArch   string
	Notes             string
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
	case "aws_ecs_service":
		analyzer = &FargateAnalyzer{}
	case "aws_lambda_function":
		analyzer = &LambdaAnalyzer{}
	case "aws_codebuild_project":
		analyzer = &CodeBuildAnalyzer{}
	case "aws_db_instance":
		analyzer = &RDSAnalyzer{}
	case "aws_rds_cluster":
		analyzer = &AuroraAnalyzer{}
	case "aws_elasticache_cluster":
		analyzer = &ElastiCacheAnalyzer{}
	case "aws_memorydb_cluster":
		analyzer = &MemoryDBAnalyzer{}
	case "aws_eks_node_group":
		analyzer = &EKSAnalyzer{}
	case "aws_emr_cluster":
		analyzer = &EMRAnalyzer{}
	case "aws_emrserverless_application":
		analyzer = &EMRServerlessAnalyzer{}
	case "aws_opensearch_domain":
		analyzer = &OpenSearchAnalyzer{}
	case "aws_msk_cluster":
		analyzer = &MSKAnalyzer{}
	case "aws_sagemaker_endpoint_configuration":
		analyzer = &SageMakerAnalyzer{}
	case "aws_gamelift_fleet":
		analyzer = &GameLiftAnalyzer{}
	default:
		return ARM64Analysis{
			ResourceType:    resource.Type,
			ResourceName:    resource.Name,
			FullAddress:     resource.GetFullAddress(),
			ARM64Compatible: false,
			Notes:           "Resource type not supported for ARM64 compatibility check",
		}
	}

	analysis := analyzer.Analyze(resource)
	analysis.FullAddress = resource.GetFullAddress()
	return analysis
}
