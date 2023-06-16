variable "SNYK_TOKEN" {
  default = ""
}

provider "snyk" {
  # example configuration here
  api_token = var.SNYK_TOKEN
  endpoint  = "https://api.snyk.io/rest"
}