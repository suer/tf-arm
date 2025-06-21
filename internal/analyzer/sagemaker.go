package analyzer

import "github.com/suer/tf-arm/internal/parser"

type SageMakerAnalyzer struct{}

func (a *SageMakerAnalyzer) SupportedType() string {
	return "aws_sagemaker_endpoint_configuration"
}

func (a *SageMakerAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: false,
		CurrentArch:     "X86_64",
	}

	for _, instance := range resource.Instances {
		if productionVariants, exists := instance.Attributes["production_variants"]; exists {
			variantsList := productionVariants.([]any)
			if len(variantsList) > 0 {
				variant := variantsList[0].(map[string]any)

				if instanceType, exists := variant["instance_type"]; exists {
					instanceTypeStr := instanceType.(string)

					if isARM64SageMakerInstanceType(instanceTypeStr) {
						analysis.ARM64Compatible = true
						analysis.CurrentArch = "ARM64"
						analysis.RecommendedArch = "ARM64"
						analysis.AlreadyUsingARM64 = true
						analysis.Notes = "Already using ARM64 instance type"
					} else if hasARM64SageMakerAlternative(instanceTypeStr) {
						analysis.ARM64Compatible = true
						analysis.RecommendedArch = getARM64SageMakerAlternative(instanceTypeStr)
						analysis.Notes = "Can migrate to ARM64 instance type: " + analysis.RecommendedArch
					} else {
						analysis.Notes = "No ARM64 compatible instance type available"
					}
				}
			}
		}
	}
	return analysis
}

type GameLiftAnalyzer struct{}

func (a *GameLiftAnalyzer) SupportedType() string {
	return "aws_gamelift_fleet"
}

func (a *GameLiftAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: false,
		CurrentArch:     "X86_64",
	}

	for _, instance := range resource.Instances {
		if ec2InstanceType, exists := instance.Attributes["ec2_instance_type"]; exists {
			instanceTypeStr := ec2InstanceType.(string)

			if isARM64GameLiftInstanceType(instanceTypeStr) {
				analysis.ARM64Compatible = true
				analysis.CurrentArch = "ARM64"
				analysis.RecommendedArch = "ARM64"
				analysis.AlreadyUsingARM64 = true
				analysis.Notes = "Already using ARM64 instance type"
			} else if hasARM64GameLiftAlternative(instanceTypeStr) {
				analysis.ARM64Compatible = true
				analysis.RecommendedArch = getARM64GameLiftAlternative(instanceTypeStr)
				analysis.Notes = "Can migrate to ARM64 instance type: " + analysis.RecommendedArch
			} else {
				analysis.Notes = "GameLift supports ARM64 with Graviton2 instances"
			}
		}
	}
	return analysis
}

func isARM64SageMakerInstanceType(instanceType string) bool {
	arm64Types := []string{
		"ml.m6g.large", "ml.m6g.xlarge", "ml.m6g.2xlarge", "ml.m6g.4xlarge",
		"ml.m6g.8xlarge", "ml.m6g.12xlarge", "ml.m6g.16xlarge",
		"ml.m6gd.large", "ml.m6gd.xlarge", "ml.m6gd.2xlarge", "ml.m6gd.4xlarge",
		"ml.m6gd.8xlarge", "ml.m6gd.12xlarge", "ml.m6gd.16xlarge",
		"ml.c6g.large", "ml.c6g.xlarge", "ml.c6g.2xlarge", "ml.c6g.4xlarge",
		"ml.c6g.8xlarge", "ml.c6g.12xlarge", "ml.c6g.16xlarge",
		"ml.c6gd.large", "ml.c6gd.xlarge", "ml.c6gd.2xlarge", "ml.c6gd.4xlarge",
		"ml.c6gd.8xlarge", "ml.c6gd.12xlarge", "ml.c6gd.16xlarge",
		"ml.r6g.large", "ml.r6g.xlarge", "ml.r6g.2xlarge", "ml.r6g.4xlarge",
		"ml.r6g.8xlarge", "ml.r6g.12xlarge", "ml.r6g.16xlarge",
		"ml.r6gd.large", "ml.r6gd.xlarge", "ml.r6gd.2xlarge", "ml.r6gd.4xlarge",
		"ml.r6gd.8xlarge", "ml.r6gd.12xlarge", "ml.r6gd.16xlarge",
	}

	for _, armType := range arm64Types {
		if instanceType == armType {
			return true
		}
	}
	return false
}

