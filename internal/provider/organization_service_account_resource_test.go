package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExampleOrganizationServiceAccountResource(t *testing.T) {
	// TODO: don't skip this
	snykOrgId := readEnvVarOrFail(t, "TEST_SNYK_ORG_ID")
	snykRoleId := readEnvVarOrSkip(t, "TEST_SNYK_ROLE_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(t) + "\n" +
					testAccExampleOrganizationResourceRaw(snykOrgId, "test-snyk-org-service-account", snykRoleId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snyk_organization_service_account.test", "organization_id", snykOrgId),
					resource.TestCheckResourceAttr("snyk_organization_service_account.test", "name", "test-snyk-org-service-account"),
					resource.TestCheckResourceAttr("snyk_organization_service_account", "role_id", snykRoleId),
				),
			},
			// No Update API so no update test
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleOrganizationServiceAccountResourceRaw(orgId string, orgName string, roleId) string {
	return fmt.Sprintf(`
resource "snyk_organization_service_account" "test" {
  organization_id = %[1]q
  name = %[2]q
  auth_type = "api_key"
  role_id = %[3]q
}`, orgId, orgName, roleId)
}
