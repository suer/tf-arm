package analyzer

import "github.com/suer/tf-arm/internal/parser"

type OpenSearchAnalyzer struct{}

func (a *OpenSearchAnalyzer) SupportedType() string {
	return "aws_opensearch_domain"
}

func (a *OpenSearchAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: false,
		CurrentArch:     "X86_64",
	}

	for _, instance := range resource.Instances {
		if clusterConfig, exists := instance.Attributes["cluster_config"]; exists {
			clusterConfigList, ok := clusterConfig.([]any)
			if !ok {
				continue
			}
			if len(clusterConfigList) > 0 {
				config, ok := clusterConfigList[0].(map[string]any)
				if !ok {
					continue
				}

				if instanceType, exists := config["instance_type"]; exists {
					instanceTypeStr, ok := instanceType.(string)
					if !ok {
						continue
					}

					if isARM64OpenSearchInstanceType(instanceTypeStr) {
						analysis.ARM64Compatible = true
						analysis.CurrentArch = "ARM64"
						analysis.RecommendedArch = "ARM64"
						analysis.AlreadyUsingARM64 = true
						analysis.Notes = "Already using ARM64 instance type"
					} else if hasARM64OpenSearchAlternative(instanceTypeStr) {
						analysis.ARM64Compatible = true
						analysis.RecommendedArch = getARM64OpenSearchAlternative(instanceTypeStr)
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

type MSKAnalyzer struct{}

func (a *MSKAnalyzer) SupportedType() string {
	return "aws_msk_cluster"
}

func (a *MSKAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: false,
		CurrentArch:     "X86_64",
	}

	for _, instance := range resource.Instances {
		if brokerNodeGroupInfo, exists := instance.Attributes["broker_node_group_info"]; exists {
			brokerNodeList, ok := brokerNodeGroupInfo.([]any)
			if !ok {
				continue
			}
			if len(brokerNodeList) > 0 {
				brokerNode, ok := brokerNodeList[0].(map[string]any)
				if !ok {
					continue
				}

				if instanceType, exists := brokerNode["instance_type"]; exists {
					instanceTypeStr, ok := instanceType.(string)
					if !ok {
						continue
					}

					if isARM64MSKInstanceType(instanceTypeStr) {
						analysis.ARM64Compatible = true
						analysis.CurrentArch = "ARM64"
						analysis.RecommendedArch = "ARM64"
						analysis.AlreadyUsingARM64 = true
						analysis.Notes = "Already using ARM64 instance type"
					} else if hasARM64MSKAlternative(instanceTypeStr) {
						analysis.ARM64Compatible = true
						analysis.RecommendedArch = getARM64MSKAlternative(instanceTypeStr)
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

func isARM64OpenSearchInstanceType(instanceType string) bool {
	arm64Types := []string{
		"t4g.small.search", "t4g.medium.search",
		"m6g.large.search", "m6g.xlarge.search", "m6g.2xlarge.search", "m6g.4xlarge.search",
		"m6g.8xlarge.search", "m6g.12xlarge.search",
		"c6g.large.search", "c6g.xlarge.search", "c6g.2xlarge.search", "c6g.4xlarge.search",
		"c6g.8xlarge.search", "c6g.12xlarge.search",
		"r6g.large.search", "r6g.xlarge.search", "r6g.2xlarge.search", "r6g.4xlarge.search",
		"r6g.8xlarge.search", "r6g.12xlarge.search",
		"r6gd.large.search", "r6gd.xlarge.search", "r6gd.2xlarge.search", "r6gd.4xlarge.search",
		"r6gd.8xlarge.search", "r6gd.12xlarge.search", "r6gd.16xlarge.search",
	}

	for _, armType := range arm64Types {
		if instanceType == armType {
			return true
		}
	}
	return false
}

func hasARM64OpenSearchAlternative(instanceType string) bool {
	_, exists := getOpenSearchX86ToArm64Map()[instanceType]
	return exists
}

func getARM64OpenSearchAlternative(instanceType string) string {
	return getOpenSearchX86ToArm64Map()[instanceType]
}

func getOpenSearchX86ToArm64Map() map[string]string {
	return map[string]string{
		"t3.small.search":    "t4g.small.search",
		"t3.medium.search":   "t4g.medium.search",
		"m5.large.search":    "m6g.large.search",
		"m5.xlarge.search":   "m6g.xlarge.search",
		"m5.2xlarge.search":  "m6g.2xlarge.search",
		"m5.4xlarge.search":  "m6g.4xlarge.search",
		"m5.8xlarge.search":  "m6g.8xlarge.search",
		"m5.12xlarge.search": "m6g.12xlarge.search",
		"c5.large.search":    "c6g.large.search",
		"c5.xlarge.search":   "c6g.xlarge.search",
		"c5.2xlarge.search":  "c6g.2xlarge.search",
		"c5.4xlarge.search":  "c6g.4xlarge.search",
		"c5.9xlarge.search":  "c6g.8xlarge.search",
		"c5.18xlarge.search": "c6g.12xlarge.search",
		"r5.large.search":    "r6g.large.search",
		"r5.xlarge.search":   "r6g.xlarge.search",
		"r5.2xlarge.search":  "r6g.2xlarge.search",
		"r5.4xlarge.search":  "r6g.4xlarge.search",
		"r5.8xlarge.search":  "r6g.8xlarge.search",
		"r5.12xlarge.search": "r6g.12xlarge.search",
	}
}

func isARM64MSKInstanceType(instanceType string) bool {
	arm64Types := []string{
		"kafka.m6g.large", "kafka.m6g.xlarge", "kafka.m6g.2xlarge", "kafka.m6g.4xlarge",
		"kafka.m6g.8xlarge", "kafka.m6g.12xlarge", "kafka.m6g.16xlarge",
	}

	for _, armType := range arm64Types {
		if instanceType == armType {
			return true
		}
	}
	return false
}

func hasARM64MSKAlternative(instanceType string) bool {
	_, exists := getMSKX86ToArm64Map()[instanceType]
	return exists
}

func getARM64MSKAlternative(instanceType string) string {
	return getMSKX86ToArm64Map()[instanceType]
}

func getMSKX86ToArm64Map() map[string]string {
	return map[string]string{
		"kafka.m5.large":    "kafka.m6g.large",
		"kafka.m5.xlarge":   "kafka.m6g.xlarge",
		"kafka.m5.2xlarge":  "kafka.m6g.2xlarge",
		"kafka.m5.4xlarge":  "kafka.m6g.4xlarge",
		"kafka.m5.8xlarge":  "kafka.m6g.8xlarge",
		"kafka.m5.12xlarge": "kafka.m6g.12xlarge",
		"kafka.m5.16xlarge": "kafka.m6g.16xlarge",
	}
}
