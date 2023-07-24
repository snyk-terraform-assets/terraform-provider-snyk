variable "snyk_token" {
  type      = string
  sensitive = true
}

provider "snyk" {
  # example configuration here
  api_token = var.snyk_token
  endpoint  = "https://api.snyk.io/rest"
}
