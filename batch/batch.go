package batch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . IBatchAPI
type IBatchAPI interface {
	CreateCollection(params *CreateCollectionRequest) (*CreateCollectionResponse, diag.Diagnostics)
	DeleteCollection(collectionID string) (*UpdateCollectionResponse, diag.Diagnostics)
	GetCollection(collectionID string) (*ReadCollectionResponse, diag.Diagnostics)
	UpdateCollection(collectionID string, params *UpdateCollectionRequest) (*UpdateCollectionResponse, diag.Diagnostics)

	CreateTeamMember(params *CreateTeamMemberRequest) (*CreateTeamMemberResponse, diag.Diagnostics)
	DeleteTeamMember(id string) diag.Diagnostics
	GetTeamMember(accountID string) (*ReadTeamMemberResponse, diag.Diagnostics)
	UpdateTeamMember(id string, params *UpdateTeamMemberRequest) (*UpdateTeamMemberResponse, diag.Diagnostics)

	GetDataLake(filters []*Filter) (map[string]interface{}, diag.Diagnostics)

	GetSchema(filters []*Filter) (map[string]interface{}, diag.Diagnostics)
}

type ApiClient struct {
	*Config
}

// ResponseError is used to unmarshal an error response from our API
type ResponseError struct {
	Code     int    `json:"code"`
	Domain   string `json:"domain"`
	Status   string `json:"status"`
	RawError string `json:"raw_error"`
	Message  string `json:"message"`
	Field    string `json:"field"`
}

// ErrorResponse is used to unmarshal an error response from our API
type ErrorResponse struct {
	Errors []*ResponseError `json:"errors"`
}

type Config struct {
	HttpClient  *http.Client
	ApiToken    string
	Version     string
	APIEndpoint string
}

func New(cfg *Config) (IBatchAPI, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return &ApiClient{
		Config: cfg,
	}, nil
}

func validateConfig(cfg *Config) error {
	if cfg.HttpClient == nil {
		cfg.HttpClient = &http.Client{Timeout: 10 * time.Second}
	}

	if cfg.ApiToken == "" {
		return errors.New("API Token cannot be empty")
	}

	return nil
}

// TODO: rate limiting, retry?
func (a *ApiClient) Request(method, endpoint string, payload []byte) ([]byte, diag.Diagnostics) {
	var diags diag.Diagnostics

	endpoint = fmt.Sprintf("%s%s", a.APIEndpoint, endpoint)
	req, err := http.NewRequest(method, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, diag.FromErr(err)
	}

	req.Header.Add("User-Agent", "terraform-provider-batchsh/"+a.Version)
	req.Header.Add("Authorization", a.ApiToken)

	r, err := a.HttpClient.Do(req)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read response body",
			Detail:   string(body),
		})
		return nil, diags
	}

	if r.StatusCode < 200 || r.StatusCode >= 400 {
		errResp := &ErrorResponse{}
		if err := json.Unmarshal(body, errResp); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to unmarshal error response",
				Detail:   string(body),
			})
			return nil, diags
		}

		for _, v := range errResp.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("API Error. Code: %d, Domain: %s", v.Code, v.Domain),
				Detail:   v.Message,
			})
		}
		return nil, diags
	}

	return body, diags
}
