package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/batchcorp/terraform-provider-batchsh/batch"
	"github.com/batchcorp/terraform-provider-batchsh/batch/batchfakes"
)

func TestDataSourceSchema(t *testing.T) {

	fakeBatch := &batchfakes.FakeIBatchAPI{}
	fakeBatch.GetSchemaStub = func([]*batch.Filter) (map[string]interface{}, diag.Diagnostics) {
		return map[string]interface{}{
			"id":          "f01bce7f-2ff1-4d4e-b2ed-9500678267a7",
			"name":        "Generic JSON",
			"type":        "json",
			"team_id":     "71d86e5f-686c-4802-84a7-f06f5925eae8",
			"shared":      false,
			"archived":    false,
			"inserted_at": "2021-09-10T17:40:45.815438Z",
			"updated_at":  "2021-09-10T17:40:45.815438Z",
		}, diag.Diagnostics{}
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactory(fakeBatch),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSchema,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.batchsh_schema.foo", "name", "Generic JSON"),
				),
			},
		},
	})
}

const testDataSourceSchema = `
data "batchsh_schema" "foo" {
  filter {
    name = "name"
	values = ["* JSON"]
  }
}
`
