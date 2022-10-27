package batch

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func (a *ApiClient) GetDataLake(filters []*Filter) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	resp, moreDiags := a.Request(http.MethodGet, "/v1/datalake", nil)
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

	lakes, diags := filterJSON(raw, filters)
	if diags.HasError() {
		return nil, diags
	}
	if len(lakes) < 1 {
		// No datalake found using filter
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to find datalake",
		})
	} else if len(lakes) > 1 {
		// Filter must find only one lake
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Filter returned more than one data lake",
		})
	}

	return lakes[0], diags
}
