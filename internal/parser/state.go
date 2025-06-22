package parser

import (
	"encoding/json"
	"fmt"
	"io"
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

	fileInfo, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to stat state file: %w", err)
	}

	if fileInfo.Size() == 0 {
		return nil, fmt.Errorf("state file is empty")
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open state file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	
	var state TerraformState
	if err := decoder.Decode(&state); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("state file is empty or invalid JSON")
		}
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Basic validation of parsed state
	if state.Version == 0 {
		return nil, fmt.Errorf("invalid state file: version is missing or zero")
	}

	return &state, nil
}
