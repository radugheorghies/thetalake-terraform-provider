# Theta Lake Terraform Provider Examples

This directory contains examples of how to use the Theta Lake Terraform Provider.

## Prerequisites

- Terraform 0.13+
- A Theta Lake API Endpoint and Token

## Usage

1.  Navigate to this directory:
    ```bash
    cd examples
    ```

2.  Initialize Terraform:
    ```bash
    terraform init
    ```

3.  Create a `terraform.tfvars` file or set environment variables for your credentials:
    ```bash
    export TF_VAR_endpoint="https://api.thetalake.ai/api/v1"
    export TF_VAR_token="your-api-token"
    ```

4.  Run `terraform plan` to see the changes:
    ```bash
    terraform plan
    ```

5.  Run `terraform apply` to apply the changes:
    ```bash
    terraform apply
    ```

## Resources

See `main.tf` for examples of all available resources and data sources.
