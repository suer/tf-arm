package analyzer

import "github.com/suer/tf-arm/internal/parser"

type RDSAnalyzer struct{}

func (a *RDSAnalyzer) SupportedType() string {
	return "aws_db_instance"
}

func (a *RDSAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: false,
		CurrentArch:     "X86_64",
	}

	for _, instance := range resource.Instances {
		if instanceClass, exists := instance.Attributes["instance_class"]; exists {
			instanceClassStr := instanceClass.(string)
			
			if isARM64RDSInstanceClass(instanceClassStr) {
				analysis.ARM64Compatible = true
				analysis.CurrentArch = "ARM64"
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "Already using ARM64 instance class"
			} else if hasARM64RDSAlternative(instanceClassStr) {
				analysis.ARM64Compatible = true
				analysis.RecommendedArch = getARM64RDSAlternative(instanceClassStr)
				analysis.Notes = "Can migrate to ARM64 instance class: " + analysis.RecommendedArch
			} else {
				analysis.Notes = "No ARM64 compatible instance class available"
			}
		}
	}
	return analysis
}

type AuroraAnalyzer struct{}

func (a *AuroraAnalyzer) SupportedType() string {
	return "aws_rds_cluster"
}

func (a *AuroraAnalyzer) Analyze(resource parser.TerraformResource) ARM64Analysis {
	analysis := ARM64Analysis{
		ResourceType:    resource.Type,
		ResourceName:    resource.Name,
		ARM64Compatible: false,
		CurrentArch:     "X86_64",
	}

	for _, instance := range resource.Instances {
		if engine, exists := instance.Attributes["engine"]; exists {
			engineStr := engine.(string)
			
			// Aurora supports ARM64 for MySQL and PostgreSQL
			if engineStr == "aurora-mysql" || engineStr == "aurora-postgresql" {
				analysis.ARM64Compatible = true
				analysis.RecommendedArch = "ARM64"
				analysis.Notes = "Aurora " + engineStr + " supports ARM64 with compatible instance classes"
			} else {
				analysis.Notes = "Engine " + engineStr + " may not support ARM64"
			}
		}
	}
	return analysis
}

func isARM64RDSInstanceClass(instanceClass string) bool {
	arm64Classes := []string{
		"db.t4g.nano", "db.t4g.micro", "db.t4g.small", "db.t4g.medium", "db.t4g.large", "db.t4g.xlarge", "db.t4g.2xlarge",
		"db.r6g.large", "db.r6g.xlarge", "db.r6g.2xlarge", "db.r6g.4xlarge", "db.r6g.8xlarge", "db.r6g.12xlarge", "db.r6g.16xlarge",
		"db.r6gd.large", "db.r6gd.xlarge", "db.r6gd.2xlarge", "db.r6gd.4xlarge", "db.r6gd.8xlarge", "db.r6gd.12xlarge", "db.r6gd.16xlarge",
	}
	
	for _, armClass := range arm64Classes {
		if instanceClass == armClass {
			return true
		}
	}
	return false
}

func hasARM64RDSAlternative(instanceClass string) bool {
	_, exists := getRDSX86ToArm64Map()[instanceClass]
	return exists
}

func getARM64RDSAlternative(instanceClass string) string {
	return getRDSX86ToArm64Map()[instanceClass]
}

func getRDSX86ToArm64Map() map[string]string {
	return map[string]string{
		"db.t3.nano":     "db.t4g.nano",
		"db.t3.micro":    "db.t4g.micro",
		"db.t3.small":    "db.t4g.small",
		"db.t3.medium":   "db.t4g.medium",
		"db.t3.large":    "db.t4g.large",
		"db.t3.xlarge":   "db.t4g.xlarge",
		"db.t3.2xlarge":  "db.t4g.2xlarge",
		"db.r5.large":    "db.r6g.large",
		"db.r5.xlarge":   "db.r6g.xlarge",
		"db.r5.2xlarge":  "db.r6g.2xlarge",
		"db.r5.4xlarge":  "db.r6g.4xlarge",
		"db.r5.8xlarge":  "db.r6g.8xlarge",
		"db.r5.12xlarge": "db.r6g.12xlarge",
		"db.r5.16xlarge": "db.r6g.16xlarge",
	}
}