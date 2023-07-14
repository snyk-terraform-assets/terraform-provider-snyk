variable "snyk_token" {
  type      = string
  sensitive = true
}

terraform {
  required_providers {
    snyk = {
      source = "registry.terraform.io/snyk-terraform-assets/snyk"
    }
  }
}

provider "snyk" {
  api_token = var.snyk_token
  endpoint  = "https://api.snyk.io/rest"
}

resource "snyk_environment" "example" {
  name            = "aws 12345"
  kind            = "aws"
  organization_id = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXXXXX"
  aws {
    role_arn = "arn:aws:iam::XXXXXXXXXXXX:role/snyk-cloud-role-XXXXXXXX"
  }
  #  azure {
  #    application_id = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  #    subscription_id = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  #    tenant_id = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  #  }
  #  google {
  #    project_id = "XXX"
  #    service_account_email = "XXX@XXX.iam.gserviceaccount.com"
  #  }

}
