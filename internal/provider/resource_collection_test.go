package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/batchcorp/terraform-provider-batchsh/batch"
	"github.com/batchcorp/terraform-provider-batchsh/batch/batchfakes"
)

func TestResourceCollection(t *testing.T) {

	id := "1fba7807-e78e-4ad4-a7de-0fb96b077d70"
	token := "0e82b53d-3115-40ea-a65c-1c6ca701a210"
	name := "My TF Managed Collection"
	notes := "Any notes you wish to keep about this collection"
	schemaID := "d0d58226-d0a3-4b54-bcec-c88c43258968"
	datalakeID := "79a60dd1-55ac-4fdf-a305-0bc5902ac0c1"
	ts := time.Unix(0, 0).Format(time.RFC3339)

	testResourceCollection := fmt.Sprintf(`
resource "batchsh_collection" "test" {
  name                  = "%s"
  notes                 = "%s"
  envelope_type         = "deep"
  envelope_root_message = "events.Message"
  schema_id             = "%s"
  datalake_id           = "%s"
}
`, name, notes, schemaID, datalakeID)

	fakeBatch := &batchfakes.FakeIBatchAPI{}
	fakeBatch.CreateCollectionStub = func(*batch.CreateCollectionRequest) (*batch.CreateCollectionResponse, diag.Diagnostics) {
		return &batch.CreateCollectionResponse{
			ID:        id,
			Token:     token,
			CreatedAt: ts,
			UpdatedAt: ts,
		}, nil
	}
	fakeBatch.GetCollectionStub = func(string) (*batch.ReadCollectionResponse, diag.Diagnostics) {
		return &batch.ReadCollectionResponse{
			ID:                  id,
			Token:               token,
			Name:                name,
			Notes:               notes,
			Archived:            false,
			SchemaID:            schemaID,
			DatalakeID:          datalakeID,
			EnvelopeType:        "deep",
			EnvelopeRootMessage: "event.Message",
			PayloadRootMessage:  "",
			PayloadFieldID:      0,
			CreatedAt:           ts,
			UpdatedAt:           ts,
		}, diag.Diagnostics{}
	}
	fakeBatch.DeleteCollectionStub = func(string) (*batch.UpdateCollectionResponse, diag.Diagnostics) {
		return &batch.UpdateCollectionResponse{
			ID:                  id,
			Name:                name,
			Notes:               notes,
			Archived:            true,
			DatalakeID:          datalakeID,
			EnvelopeType:        "deep",
			EnvelopeRootMessage: "events.Message",
			PayloadRootMessage:  "",
			PayloadFieldID:      0,
			UpdatedAt:           ts,
		}, diag.Diagnostics{}
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactory(fakeBatch),
		Steps: []resource.TestStep{
			{
				ResourceName:  "batchsh_collection.test",
				Config:        testResourceCollection,
				ImportState:   true,
				ImportStateId: id,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"batchsh_collection.test", "id", id),
					resource.TestCheckResourceAttr(
						"batchsh_collection.test", "token", token),
					resource.TestCheckResourceAttr(
						"batchsh_collection.test", "notes", notes),
					resource.TestCheckResourceAttr(
						"batchsh_collection.test", "schema_id", schemaID),
					resource.TestCheckResourceAttr(
						"batchsh_collection.test", "datalake_id", datalakeID),
					resource.TestCheckResourceAttr(
						"batchsh_collection.test", "envelope_type", "deep"),
					resource.TestCheckResourceAttr(
						"batchsh_collection.test", "envelope_root_message", "events.Message"),
					resource.TestCheckResourceAttr(
						"batchsh_collection.test", "created_at", ts),
					resource.TestCheckResourceAttr(
						"batchsh_collection.test", "updated_at", ts),
				),
			},
		},
	})
}
