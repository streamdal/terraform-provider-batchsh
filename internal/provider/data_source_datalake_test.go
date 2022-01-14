package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/batchcorp/terraform-provider-batchsh/batch"
	"github.com/batchcorp/terraform-provider-batchsh/batch/batchfakes"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataSourceDatalake(t *testing.T) {

	fakeBatch := &batchfakes.FakeIBatchAPI{}
	fakeBatch.GetDataLakeStub = func([]*batch.Filter) (map[string]interface{}, diag.Diagnostics) {
		return map[string]interface{}{
			"id":          "002f7f50-aadb-416d-aae6-56a1103273ba",
			"type":        "aws",
			"name":        "Default DataLake",
			"team_id":     "71d86e5f-686c-4802-84a7-f06f5925eae8",
			"status":      "active",
			"status_full": "",
			"archived":    false,
			"inserted_at": "2021-09-10T17:40:45.815438Z",
			"updated_at":  "2021-09-10T17:40:45.930814Z",
		}, diag.Diagnostics{}
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactory(fakeBatch),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDatalake,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.batchsh_datalake.foo", "name", "Default DataLake"),
				),
			},
		},
	})
}

const testDataSourceDatalake = `
data "batchsh_datalake" "foo" {
  filter {
    name = "name"
	values = ["Generic JSON"]
  }
}
`
