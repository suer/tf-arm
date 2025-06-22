package analyzer

import "github.com/suer/tf-arm/internal/parser"

type ElastiCacheAnalyzer struct{}

func (a *ElastiCacheAnalyzer) SupportedType() string {
	return "aws_elasticache_cluster"
}

func (a *ElastiCacheAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: false,
		CurrentArch:     "X86_64",
	}

	for _, instance := range resource.Instances {
		if nodeType, exists := instance.Attributes["node_type"]; exists {
			nodeTypeStr, ok := nodeType.(string)
			if !ok {
				continue
			}

			if isARM64ElastiCacheNodeType(nodeTypeStr) {
				analysis.ARM64Compatible = true
				analysis.CurrentArch = "ARM64"
				analysis.RecommendedArch = "ARM64"
				analysis.AlreadyUsingARM64 = true
				analysis.Notes = "Already using ARM64 node type"
			} else if hasARM64ElastiCacheAlternative(nodeTypeStr) {
				analysis.ARM64Compatible = true
				analysis.RecommendedArch = getARM64ElastiCacheAlternative(nodeTypeStr)
				analysis.Notes = "Can migrate to ARM64 node type: " + analysis.RecommendedArch
			} else {
				analysis.Notes = "No ARM64 compatible node type available"
			}
		}
	}
	return analysis
}

type MemoryDBAnalyzer struct{}

func (a *MemoryDBAnalyzer) SupportedType() string {
	return "aws_memorydb_cluster"
}

func (a *MemoryDBAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: false,
		CurrentArch:     "X86_64",
	}

	for _, instance := range resource.Instances {
		if nodeType, exists := instance.Attributes["node_type"]; exists {
			nodeTypeStr, ok := nodeType.(string)
			if !ok {
				continue
			}

			if isARM64MemoryDBNodeType(nodeTypeStr) {
				analysis.ARM64Compatible = true
				analysis.CurrentArch = "ARM64"
				analysis.RecommendedArch = "ARM64"
				analysis.AlreadyUsingARM64 = true
				analysis.Notes = "Already using ARM64 node type"
			} else if hasARM64MemoryDBAlternative(nodeTypeStr) {
				analysis.ARM64Compatible = true
				analysis.RecommendedArch = getARM64MemoryDBAlternative(nodeTypeStr)
				analysis.Notes = "Can migrate to ARM64 node type: " + analysis.RecommendedArch
			} else {
				analysis.Notes = "No ARM64 compatible node type available"
			}
		}
	}
	return analysis
}

func isARM64ElastiCacheNodeType(nodeType string) bool {
	arm64NodeTypes := []string{
		"cache.r6g.large", "cache.r6g.xlarge", "cache.r6g.2xlarge", "cache.r6g.4xlarge",
		"cache.r6g.8xlarge", "cache.r6g.12xlarge", "cache.r6g.16xlarge",
		"cache.r6gd.large", "cache.r6gd.xlarge", "cache.r6gd.2xlarge", "cache.r6gd.4xlarge",
		"cache.r6gd.8xlarge", "cache.r6gd.12xlarge", "cache.r6gd.16xlarge",
		"cache.t4g.nano", "cache.t4g.micro", "cache.t4g.small", "cache.t4g.medium",
	}

	for _, armType := range arm64NodeTypes {
		if nodeType == armType {
			return true
		}
	}
	return false
}

func hasARM64ElastiCacheAlternative(nodeType string) bool {
	_, exists := getElastiCacheX86ToArm64Map()[nodeType]
	return exists
}

func getARM64ElastiCacheAlternative(nodeType string) string {
	return getElastiCacheX86ToArm64Map()[nodeType]
}

func getElastiCacheX86ToArm64Map() map[string]string {
	return map[string]string{
		"cache.t3.nano":     "cache.t4g.nano",
		"cache.t3.micro":    "cache.t4g.micro",
		"cache.t3.small":    "cache.t4g.small",
		"cache.t3.medium":   "cache.t4g.medium",
		"cache.r5.large":    "cache.r6g.large",
		"cache.r5.xlarge":   "cache.r6g.xlarge",
		"cache.r5.2xlarge":  "cache.r6g.2xlarge",
		"cache.r5.4xlarge":  "cache.r6g.4xlarge",
		"cache.r5.8xlarge":  "cache.r6g.8xlarge",
		"cache.r5.12xlarge": "cache.r6g.12xlarge",
		"cache.r5.16xlarge": "cache.r6g.16xlarge",
	}
}

func isARM64MemoryDBNodeType(nodeType string) bool {
	arm64NodeTypes := []string{
		"db.r6g.large", "db.r6g.xlarge", "db.r6g.2xlarge", "db.r6g.4xlarge",
		"db.r6g.8xlarge", "db.r6g.12xlarge", "db.r6g.16xlarge",
		"db.r6gd.large", "db.r6gd.xlarge", "db.r6gd.2xlarge", "db.r6gd.4xlarge",
		"db.r6gd.8xlarge", "db.r6gd.12xlarge", "db.r6gd.16xlarge",
		"db.t4g.small", "db.t4g.medium",
	}

	for _, armType := range arm64NodeTypes {
		if nodeType == armType {
			return true
		}
	}
	return false
}

func hasARM64MemoryDBAlternative(nodeType string) bool {
	_, exists := getMemoryDBX86ToArm64Map()[nodeType]
	return exists
}

func getARM64MemoryDBAlternative(nodeType string) string {
	return getMemoryDBX86ToArm64Map()[nodeType]
}

func getMemoryDBX86ToArm64Map() map[string]string {
	return map[string]string{
		"db.t3.small":    "db.t4g.small",
		"db.t3.medium":   "db.t4g.medium",
		"db.r5.large":    "db.r6g.large",
		"db.r5.xlarge":   "db.r6g.xlarge",
		"db.r5.2xlarge":  "db.r6g.2xlarge",
		"db.r5.4xlarge":  "db.r6g.4xlarge",
		"db.r5.8xlarge":  "db.r6g.8xlarge",
		"db.r5.12xlarge": "db.r6g.12xlarge",
		"db.r5.16xlarge": "db.r6g.16xlarge",
	}
}
