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
