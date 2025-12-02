terraform {
  required_providers {
    thetalake = {
      source = "radugheorghies/thetalake"
    }
  }
}

provider "thetalake" {
  endpoint = "https://api.thetalake.ai/api/v1" # Replace with actual if different
  token    = "your-api-token"
}

# Resources

resource "thetalake_case" "test" {
  name        = "Test Case"
  number      = "CASE-123"
  visibility  = "PRIVATE"
  description = "A test case created via Terraform"
}

resource "thetalake_user" "test" {
  name     = "Test User"
  email    = "test_user@example.com"
  password = "SecurePassword123!"
  role_id  = 1
}

resource "thetalake_directory_group" "test" {
  name        = "Test Group"
  description = "A test directory group"
  external_id = "EXT-GROUP-001"
}

resource "thetalake_tag" "test" {
  name        = "Test Tag"
  description = "A test tag"
}

resource "thetalake_retention_policy" "test" {
  name                  = "Test Policy"
  description           = "A test retention policy"
  retention_period_days = 365
}

resource "thetalake_legal_hold" "test" {
  name        = "Test Hold"
  description = "A test legal hold"
  case_id     = 123 # Replace with valid Case ID if needed
}

resource "thetalake_export" "test" {
  name        = "Test Export"
  description = "A test export"
  query_id    = 456 # Replace with valid Query ID
  format      = "CSV"
}

resource "thetalake_integration_state" "test" {
  integration_id = "789" # Replace with valid Integration ID
  paused         = false
}

resource "thetalake_record" "test" {
  id           = "rec_12345" # Replace with valid Record ID
  review_state = "reviewed"
  comment      = "Reviewed via Terraform"
}

resource "thetalake_case_record" "test" {
  case_id   = thetalake_case.test.id
  record_id = "rec_12345" # Replace with valid Record ID
}

# Data Sources

data "thetalake_case" "example" {
  id = thetalake_case.test.id
}

data "thetalake_user" "example" {
  id = thetalake_user.test.id
}

data "thetalake_directory_group" "example" {
  id = thetalake_directory_group.test.id
}

data "thetalake_tag" "example" {
  id = thetalake_tag.test.id
}

data "thetalake_retention_policy" "example" {
  id = thetalake_retention_policy.test.id
}

data "thetalake_legal_hold" "example" {
  id = thetalake_legal_hold.test.id
}

data "thetalake_export" "example" {
  id = thetalake_export.test.id
}

data "thetalake_record" "example" {
  id = "rec_12345" # Replace with valid Record ID
}

data "thetalake_integration_state" "example" {
  integration_id = "789" # Replace with valid Integration ID
}

data "thetalake_analysis" "example" {
  id = "707" # Replace with valid Analysis ID
}

data "thetalake_audit_logs" "example" {}

data "thetalake_events" "example" {}

data "thetalake_system_status" "example" {}

data "thetalake_analysis_policies" "example" {}

data "thetalake_analysis_policy_hits" "example" {}
