package batch

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type ReadTeamMemberResponse struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	AvailableTeams []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"available_teams"`
	Roles []string `json:"roles"`
}

type CreateTeamMemberRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

type UpdateTeamMemberRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

type UpdateTeamMemberResponse struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	AvailableTeams []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"available_teams"`
	Roles []string `json:"roles"`
}

type CreateTeamMemberResponse struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	AvailableTeams []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"available_teams"`
	Roles []string `json:"roles"`
}

func (a *ApiClient) GetTeamMember(accountID string) (*ReadTeamMemberResponse, diag.Diagnostics) {
	var diags diag.Diagnostics

	body, moreDiags := a.Request(http.MethodGet, fmt.Sprintf("/v1/team/member/%s", accountID), nil)
	if moreDiags.HasError() {
		diags = append(diags, moreDiags...)
		return nil, diags
	}

	resp := &ReadTeamMemberResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to unmarshal read team member response",
			Detail:   string(body),
		})
		return nil, diags
	}

	return resp, diags
}

func (a *ApiClient) CreateTeamMember(params *CreateTeamMemberRequest) (*CreateTeamMemberResponse, diag.Diagnostics) {
	var diags diag.Diagnostics

	payload, err := json.Marshal(params)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to marshal create team member request",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	body, moreDiags := a.Request(http.MethodPost, "/v1/team/member", payload)
	if moreDiags.HasError() {
		return nil, append(diags, moreDiags...)
	}

	resp := &CreateTeamMemberResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to marshal create team member response",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return resp, diags
}

func (a *ApiClient) UpdateTeamMember(id string, params *UpdateTeamMemberRequest) (*UpdateTeamMemberResponse, diag.Diagnostics) {
	var diags diag.Diagnostics

	payload, err := json.Marshal(params)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to marshal update team member request",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	body, moreDiags := a.Request(http.MethodPut, "/v1/team/member/"+id, payload)
	if moreDiags.HasError() {
		return nil, append(diags, moreDiags...)
	}

	resp := &UpdateTeamMemberResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to marshal update team member response",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return resp, diags
}

func (a *ApiClient) DeleteTeamMember(id string) diag.Diagnostics {
	var diags diag.Diagnostics

	_, moreDiags := a.Request(http.MethodDelete, fmt.Sprintf("/v1/team/member/%s", id), nil)
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	return diags
}
