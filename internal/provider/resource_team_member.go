package provider

import (
	"context"

	"github.com/batchcorp/terraform-provider-batchsh/batch"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeamMember() *schema.Resource {
	return &schema.Resource{
		Description: "Team Members",

		CreateContext: resourceTeamMemberCreate,
		ReadContext:   resourceTeamMemberRead,
		UpdateContext: resourceTeamMemberUpdate,
		DeleteContext: resourceTeamMemberDelete,

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Team Member ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Member Name",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"email": {
				Description: "Member Email",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"password": {
				Description: "Member Password",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"roles": {
				Description: "Member Roles",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTeamMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(batch.IBatchAPI)

	resp, moreDiags := client.GetTeamMember(d.Id())
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	d.SetId(resp.Id)
	d.Set("name", resp.Name)
	d.Set("email", resp.Email)
	d.Set("roles", resp.Roles)

	return diags
}

func resourceTeamMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	params := &batch.CreateTeamMemberRequest{
		Name:     d.Get("name").(string),
		Email:    d.Get("email").(string),
		Password: d.Get("password").(string),
		Roles:    flattenRoles(d.Get("roles")),
	}

	client := m.(batch.IBatchAPI)

	resp, moreDiags := client.CreateTeamMember(params)
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	d.SetId(resp.Id)

	return diags
}

func resourceTeamMemberUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	params := &batch.UpdateTeamMemberRequest{
		Name:     d.Get("name").(string),
		Email:    d.Get("email").(string),
		Password: d.Get("password").(string),
		Roles:    flattenRoles(d.Get("roles")),
	}

	client := m.(batch.IBatchAPI)

	resp, moreDiags := client.UpdateTeamMember(d.Id(), params)
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	d.SetId(resp.Id)

	return diags
}

func resourceTeamMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(batch.IBatchAPI)

	moreDiags := client.DeleteTeamMember(d.Id())
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	d.SetId("")

	return diags
}

// flattenRoles interface{} to []string for team member roles
func flattenRoles(input interface{}) []string {
	roles := input.([]interface{})

	var result []string

	for _, role := range roles {
		result = append(result, role.(string))
	}

	return result
}
