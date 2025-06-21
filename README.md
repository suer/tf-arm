# tf-arm: Terraform State ARM64 Analyzer

A command-line tool that analyzes Terraform state files to identify AWS resources that can be migrated to ARM64 architecture for cost optimization.

## Installation

### Build from Source

```bash
git clone https://github.com/suer/tf-arm
cd tf-arm
go build -o tf-arm ./cmd
```

## Usage

### Basic Usage

```bash
./tf-arm <terraform-state-file>
```

### Example

```bash
./tf-arm terraform.tfstate
```
