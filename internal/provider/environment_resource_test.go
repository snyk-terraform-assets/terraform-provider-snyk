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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAwsEnvironment(t *testing.T) {
	snykOrgId := readEnvVarOrFail(t, "TEST_SNYK_ORG_ID")
	awsArn := readEnvVarOrSkip(t, "TEST_AWS_ARN")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(t) + "\n" +
					testAccExampleResourceConfigForAws("initial", snykOrgId, awsArn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snyk_environment.test", "name", "initial"),
					resource.TestCheckResourceAttr("snyk_environment.test", "kind", "aws"),
					resource.TestCheckResourceAttr("snyk_environment.test", "organization_id", snykOrgId),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(t) + "\n" +
					testAccExampleResourceConfigForAws("updated", snykOrgId, awsArn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snyk_environment.test", "name", "updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAzureEnvironment(t *testing.T) {
	snykOrgId := readEnvVarOrFail(t, "TEST_SNYK_ORG_ID")
	azureApplicationId := readEnvVarOrSkip(t, "TEST_AZURE_APPLICATION_ID")
	azureSubscriptionId := readEnvVarOrSkip(t, "TEST_AZURE_SUBSCRIPTION_ID")
	azureTenantId := readEnvVarOrSkip(t, "TEST_AZURE_TENANT_ID")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(t) + "\n" +
					testAccExampleResourceConfigForAzure("initial", snykOrgId, azureApplicationId, azureSubscriptionId, azureTenantId),
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
				Config: testAccProviderConfig(t) + "\n" +
					testAccExampleResourceConfigForAzure("updated", snykOrgId, azureApplicationId, azureSubscriptionId, azureTenantId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snyk_environment.test_azure", "name", "updated"),
					resource.TestCheckResourceAttr("snyk_environment.test_azure", "azure.application_id", azureApplicationId),
					resource.TestCheckResourceAttr("snyk_environment.test_azure", "azure.subscription_id", azureSubscriptionId),
					resource.TestCheckResourceAttr("snyk_environment.test_azure", "azure.tenant_id", azureTenantId),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccGoogleEnvironment(t *testing.T) {
	snykOrgId := readEnvVarOrFail(t, "TEST_SNYK_ORG_ID")
	googleProjectId := readEnvVarOrSkip(t, "TEST_GOOGLE_PROJECT_ID")
	googleServiceAccountEmail := readEnvVarOrSkip(t, "TEST_GOOGLE_SERVICE_ACCOUNT_EMAIL")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(t) + "\n" +
					testAccExampleResourceConfigForGoogle("initial", snykOrgId, googleProjectId, googleServiceAccountEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snyk_environment.test_google", "name", "initial"),
					resource.TestCheckResourceAttr("snyk_environment.test_google", "kind", "google"),
					resource.TestCheckResourceAttr("snyk_environment.test_google", "google.project_id", googleProjectId),
					resource.TestCheckResourceAttr("snyk_environment.test_google", "google.service_account_email", googleServiceAccountEmail),
				),
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig(t) + "\n" +
					testAccExampleResourceConfigForGoogle("updated", snykOrgId, googleProjectId, googleServiceAccountEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snyk_environment.test_google", "name", "updated"),
					resource.TestCheckResourceAttr("snyk_environment.test_google", "google.project_id", googleProjectId),
					resource.TestCheckResourceAttr("snyk_environment.test_google", "google.service_account_email", googleServiceAccountEmail),
				),
			},
			// Delete testing automatically occurs in TestCase
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
