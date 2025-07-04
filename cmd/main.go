package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/suer/tf-arm/internal/analyzer"
	"github.com/suer/tf-arm/internal/parser"
	"github.com/suer/tf-arm/internal/reporter"
)

var version = "dev"
var showVersion bool
var outputFormat string
var exitCode int

type JSONOutput struct {
	Summary struct {
		TotalAnalyzed      int     `json:"total_analyzed"`
		ARM64Compatible    int     `json:"arm64_compatible"`
		Migrateable        int     `json:"migrateable"`
		CompatibilityRate  float64 `json:"compatibility_rate"`
		MigrateablePercent float64 `json:"migrateable_percent"`
	} `json:"summary"`
	Resources []analyzer.ARM64Analysis `json:"resources"`
}

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
		analyzeStateFile(stateFile, outputFormat, exitCode)
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show version information")
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "Output format (text or json)")
	rootCmd.Flags().IntVar(&exitCode, "exit-code", 0, "Exit with specified code when ARM64 compatible resources are found")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func canMigrateToARM64(analysis analyzer.ARM64Analysis) bool {
	return analysis.ARM64Compatible && !analysis.AlreadyUsingARM64
}

func calculateMigrateablePercent(migrateableCount, arm64CompatibleCount int) float64 {
	if arm64CompatibleCount == 0 {
		return 0
	}
	return float64(migrateableCount) / float64(arm64CompatibleCount) * 100
}

func analyzeStateFile(stateFile, format string, exitCode int) {
	// Validate file exists and is accessible
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		fmt.Printf("Error: State file '%s' does not exist\n", stateFile)
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("Error accessing state file '%s': %v\n", stateFile, err)
		os.Exit(1)
	}

	state, err := parser.ParseStateFile(stateFile)
	if err != nil {
		fmt.Printf("Error parsing state file: %v\n", err)
		os.Exit(1)
	}

	var arm64CompatibleCount int
	var migrateableCount int
	var totalAnalyzedCount int
	var analyses []analyzer.ARM64Analysis

	for _, resource := range state.Resources {
		if resource.Mode != "managed" {
			continue
		}

		analysis := analyzer.AnalyzeResource(resource)

		if analysis.Supported {
			totalAnalyzedCount++
			analyses = append(analyses, analysis)

			if analysis.ARM64Compatible {
				arm64CompatibleCount++
				// Check if resource is ARM64-compatible but not currently using ARM64
				if canMigrateToARM64(analysis) {
					migrateableCount++
				}
			}
		}
	}

	if format == "json" {
		output := JSONOutput{
			Resources: analyses,
		}
		output.Summary.TotalAnalyzed = totalAnalyzedCount
		output.Summary.ARM64Compatible = arm64CompatibleCount
		output.Summary.Migrateable = migrateableCount
		if totalAnalyzedCount > 0 {
			output.Summary.CompatibilityRate = float64(arm64CompatibleCount) / float64(totalAnalyzedCount) * 100
		}
		output.Summary.MigrateablePercent = calculateMigrateablePercent(migrateableCount, arm64CompatibleCount)

		jsonData, err := json.Marshal(output)
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(jsonData))
	} else {
		fmt.Println("tf-arm: Terraform State ARM64 Analyzer")
		fmt.Printf("Analyzing Terraform state file: %s\n", stateFile)
		fmt.Println("")

		rep := reporter.New()
		rep.PrintHeader(len(state.Resources))

		for _, analysis := range analyses {
			rep.PrintAnalysis(analysis)
		}

		rep.PrintSummary(totalAnalyzedCount, arm64CompatibleCount, migrateableCount)
	}

	if exitCode != 0 && migrateableCount > 0 {
		os.Exit(exitCode)
	}
}
