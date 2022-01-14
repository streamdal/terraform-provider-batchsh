package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/batchcorp/terraform-provider-batchsh/batch"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){}

// ProviderFactory us used to return a map of providers with the Batch API client injected into it
func ProviderFactory(fakeApiClient batch.IBatchAPI) map[string]func() (*schema.Provider, error) {
	factories := map[string]func() (*schema.Provider, error){}
	factories["batchsh"] = func() (*schema.Provider, error) {
		p := New("dev", "test")()

		p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			return fakeApiClient, nil
		}

		return p, nil
	}

	return factories
}

func TestProvider(t *testing.T) {
	if err := New("dev", "test")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
