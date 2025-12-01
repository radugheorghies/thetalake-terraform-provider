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

resource "thetalake_case" "test" {
  name       = "Test Case"
  number     = "CASE-123"
  visibility = "PRIVATE"
}

resource "thetalake_user" "test" {
  name  = "Test User"
  email = "test@example.com"
}

resource "thetalake_directory_group" "test" {
  name        = "Test Group"
  description = "A test directory group"
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
  case_id     = 123
}

resource "thetalake_export" "test" {
  name        = "Test Export"
  description = "A test export"
  query_id    = 456
  format      = "CSV"
}

resource "thetalake_integration_state" "test" {
  integration_id = "789"
  paused         = false
}

data "thetalake_case" "example" {
  id = "123"
}

data "thetalake_user" "example" {
  id = "456"
}

data "thetalake_directory_group" "example" {
  id = "789"
}

data "thetalake_tag" "example" {
  id = "101"
}

data "thetalake_retention_policy" "example" {
  id = "202"
}

data "thetalake_legal_hold" "example" {
  id = "303"
}

data "thetalake_export" "example" {
  id = "404"
}

data "thetalake_record" "example" {
  id = "505"
}

data "thetalake_integration_state" "example" {
  integration_id = "606"
}

data "thetalake_analysis" "example" {
  id = "707"
}
