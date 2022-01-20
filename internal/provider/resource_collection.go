package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/batchcorp/terraform-provider-batchsh/batch"
)

func resourceCollection() *schema.Resource {
	return &schema.Resource{
		Description: "Message Collections",

		CreateContext: resourceCollectionCreate,
		ReadContext:   resourceCollectionRead,
		UpdateContext: resourceCollectionUpdate,
		DeleteContext: resourceCollectionDelete,

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Collection ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name",
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
			},
			"notes": {
				Description: "Customer Notes",
				Type:        schema.TypeString,
				Computed:    false,
				Optional:    true,
			},
			"token": {
				Description: "Ingestion token",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"paused": {
				Description: "Is the collection paused",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"archived": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"envelope_type": {
				Description: "Protobuf envelope type: deep or shallow",
				Type:        schema.TypeString,
				Default:     "deep",
				Optional:    true,
			},
			"envelope_root_message": {
				Description: "Protobuf message name of the root message",
				Type:        schema.TypeString,
				Computed:    false,
				Optional:    true,
			},
			"payload_root_message": {
				Description: "Protobuf message name of the shallow envelope payload",
				Type:        schema.TypeString,
				Computed:    false,
				Optional:    true,
			},
			"payload_field_id": {
				Description: "Shallow envelope's protobuf field ID",
				Type:        schema.TypeInt,
				Computed:    false,
				Optional:    true,
			},
			"schema_id": {
				Description: "Batch Schema ID",
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
			},
			"datalake_id": {
				Description: "Batch Datalake ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"inserted_at": {
				Description: "Date the collection was created",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "Date the collection was last updated",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCollectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(batch.IBatchAPI)

	c, moreDiags := client.GetCollection(d.Id())
	if diags.HasError() {
		return append(diags, moreDiags...)
	}

	// Archived collections should be treated as deleted
	if c.Archived {
		d.SetId("")
		return diags
	}

	// These are currently the only values that can be updated on a c
	d.Set("name", c.Name)
	d.Set("notes", c.Notes)

	return diags
}

func resourceCollectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	params := &batch.CreateCollectionRequest{
		Name:                d.Get("name").(string),
		Notes:               d.Get("notes").(string),
		EnvelopeType:        d.Get("envelope_type").(string),
		EnvelopeRootMessage: d.Get("envelope_root_message").(string),
		PayloadRootMessage:  d.Get("payload_root_message").(string),
		PayloadFieldID:      d.Get("payload_field_id").(int),
		SchemaID:            d.Get("schema_id").(string),
		DatalakeID:          d.Get("datalake_id").(string),
	}

	client := m.(batch.IBatchAPI)

	resp, moreDiags := client.CreateCollection(params)
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	d.SetId(resp.ID)
	d.Set("token", resp.Token)
	d.Set("created_at", resp.CreatedAt)
	d.Set("updated_at", resp.UpdatedAt)

	return diags
}

func resourceCollectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	params := &batch.UpdateCollectionRequest{
		CollectionID: d.Id(),
		Name:         d.Get("name").(string),
		Notes:        d.Get("notes").(string),
	}

	client := m.(batch.IBatchAPI)

	resp, moreDiags := client.UpdateCollection(d.Id(), params)
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	d.Set("name", resp.Name)
	d.Set("notes", resp.Notes)
	d.Set("updated_at", resp.UpdatedAt)

	return diags
}

func resourceCollectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(batch.IBatchAPI)

	resp, moreDiags := client.DeleteCollection(d.Id())
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	// Archive successful,set ID to empty to indicate resource no longer exists
	if resp.Archived == true {
		d.SetId("")
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to archive collection",
		})
	}

	return diags
}
