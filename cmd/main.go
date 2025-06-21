package main

import (
	"fmt"
	"os"

	"github.com/suer/tf-arm/internal/analyzer"
	"github.com/suer/tf-arm/internal/parser"
	"github.com/suer/tf-arm/internal/reporter"
)

func main() {
	fmt.Println("tf-arm: Terraform State ARM64 Analyzer")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Handle help flags
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage()
		os.Exit(0)
	}

	// Handle version flag
	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Println("tf-arm version 1.0.0")
		os.Exit(0)
	}

	stateFile := os.Args[1]
	fmt.Printf("Analyzing Terraform state file: %s\n", stateFile)
	fmt.Println("")

	state, err := parser.ParseStateFile(stateFile)
	if err != nil {
		fmt.Printf("Error parsing state file: %v\n", err)
		os.Exit(1)
	}

	rep := reporter.New()
	rep.PrintHeader(len(state.Resources))

	var arm64CompatibleCount int
	var totalAnalyzedCount int

	for _, resource := range state.Resources {
		if resource.Mode != "managed" {
			continue
		}

		analysis := analyzer.AnalyzeResource(resource)

		if analysis.Notes != "Resource type not supported for ARM64 compatibility check" {
			totalAnalyzedCount++
			rep.PrintAnalysis(analysis)

			if analysis.ARM64Compatible {
				arm64CompatibleCount++
			}
		}
	}

	rep.PrintSummary(totalAnalyzedCount, arm64CompatibleCount)
}

func printUsage() {
	fmt.Println("Usage: tf-arm [OPTIONS] <terraform-state-file>")
	fmt.Println("")
	fmt.Println("This tool analyzes Terraform state files to identify AWS resources")
	fmt.Println("that can be migrated to ARM64 architecture for cost optimization.")
	fmt.Println("")
	fmt.Println("Arguments:")
	fmt.Println("  <terraform-state-file>    Path to the Terraform state file to analyze")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help               Show this help message and exit")
	fmt.Println("  -v, --version            Show version information and exit")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  tf-arm terraform.tfstate")
	fmt.Println("  tf-arm infrastructure.tfstate")
	fmt.Println("")
	fmt.Println("Supported AWS Services:")
	fmt.Println("  - Amazon EC2 (aws_instance, aws_launch_template)")
	fmt.Println("  - AWS Lambda (aws_lambda_function)")
	fmt.Println("  - Amazon ECS (aws_ecs_task_definition, aws_ecs_service)")
	fmt.Println("  - Amazon RDS (aws_db_instance, aws_rds_cluster)")
	fmt.Println("  - Amazon ElastiCache (aws_elasticache_cluster)")
	fmt.Println("  - Amazon MemoryDB (aws_memorydb_cluster)")
	fmt.Println("  - Amazon EKS (aws_eks_node_group)")
	fmt.Println("  - Amazon EMR (aws_emr_cluster, aws_emrserverless_application)")
	fmt.Println("  - Amazon OpenSearch (aws_opensearch_domain)")
	fmt.Println("  - Amazon MSK (aws_msk_cluster)")
	fmt.Println("  - AWS CodeBuild (aws_codebuild_project)")
	fmt.Println("  - Amazon SageMaker (aws_sagemaker_endpoint_configuration)")
	fmt.Println("  - Amazon GameLift (aws_gamelift_fleet)")
}
