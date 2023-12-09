package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccProjectInboundDataFilterResource(t *testing.T) {
	rn := "sentry_project_inbound_data_filter.test"
	teamSlug := acctest.RandomWithPrefix("tf-team")
	projectSlug := acctest.RandomWithPrefix("tf-project")
	filterId := "browser-extensions"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTeamMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectInboundDataFilterConfig(teamSlug, projectSlug, filterId, "active = true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project_slug", projectSlug),
					resource.TestCheckResourceAttr(rn, "filter_id", filterId),
					resource.TestCheckResourceAttr(rn, "active", "true"),
					resource.TestCheckNoResourceAttr(rn, "subfilters"),
				),
			},
			{
				Config: testAccProjectInboundDataFilterConfig(teamSlug, projectSlug, filterId, "active = false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project_slug", projectSlug),
					resource.TestCheckResourceAttr(rn, "filter_id", filterId),
					resource.TestCheckResourceAttr(rn, "active", "false"),
					resource.TestCheckNoResourceAttr(rn, "subfilters"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProjectInboundDataFilterResource_LegacyBrowser(t *testing.T) {
	rn := "sentry_project_inbound_data_filter.test"
	teamSlug := acctest.RandomWithPrefix("tf-team")
	projectSlug := acctest.RandomWithPrefix("tf-project")
	filterId := "legacy-browsers"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTeamMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectInboundDataFilterConfig(teamSlug, projectSlug, filterId, "subfilters = [\"ie_pre_9\"]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project_slug", projectSlug),
					resource.TestCheckResourceAttr(rn, "filter_id", filterId),
					resource.TestCheckNoResourceAttr(rn, "active"),
					resource.TestCheckResourceAttr(rn, "subfilters.#", "1"),
					resource.TestCheckResourceAttr(rn, "subfilters.0", "ie_pre_9"),
				),
			},
			{
				Config: testAccProjectInboundDataFilterConfig(teamSlug, projectSlug, filterId, "subfilters = [\"safari_pre_6\"]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project_slug", projectSlug),
					resource.TestCheckResourceAttr(rn, "filter_id", filterId),
					resource.TestCheckNoResourceAttr(rn, "active"),
					resource.TestCheckResourceAttr(rn, "subfilters.#", "1"),
					resource.TestCheckResourceAttr(rn, "subfilters.0", "safari_pre_6"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectInboundDataFilterConfig(teamName, projectName, filterId, body string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
  organization = data.sentry_organization.test.id
  name         = "%[1]s"
  slug         = "%[1]s"
}

resource "sentry_project" "test" {
  organization = sentry_team.test.organization
  teams        = [sentry_team.test.id]
  name         = "%[2]s"
  platform     = "go"
}

resource "sentry_project_inbound_data_filter" "test" {
  organization = sentry_team.test.organization
  project_slug = sentry_project.test.slug
  filter_id    = "%[3]s"
  %[4]s
}
`, teamName, projectName, filterId, body)
}
