package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExampleOrganizationResource(t *testing.T) {
	snykGroupId := readEnvironmentVar("TEST_SNYK_GROUP_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccExampleOrganizationResourceRaw("Test snyk org", snykGroupId), // TODO: group id
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snyk_organization.test", "name", "Test snyk org"),
					resource.TestCheckResourceAttr("snyk_organization.test", "group_id", snykGroupId),
				),
			},
			// No Update API so no update test
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleOrganizationResourceRaw(orgName string, groupId string) string {
	return fmt.Sprintf(`
resource "snyk_organization" "test" {
  name = %[1]q
  group_id = %[2]q
}`, orgName, groupId)
}
