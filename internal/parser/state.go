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
	Provider  string             `json:"provider"`
	Instances []ResourceInstance `json:"instances"`
}

type ResourceInstance struct {
	Attributes map[string]any `json:"attributes"`
}

func ParseStateFile(filename string) (*TerraformState, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state TerraformState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &state, nil
}
