package batch

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/minio/pkg/wildcard"
)

type Filter struct {
	Name   string
	Values []string
}

// filterJSON takes a slice of map[string]interface{}, which represents the output of an API call to
// https://api.batch.sh endpoint which returns a JSON array of objects(collections/schemas/etc). It then applies
// the given filters to the JSON array and returns only entries in the data param which have a field that
// matches a filter.
func filterJSON(data []map[string]interface{}, filters []*Filter) ([]map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	found := make([]map[string]interface{}, 0)

	for _, filter := range filters {
		for _, item := range data {
			// Can't match on non-existent keys
			if _, ok := item[filter.Name]; !ok {
				return nil, diag.FromErr(fmt.Errorf("%s is not a valid key", filter.Name))
			}

			checkVal := fmt.Sprintf("%s", item[filter.Name])
			for _, val := range filter.Values {
				// Wildcard match
				if matches(val, checkVal) {
					found = append(found, item)
				}
			}
		}
	}

	return found, diags
}

func matches(val, checkVal string) bool {
	if strings.Contains(val, "*") && wildcard.MatchSimple(val, checkVal) {
		return true
	}

	if val == checkVal {
		return true
	}

	return false
}
