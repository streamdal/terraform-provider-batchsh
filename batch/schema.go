package batch

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func (a *ApiClient) GetSchema(filters []*Filter) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	resp, moreDiags := a.Request(http.MethodGet, "/v1/schema", nil)
	if moreDiags.HasError() {
		return nil, append(diags, moreDiags...)
	}

	raw := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(resp, &raw); err != nil {
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse response",
			Detail:   err.Error(),
		})
	}

	schemas, diags := filterJSON(raw, filters)
	if diags.HasError() {
		return nil, diags
	}
	if len(schemas) < 1 {
		// No schema found using filter
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to find schema",
		})
	} else if len(schemas) > 1 {
		// Filter must find only one schema
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Filter returned more than one schema",
		})
	}

	return schemas[0], diags
}
