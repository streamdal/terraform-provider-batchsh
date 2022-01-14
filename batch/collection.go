package batch

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// ReadCollectionResponse is used to unmarshal a read response from the API
type ReadCollectionResponse struct {
	ID                  string `json:"id"`
	Token               string `json:"token"`
	Name                string `json:"name"`
	Notes               string `json:"notes"`
	Archived            bool   `json:"archived"`
	SchemaID            string `json:"schema_id"`
	DatalakeID          string `json:"datalake_id"`
	EnvelopeType        string `json:"envelope_type"`
	EnvelopeRootMessage string `json:"envelope_root_message"`
	PayloadRootMessage  string `json:"payload_root_message"`
	PayloadFieldID      int    `json:"payload_field_id"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

type CreateCollectionResponse struct {
	ID        string `json:"id"`
	Token     string `json:"token"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateCollectionRequest struct {
	Name                string `json:"name"`
	Notes               string `json:"notes"`
	SchemaID            string `json:"schema_id"`
	DatalakeID          string `json:"datalake_id"`
	EnvelopeType        string `json:"envelope_type"`
	EnvelopeRootMessage string `json:"envelope_root_message"`
	PayloadRootMessage  string `json:"payload_root_message"`
	PayloadFieldID      int    `json:"payload_field_id"`
}

type UpdateCollectionRequest struct {
	CollectionID string `json:"collection_id"`
	Name         string `json:"name,omitempty"`
	Notes        string `json:"notes,omitempty"`
	Archived     bool   `json:"archived,omitempty"`
}

type UpdateCollectionResponse struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Notes               string `json:"notes"`
	SchemaID            string `json:"schema_id"`
	DatalakeID          string `json:"datalake_id"`
	EnvelopeType        string `json:"envelope_type"`
	EnvelopeRootMessage string `json:"envelope_root_message"`
	PayloadRootMessage  string `json:"payload_root_message"`
	PayloadFieldID      int    `json:"payload_field_id"`
	Paused              bool   `json:"paused"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	Archived            bool   `json:"archived"`
}

// GetCollection obtains data for a single collection
func (a *ApiClient) GetCollection(collectionID string) (*ReadCollectionResponse, diag.Diagnostics) {
	var diags diag.Diagnostics

	body, moreDiags := a.Request("GET", fmt.Sprintf("/collection/%s", collectionID), nil)
	if moreDiags.HasError() {
		diags = append(diags, moreDiags...)
		return nil, diags
	}

	resp := &ReadCollectionResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to unmarshal read collection response",
			Detail:   string(body),
		})
		return nil, diags
	}

	return resp, diags
}

func (a *ApiClient) CreateCollection(params *CreateCollectionRequest) (*CreateCollectionResponse, diag.Diagnostics) {
	var diags diag.Diagnostics

	reqBody, err := json.Marshal(params)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to marshal create collection request",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	body, moreDiags := a.Request("POST", "/collection", reqBody)
	if moreDiags.HasError() {
		diags = append(diags, moreDiags...)
		return nil, diags
	}

	resp := &CreateCollectionResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to marshal create collection response",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return resp, diags
}

func (a *ApiClient) UpdateCollection(collectionID string, params *UpdateCollectionRequest) (*UpdateCollectionResponse, diag.Diagnostics) {
	var diags diag.Diagnostics

	reqBody, err := json.Marshal(params)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to marshal update collection request",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	body, moreDiags := a.Request("PUT", "/collection", reqBody)
	if moreDiags.HasError() {
		diags = append(diags, moreDiags...)
		return nil, diags
	}

	resp := &UpdateCollectionResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to unmarshal update collection response",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return resp, diags
}

func (a *ApiClient) DeleteCollection(collectionID string) (*UpdateCollectionResponse, diag.Diagnostics) {
	var diags diag.Diagnostics

	params := &UpdateCollectionRequest{
		CollectionID: collectionID,
		Archived:     true,
	}

	reqBody, err := json.Marshal(params)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to marshal delete collection request",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	body, moreDiags := a.Request("PUT", "/collection", reqBody)
	if moreDiags.HasError() {
		diags = append(diags, moreDiags...)
		return nil, diags
	}

	resp := &UpdateCollectionResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to unmarshal delete collection response",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return resp, diags
}
