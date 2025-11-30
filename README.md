# Theta Lake Terraform Provider

This is a Terraform provider for managing resources in [Theta Lake](https://thetalake.com/). It allows you to manage Cases, Users, Directory Groups, Retention Policies, Legal Holds, Tags, Exports, Records, and more using Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20 (to build the provider plugin)

## Building the Provider

1. Clone the repository:
   ```bash
   git clone https://github.com/radugheorghies/thetalake-terraform-provider.git
   cd thetalake-terraform-provider
   ```

2. Build the provider:
   ```bash
   go build -o terraform-provider-thetalake
   ```

## Local Development & Installation

Since this provider is not yet published to the Terraform Registry, you need to configure Terraform to use the locally built binary.

1. Create a `dev_override.tfrc` file (or update your `~/.terraformrc`):

   ```hcl
   provider_installation {
     dev_overrides {
       "radugheorghies/thetalake" = "/path/to/your/go/bin" # Or the directory where you built the binary
     }
     # For all other providers, install them directly from their origin provider
     # registries as normal. If you omit this, Terraform will _only_ use
     # the dev_overrides block, and you will be unable to use other providers.
     direct {}
   }
   ```
   *Note: Replace `/path/to/your/go/bin` with the absolute path to the directory containing the `terraform-provider-thetalake` binary.*

2. Tell Terraform to use this configuration file:
   ```bash
   export TF_CLI_CONFIG_FILE=$(pwd)/dev_override.tfrc
   ```

## Usage

### Provider Configuration

```hcl
terraform {
  required_providers {
    thetalake = {
      source  = "radugheorghies/thetalake"
      version = "~> 1.0.0" # Version is ignored when using dev_overrides
    }
  }
}

provider "thetalake" {
  token    = "YOUR_API_TOKEN"
  endpoint = "https://api.thetalake.ai/api/v1"
}
```

### Resource Examples

**Manage a Case**
```hcl
resource "thetalake_case" "example" {
  name        = "Investigation 2023"
  number      = "CASE-2023-001"
  visibility  = "PRIVATE"
  description = "Internal investigation case"
  status      = "OPEN"
}
```

**Manage a User**
```hcl
resource "thetalake_user" "analyst" {
  name                  = "Jane Doe"
  email                 = "jane.doe@example.com"
  password              = "SecureP@ssw0rd!"
  password_confirmation = "SecureP@ssw0rd!"
  role_id               = 1
}
```

**Manage a Retention Policy**
```hcl
resource "thetalake_retention_policy" "seven_years" {
  name                  = "7 Year Retention"
  description           = "Retain data for 7 years"
  retention_period_days = 2555
}
```

**Manage an Export**
```hcl
resource "thetalake_export" "monthly_report" {
  name   = "Monthly Compliance Report"
  format = "csv"
}
```

### Data Source Examples

**Read Audit Logs**
```hcl
data "thetalake_audit_logs" "recent" {}

output "logs" {
  value = data.thetalake_audit_logs.recent.logs
}
```

**Check System Status**
```hcl
data "thetalake_system_status" "current" {}

output "status" {
  value = data.thetalake_system_status.current.status
}
```

## Available Resources & Data Sources

**Resources:**
* `thetalake_case`
* `thetalake_case_record`
* `thetalake_user`
* `thetalake_directory_group`
* `thetalake_retention_policy`
* `thetalake_legal_hold`
* `thetalake_tag`
* `thetalake_integration_state`
* `thetalake_export`
* `thetalake_record`

**Data Sources:**
* `thetalake_audit_logs`
* `thetalake_events`
* `thetalake_analysis_policies`
* `thetalake_analysis_policy_hits`
* `thetalake_system_status`

## Testing

To run acceptance tests (requires a valid API token):

```bash
export THETALAKE_TOKEN="your-token"
export THETALAKE_ENDPOINT="https://api.thetalake.ai/api/v1"
TF_ACC=1 go test -v ./...
```
