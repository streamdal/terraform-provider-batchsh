package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/batchcorp/terraform-provider-batchsh/batch"
	"github.com/batchcorp/terraform-provider-batchsh/batch/batchfakes"
)

func TestResourceTeamMember(t *testing.T) {

	id := "76db6d2e-cead-410d-ab04-c6b95a5ebb6d"

	fakeBatch := &batchfakes.FakeIBatchAPI{}
	fakeBatch.CreateTeamMemberStub = func(*batch.CreateTeamMemberRequest) (*batch.CreateTeamMemberResponse, diag.Diagnostics) {
		return &batch.CreateTeamMemberResponse{
			Id:    id,
			Name:  "Johnny User",
			Email: "johnny@batch.sh",
			Roles: []string{"member"},
		}, diag.Diagnostics{}
	}
	fakeBatch.GetTeamMemberStub = func(string) (*batch.ReadTeamMemberResponse, diag.Diagnostics) {
		return &batch.ReadTeamMemberResponse{
			Id:    id,
			Name:  "Johnny User",
			Email: "johnny@batch.sh",
			Roles: []string{"member"},
		}, diag.Diagnostics{}
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactory(fakeBatch),
		Steps: []resource.TestStep{
			{
				ResourceName:  "batchsh_team_member.johnny",
				Config:        testResourceTeamMember,
				ImportState:   true,
				ImportStateId: id,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"batchsh_team_member.johnny", "id", id),
					resource.TestCheckResourceAttr(
						"batchsh_team_member.johnny", "name", "Johnny User"),
					resource.TestCheckResourceAttr(
						"batchsh_team_member.johnny", "email", "johnny@batch.sh"),
				),
			},
		},
	})
}

const testResourceTeamMember = `
resource "batchsh_team_member" "johnny" {
  name     = "Johnny User"
  email    = "johnny@batch.sh"
  password = "./password123"
  roles    = ["member"]
}
`
