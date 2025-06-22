package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

type TerraformState struct {
	Version   int                 `json:"version"`
	Resources []TerraformResource `json:"resources"`
}

type TerraformResource struct {
	Mode      string             `json:"mode"`
	Type      string             `json:"type"`
	Name      string             `json:"name"`
	Module    string             `json:"module,omitempty"`
	Provider  string             `json:"provider"`
	Instances []ResourceInstance `json:"instances"`
}

type ResourceInstance struct {
	Attributes map[string]interface{} `json:"attributes"`
}

// GetFullAddress returns the full Terraform address for the resource
func (r *TerraformResource) GetFullAddress() string {
	if r.Module != "" {
		return fmt.Sprintf("%s.%s.%s", r.Module, r.Type, r.Name)
	}
	return fmt.Sprintf("%s.%s", r.Type, r.Name)
}

func ParseStateFile(filename string) (*TerraformState, error) {
	if filename == "" {
		return nil, fmt.Errorf("filename cannot be empty")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("state file is empty")
	}

	var state TerraformState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Basic validation of parsed state
	if state.Version == 0 {
		return nil, fmt.Errorf("invalid state file: version is missing or zero")
	}

	return &state, nil
}
