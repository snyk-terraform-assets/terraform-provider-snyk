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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExampleResource(t *testing.T) {
	snykOrgId := readEnvironmentVar("TEST_SNYK_ORG_ID")
	awsArn := readEnvironmentVar("TEST_AWS_ARN")
	azureApplicationId := readEnvironmentVar("TEST_AZURE_APPLICATION_ID")
	azureSubscriptionId := readEnvironmentVar("TEST_AZURE_SUBSCRIPTION_ID")
	azureTenantId := readEnvironmentVar("TEST_AZURE_TENANT_ID")

	googleProjectId := readEnvironmentVar("TEST_GOOGLE_PROJECT_ID")
	googleServiceAccountEmail := readEnvironmentVar("TEST_GOOGLE_SERVICE_ACCOUNT_EMAIL")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccExampleResourceConfigForAws("initial", snykOrgId, awsArn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snyk_environment.test", "name", "initial"),
					resource.TestCheckResourceAttr("snyk_environment.test", "kind", "aws"),
					resource.TestCheckResourceAttr("snyk_environment.test", "organization_id", snykOrgId),
				),
			},
			//TODO See the issue here: https://github.com/hashicorp/terraform-plugin-framework/issues/677
			//// ImportState testing
			//{
			//	ResourceName:      "snyk_environment.test",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	// This is not normally necessary, but is here because this
			//	// example code does not have an actual upstream service.
			//	// Once the Read method is able to refresh information from
			//	// the upstream service, this can be removed.
			//	ImportStateVerifyIgnore: []string{"configurable_attribute"},
			//},
			// Update and Read testing
			{
				Config: providerConfig + testAccExampleResourceConfigForAws("updated", snykOrgId, awsArn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snyk_environment.test", "name", "updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
			// Create and Read testing
			{
				Config: providerConfig + testAccExampleResourceConfigForAzure("initial", snykOrgId, azureApplicationId, azureSubscriptionId, azureTenantId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snyk_environment.test_azure", "name", "initial"),
					resource.TestCheckResourceAttr("snyk_environment.test_azure", "kind", "azure"),
					resource.TestCheckResourceAttr("snyk_environment.test_azure", "organization_id", snykOrgId),
					resource.TestCheckResourceAttr("snyk_environment.test_azure", "azure.application_id", azureApplicationId),
					resource.TestCheckResourceAttr("snyk_environment.test_azure", "azure.subscription_id", azureSubscriptionId),
					resource.TestCheckResourceAttr("snyk_environment.test_azure", "azure.tenant_id", azureTenantId),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccExampleResourceConfigForAzure("updated", snykOrgId, azureApplicationId, azureSubscriptionId, azureTenantId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snyk_environment.test_azure", "name", "updated"),
					resource.TestCheckResourceAttr("snyk_environment.test_azure", "azure.application_id", azureApplicationId),
					resource.TestCheckResourceAttr("snyk_environment.test_azure", "azure.subscription_id", azureSubscriptionId),
					resource.TestCheckResourceAttr("snyk_environment.test_azure", "azure.tenant_id", azureTenantId),
				),
			},
			// Delete testing automatically occurs in TestCase
			// Create and Read testing
			{
				Config: providerConfig + testAccExampleResourceConfigForGoogle("initial", snykOrgId, googleProjectId, googleServiceAccountEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snyk_environment.test_google", "name", "initial"),
					resource.TestCheckResourceAttr("snyk_environment.test_google", "kind", "google"),
					resource.TestCheckResourceAttr("snyk_environment.test_google", "google.project_id", googleProjectId),
					resource.TestCheckResourceAttr("snyk_environment.test_google", "google.service_account_email", googleServiceAccountEmail),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccExampleResourceConfigForGoogle("updated", snykOrgId, googleProjectId, googleServiceAccountEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snyk_environment.test_google", "name", "updated"),
					resource.TestCheckResourceAttr("snyk_environment.test_google", "google.project_id", googleProjectId),
					resource.TestCheckResourceAttr("snyk_environment.test_google", "google.service_account_email", googleServiceAccountEmail),
				),
			},
		},
	})
}

func testAccExampleResourceConfigForAws(envName string, orgId string, awsArn string) string {
	return fmt.Sprintf(`
resource "snyk_environment" "test" {
  name = %[1]q
  kind = "aws"
  organization_id = %[2]q
  aws {
    role_arn = %[3]q
  }
}`, envName, orgId, awsArn)
}

func testAccExampleResourceConfigForAzure(envName string, orgId string, applicationId string,
	subscriptionId string, tenantId string) string {
	return fmt.Sprintf(`
resource "snyk_environment" "test_azure" {
  name = %[1]q
  kind = "azure"
  organization_id = %[2]q
  azure {
    application_id = %[3]q
    subscription_id = %[4]q
    tenant_id = %[5]q
  }
}`, envName, orgId, applicationId, subscriptionId, tenantId)
}

func testAccExampleResourceConfigForGoogle(envName string, orgId string,
	projectId string, serviceAccountEmail string) string {
	return fmt.Sprintf(`
resource "snyk_environment" "test_google" {
  name = %[1]q
  kind = "google"
  organization_id = %[2]q
  google {
    project_id = %[3]q
    service_account_email = %[4]q
  }
}`, envName, orgId, projectId, serviceAccountEmail)
}

func readEnvironmentVar(key string) string {
	res := os.Getenv(key)
	if res == "" {
		panic(fmt.Sprintf("Env variable %s not set!", key))
	}
	return res

}
