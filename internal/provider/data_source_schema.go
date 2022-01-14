package provider

import (
	"context"

	"github.com/batchcorp/terraform-provider-batchsh/batch"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSoureSchemaRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"id": {
				Description: "Batch Schema ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "Schema type",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Schema name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"inserted_at": {
				Description: "When schema was created",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "When schema was last updated",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSoureSchemaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	schema, moreDiags := client.GetSchema(filters)
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	d.SetId(schema["id"].(string))
	d.Set("name", schema["name"].(string))
	d.Set("type", schema["type"].(string))
	d.Set("inserted_at", schema["inserted_at"].(string))
	d.Set("updated_at", schema["updated_at"].(string))

	return diags
}
