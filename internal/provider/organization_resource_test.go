package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExampleOrganizationResource(t *testing.T) {
	// TODO: don't skip this
	snykGroupId := readEnvVarOrSkip(t, "TEST_SNYK_GROUP_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig(t) + "\n" +
					testAccExampleOrganizationResourceRaw("Test snyk org", snykGroupId),
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