func hasARM64SageMakerAlternative(instanceType string) bool {
	_, exists := getSageMakerX86ToArm64Map()[instanceType]
	return exists
}

func getARM64SageMakerAlternative(instanceType string) string {
	return getSageMakerX86ToArm64Map()[instanceType]
}

func getSageMakerX86ToArm64Map() map[string]string {
	return map[string]string{
		"ml.m5.large":    "ml.m6g.large",
		"ml.m5.xlarge":   "ml.m6g.xlarge",
		"ml.m5.2xlarge":  "ml.m6g.2xlarge",
		"ml.m5.4xlarge":  "ml.m6g.4xlarge",
		"ml.m5.8xlarge":  "ml.m6g.8xlarge",
		"ml.m5.12xlarge": "ml.m6g.12xlarge",
		"ml.m5.16xlarge": "ml.m6g.16xlarge",
		"ml.c5.large":    "ml.c6g.large",
		"ml.c5.xlarge":   "ml.c6g.xlarge",
		"ml.c5.2xlarge":  "ml.c6g.2xlarge",
		"ml.c5.4xlarge":  "ml.c6g.4xlarge",
		"ml.c5.9xlarge":  "ml.c6g.8xlarge",
		"ml.c5.18xlarge": "ml.c6g.16xlarge",
		"ml.r5.large":    "ml.r6g.large",
		"ml.r5.xlarge":   "ml.r6g.xlarge",
		"ml.r5.2xlarge":  "ml.r6g.2xlarge",
		"ml.r5.4xlarge":  "ml.r6g.4xlarge",
		"ml.r5.8xlarge":  "ml.r6g.8xlarge",
		"ml.r5.12xlarge": "ml.r6g.12xlarge",
		"ml.r5.16xlarge": "ml.r6g.16xlarge",
	}
}

func isARM64GameLiftInstanceType(instanceType string) bool {
	// GameLift supports ARM64 with Graviton2 instances
	arm64Types := []string{
		"c6g.large", "c6g.xlarge", "c6g.2xlarge", "c6g.4xlarge",
		"c6g.8xlarge", "c6g.12xlarge", "c6g.16xlarge",
		"m6g.large", "m6g.xlarge", "m6g.2xlarge", "m6g.4xlarge",
		"m6g.8xlarge", "m6g.12xlarge", "m6g.16xlarge",
		"r6g.large", "r6g.xlarge", "r6g.2xlarge", "r6g.4xlarge",
		"r6g.8xlarge", "r6g.12xlarge", "r6g.16xlarge",
	}

	for _, armType := range arm64Types {
		if instanceType == armType {
			return true
		}
	}
	return false
}

func hasARM64GameLiftAlternative(instanceType string) bool {
	_, exists := getGameLiftX86ToArm64Map()[instanceType]
	return exists
}

func getARM64GameLiftAlternative(instanceType string) string {
	return getGameLiftX86ToArm64Map()[instanceType]
}

func getGameLiftX86ToArm64Map() map[string]string {
	return map[string]string{
		"c5.large":    "c6g.large",
		"c5.xlarge":   "c6g.xlarge",
		"c5.2xlarge":  "c6g.2xlarge",
		"c5.4xlarge":  "c6g.4xlarge",
		"c5.9xlarge":  "c6g.8xlarge",
		"c5.12xlarge": "c6g.12xlarge",
		"c5.18xlarge": "c6g.16xlarge",
		"m5.large":    "m6g.large",
		"m5.xlarge":   "m6g.xlarge",
		"m5.2xlarge":  "m6g.2xlarge",
		"m5.4xlarge":  "m6g.4xlarge",
		"m5.8xlarge":  "m6g.8xlarge",
		"m5.12xlarge": "m6g.12xlarge",
		"m5.16xlarge": "m6g.16xlarge",
		"r5.large":    "r6g.large",
		"r5.xlarge":   "r6g.xlarge",
		"r5.2xlarge":  "r6g.2xlarge",
		"r5.4xlarge":  "r6g.4xlarge",
		"r5.8xlarge":  "r6g.8xlarge",
		"r5.12xlarge": "r6g.12xlarge",
		"r5.16xlarge": "r6g.16xlarge",
	}
}
