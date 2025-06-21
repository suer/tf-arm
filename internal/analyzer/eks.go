package analyzer

import "github.com/suer/tf-arm/internal/parser"

type EKSAnalyzer struct{}

func (a *EKSAnalyzer) SupportedType() string {
	return "aws_eks_node_group"
}

func (a *EKSAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: true,
		CurrentArch:     "X86_64",
	}

	for _, instance := range resource.Instances {
		if instanceTypes, exists := instance.Attributes["instance_types"]; exists {
			instanceTypesList := instanceTypes.([]any)
			if len(instanceTypesList) > 0 {
				instanceType := instanceTypesList[0].(string)
				
				if isARM64InstanceType(instanceType) {
					analysis.CurrentArch = "ARM64"
					analysis.RecommendedArch = "ARM64"
					analysis.Notes = "Already using ARM64 instance type"
				} else if hasARM64Alternative(instanceType) {
					analysis.RecommendedArch = getARM64Alternative(instanceType)
					analysis.Notes = "Can migrate to ARM64 instance type: " + analysis.RecommendedArch
				} else {
					analysis.Notes = "Can use ARM64 instance types for EKS node group"
				}
			}
		} else {
			analysis.Notes = "Can specify ARM64 instance types for EKS node group"
		}
		
		// Check AMI type
		if amiType, exists := instance.Attributes["ami_type"]; exists {
			amiTypeStr := amiType.(string)
			if amiTypeStr == "AL2_ARM_64" {
				analysis.CurrentArch = "ARM64"
				analysis.Notes = "Already using ARM64 AMI type"
			} else {
				analysis.Notes += " | Consider using AL2_ARM_64 AMI type"
			}
		}
	}
	return analysis
}

type FargateAnalyzer struct{}

func (a *FargateAnalyzer) SupportedType() string {
	return "aws_ecs_service"
}

func (a *FargateAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: true,
		CurrentArch:     "X86_64 (default)",
		RecommendedArch: "ARM64",
	}

	for _, instance := range resource.Instances {
		if launchType, exists := instance.Attributes["launch_type"]; exists {
			if launchType == "FARGATE" {
				analysis.Notes = "Fargate supports ARM64. Check task definition cpu_architecture"
			} else {
				analysis.Notes = "Not using Fargate launch type"
				analysis.ARM64Compatible = false
			}
		} else {
			// Check if using capacity provider strategy for Fargate
			if capacityProviderStrategy, exists := instance.Attributes["capacity_provider_strategy"]; exists {
				strategies := capacityProviderStrategy.([]any)
				for _, strategy := range strategies {
					strategyMap := strategy.(map[string]any)
					if provider, exists := strategyMap["capacity_provider"]; exists {
						if provider == "FARGATE" || provider == "FARGATE_SPOT" {
							analysis.Notes = "Using Fargate capacity provider. Check task definition cpu_architecture"
							return analysis
						}
					}
				}
			}
			analysis.Notes = "Service configuration unclear for ARM64 compatibility"
		}
	}
	return analysis
}