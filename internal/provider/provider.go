package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/batchcorp/terraform-provider-batchsh/batch"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" (Default: `%v`)", s.Default)
		}
		return strings.TrimSpace(desc)
	}
}

func New(version, apiToken string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"token": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("BATCHSH_TOKEN", apiToken),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"batchsh_datalake": dataSourceDatalake(),
				"batchsh_schema":   dataSourceSchema(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"batchsh_collection":  resourceCollection(),
				"batchsh_team_member": resourceTeamMember(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		token := d.Get("token").(string)
		token = fmt.Sprintf("Bearer %s", token)

		b, err := batch.New(&batch.Config{
			ApiToken: token,
			Version:  version,
		})
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return b, nil
	}
}

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Description: "Field name to filter on",
					Type:        schema.TypeString,
					Required:    true,
				},

				"values": {
					Description: "Value(s) to filter by. Wildcards '*' are supported.",
					Type:        schema.TypeList,
					Required:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func buildFiltersDataSource(set *schema.Set) []*batch.Filter {
	var filters []*batch.Filter
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}
		filters = append(filters, &batch.Filter{
			Name:   m["name"].(string),
			Values: filterValues,
		})
	}
	return filters
}
