package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/suer/tf-arm/internal/analyzer"
	"github.com/suer/tf-arm/internal/parser"
	"github.com/suer/tf-arm/internal/reporter"
)

var version = "1.0.0"
var showVersion bool

var rootCmd = &cobra.Command{
	Use:   "tf-arm [state-file]",
	Short: "Terraform State ARM64 Analyzer",
	Long: `tf-arm analyzes Terraform state files to identify AWS resources
that can be migrated to ARM64 architecture for cost optimization.

Supported AWS Services:
  - Amazon EC2 (aws_instance, aws_launch_template)
  - AWS Lambda (aws_lambda_function)
  - Amazon ECS (aws_ecs_task_definition, aws_ecs_service)
  - Amazon RDS (aws_db_instance, aws_rds_cluster)
  - Amazon ElastiCache (aws_elasticache_cluster)
  - Amazon MemoryDB (aws_memorydb_cluster)
  - Amazon EKS (aws_eks_node_group)
  - Amazon EMR (aws_emr_cluster, aws_emrserverless_application)
  - Amazon OpenSearch (aws_opensearch_domain)
  - Amazon MSK (aws_msk_cluster)
  - AWS CodeBuild (aws_codebuild_project)
  - Amazon SageMaker (aws_sagemaker_endpoint_configuration)
  - Amazon GameLift (aws_gamelift_fleet)`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Printf("tf-arm version %s\n", version)
			return
		}

		if len(args) == 0 {
			cmd.Help()
			return
		}

		stateFile := args[0]
		analyzeStateFile(stateFile)
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show version information")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func analyzeStateFile(stateFile string) {
	fmt.Println("tf-arm: Terraform State ARM64 Analyzer")
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
