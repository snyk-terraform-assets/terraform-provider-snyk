// Copyright 2023 Snyk Limited All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"scaffolding": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccProviderConfig(t *testing.T) string {
	apiToken := readEnvVarOrFail(t, "TEST_SNYK_TOKEN")
	endpoint := os.Getenv("TEST_SNYK_API")
	if endpoint == "" {
		endpoint = "https://api.snyk.io/rest"
	}

	return fmt.Sprintf(`
terraform {
  required_providers {
    snyk = {
      source = "registry.terraform.io/snyk-terraform-assets/snyk"
    }
  }
}

provider "snyk" {
  api_token = %[1]q
  endpoint  = %[2]q
}`, apiToken, endpoint)
}

// readEnvVarOrFail reads the requested environment variable.
// If this variable is not present, the test fails.
func readEnvVarOrFail(t *testing.T, key string) string {
	val := os.Getenv(key)
	if val == "" {
		t.Fatalf("Missing environment variable %s", key)
	}
	return val
}

// readEnvVarOrSkip reads the requested environment variable.
// If this variable is not present, the test is skipped.
func readEnvVarOrSkip(t *testing.T, key string) string {
	val := os.Getenv(key)
	if val == "" {
		t.Skipf("Missing environment variable %s", key)
	}
	return val
}
