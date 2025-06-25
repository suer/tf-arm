package analyzer

import "github.com/suer/tf-arm/internal/parser"

type EMRAnalyzer struct{}

func (a *EMRAnalyzer) SupportedType() string {
	return "aws_emr_cluster"
}

func (a *EMRAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: false,
		CurrentArch:     "X86_64",
	}

	for _, instance := range resource.Instances {
		if instanceGroups, exists := instance.Attributes["master_instance_group"]; exists {
			instanceGroupList, ok := instanceGroups.([]any)
			if !ok {
				continue
			}
			if len(instanceGroupList) > 0 {
				instanceGroup, ok := instanceGroupList[0].(map[string]any)
				if !ok {
					continue
				}
				if instanceType, exists := instanceGroup["instance_type"]; exists {
					instanceTypeStr, ok := instanceType.(string)
					if !ok {
						continue
					}

					if isARM64EMRInstanceType(instanceTypeStr) {
						analysis.ARM64Compatible = true
						analysis.CurrentArch = "ARM64"
						analysis.RecommendedArch = "ARM64"
						analysis.AlreadyUsingARM64 = true
						analysis.Notes = "Already using ARM64 instance type for master node"
					} else if hasARM64EMRAlternative(instanceTypeStr) {
						analysis.ARM64Compatible = true
						analysis.RecommendedArch = getARM64EMRAlternative(instanceTypeStr)
						analysis.Notes = "Can migrate to ARM64 instance type: " + analysis.RecommendedArch
					} else {
						analysis.Notes = "EMR supports ARM64 with Graviton2 instances"
					}
				}
			}
		}
	}
	return analysis
}

type EMRServerlessAnalyzer struct{}

func (a *EMRServerlessAnalyzer) SupportedType() string {
	return "aws_emrserverless_application"
}

func (a *EMRServerlessAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: true,
		CurrentArch:     "X86_64 (default)",
		RecommendedArch: "ARM64",
		Notes:           "EMR Serverless supports ARM64 architecture for cost optimization",
	}

	for _, instance := range resource.Instances {
		if architecture, exists := instance.Attributes["architecture"]; exists {
			if architecture == "ARM64" {
				analysis.CurrentArch = "ARM64"
				analysis.AlreadyUsingARM64 = true
				analysis.Notes = "Already using ARM64 architecture"
			}
		}
	}
	return analysis
}

func isARM64EMRInstanceType(instanceType string) bool {
	// EMR supports ARM64 with Graviton2 and Graviton3 instances
	arm64Types := []string{
		// Graviton2
		"m6g.xlarge", "m6g.2xlarge", "m6g.4xlarge", "m6g.8xlarge", "m6g.12xlarge", "m6g.16xlarge",
		"m6gd.xlarge", "m6gd.2xlarge", "m6gd.4xlarge", "m6gd.8xlarge", "m6gd.12xlarge", "m6gd.16xlarge",
		"c6g.xlarge", "c6g.2xlarge", "c6g.4xlarge", "c6g.8xlarge", "c6g.12xlarge", "c6g.16xlarge",
		"c6gd.xlarge", "c6gd.2xlarge", "c6gd.4xlarge", "c6gd.8xlarge", "c6gd.12xlarge", "c6gd.16xlarge",
		"r6g.xlarge", "r6g.2xlarge", "r6g.4xlarge", "r6g.8xlarge", "r6g.12xlarge", "r6g.16xlarge",
		"r6gd.xlarge", "r6gd.2xlarge", "r6gd.4xlarge", "r6gd.8xlarge", "r6gd.12xlarge", "r6gd.16xlarge",
		// Graviton3
		"c7g.xlarge", "c7g.2xlarge", "c7g.4xlarge", "c7g.8xlarge", "c7g.12xlarge", "c7g.16xlarge",
		"m7g.xlarge", "m7g.2xlarge", "m7g.4xlarge", "m7g.8xlarge", "m7g.12xlarge", "m7g.16xlarge",
		"r7g.xlarge", "r7g.2xlarge", "r7g.4xlarge", "r7g.8xlarge", "r7g.12xlarge", "r7g.16xlarge",
	}

	for _, armType := range arm64Types {
		if instanceType == armType {
			return true
		}
	}
	return false
}

func hasARM64EMRAlternative(instanceType string) bool {
	_, exists := getEMRX86ToArm64Map()[instanceType]
	return exists
}

func getARM64EMRAlternative(instanceType string) string {
	return getEMRX86ToArm64Map()[instanceType]
}

func getEMRX86ToArm64Map() map[string]string {
	return map[string]string{
		// M5 -> M7g (Graviton3 - better performance than M6g)
		"m5.xlarge":   "m7g.xlarge",
		"m5.2xlarge":  "m7g.2xlarge",
		"m5.4xlarge":  "m7g.4xlarge",
		"m5.8xlarge":  "m7g.8xlarge",
		"m5.12xlarge": "m7g.12xlarge",
		"m5.16xlarge": "m7g.16xlarge",
		// C5 -> C7g (Graviton3 - better performance than C6g)
		"c5.xlarge":   "c7g.xlarge",
		"c5.2xlarge":  "c7g.2xlarge",
		"c5.4xlarge":  "c7g.4xlarge",
		"c5.9xlarge":  "c7g.8xlarge",
		"c5.12xlarge": "c7g.12xlarge",
		"c5.18xlarge": "c7g.16xlarge",
		// R5 -> R7g (Graviton3 - better performance than R6g)
		"r5.xlarge":   "r7g.xlarge",
		"r5.2xlarge":  "r7g.2xlarge",
		"r5.4xlarge":  "r7g.4xlarge",
		"r5.8xlarge":  "r7g.8xlarge",
		"r5.12xlarge": "r7g.12xlarge",
		"r5.16xlarge": "r7g.16xlarge",
	}
}
