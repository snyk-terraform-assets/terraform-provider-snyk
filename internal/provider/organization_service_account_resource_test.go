package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExampleOrganizationServiceAccountResource(t *testing.T) {
	snykOrgId := readEnvVarOrFail(t, "TEST_SNYK_ORG_ID")
	snykRoleId := readEnvVarOrSkip(t, "TEST_SNYK_ROLE_ID")

  snykName :=fmt.Sprintf("Test snyk service account%d", time.Now().Unix())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(t) + "\n" + fmt.Sprintf(`
resource "snyk_organization_service_account" "test" {
  organization_id = %[1]q
  name = %[2]q
  auth_type = "api_key"
  role_id = %[3]q
}`, snykOrgId, snykName , snykRoleId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snyk_organization_service_account.test", "organization_id", snykOrgId),
					resource.TestCheckResourceAttr("snyk_organization_service_account.test", "name", snykName),
					resource.TestCheckResourceAttr("snyk_organization_service_account.test", "role_id", snykRoleId),
				),
			},
			// No Update API so no update test
			// Delete testing automatically occurs in TestCase
		},
	})
}
