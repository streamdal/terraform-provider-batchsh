package provider

import (
	"context"

	"github.com/batchcorp/terraform-provider-batchsh/batch"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDatalake() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSoureDataLakeRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"id": {
				Description: "Batch Data Lake ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "Data Lake type: aws | azure | gcp",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Data Lake name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: "Data Lake status slug",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status_full": {
				Description: "Readable data lake status",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"inserted_at": {
				Description: "When Data Lake was created",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "When Data Lake was last updated",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSoureDataLakeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var filters []*batch.Filter

	client := m.(batch.IBatchAPI)

	if v, ok := d.GetOk("filter"); ok {
		filters = buildFiltersDataSource(v.(*schema.Set))
	} else {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "No filters defined",
			Detail:   "At least one filter must be defined",
		})
	}

	lake, moreDiags := client.GetDataLake(filters)
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	d.SetId(lake["id"].(string))
	d.Set("type", lake["type"].(string))
	d.Set("name", lake["name"].(string))
	d.Set("status", lake["status"].(string))
	d.Set("status_full", lake["status_full"].(string))
	d.Set("inserted_at", lake["inserted_at"].(string))
	d.Set("updated_at", lake["updated_at"].(string))

	return diags
}
