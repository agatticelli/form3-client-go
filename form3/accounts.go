package form3

import (
	"context"
	"net/http"
)

// HTTP entities
// Ref: https://www.api-docs.form3.tech/api/schemes/fps-direct/accounts/accounts/create-an-account
type CreateAccountRequest = Form3BodyRequest[CreateAccountData]
type CreateAccountData struct {
	ID             string                   `json:"id,omitempty"`
	OrganisationID string                   `json:"organisation_id,omitempty"`
	Type           string                   `json:"type,omitempty"`
	Attributes     *CreateAccountAttributes `json:"attributes,omitempty"`
}
type CreateAccountAttributes struct {
	BankID                  string    `json:"bank_id"`
	BankIDCode              string    `json:"bank_id_code"`
	Bic                     string    `json:"bic"`
	Country                 string    `json:"country"`
	Name                    []string  `json:"name,omitempty"`
	AccountClassification   *string   `json:"account_classification,omitempty"`
	AccountNumber           *string   `json:"account_number,omitempty"`
	AlternativeNames        *[]string `json:"alternative_names,omitempty"`
	BaseCurrency            *string   `json:"base_currency,omitempty"`
	Iban                    *string   `json:"iban,omitempty"`
	JointAccount            *bool     `json:"joint_account,omitempty"`
	SecondaryIdentification *string   `json:"secondary_identification,omitempty"`
}
type CreateAccountResponse = Form3BodyResponse[Account]

// Ref: https://www.api-docs.form3.tech/api/schemes/fps-direct/accounts/accounts/fetch-an-account
type FetchAccountResponse = Form3BodyResponse[Account]

// Business models
type Account struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	CreatedOn      string             `json:"created_on,omitempty"`
	ModifiedOn     string             `json:"modified_on,omitempty"`
	Version        *int64             `json:"version,omitempty"`
}
type AccountAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
}

func (c *Client) CreateAccount(ctx context.Context, ID string, organisationID string, attributes *CreateAccountAttributes) (*Account, *Form3BodyResponseLinks, error) {
	formData := CreateAccountRequest{
		Data: CreateAccountData{
			ID:             ID,
			OrganisationID: organisationID,
			Type:           "accounts",
			Attributes:     attributes,
		},
	}

	accountResponse := CreateAccountResponse{}
	err := c.Do(ctx, http.MethodPost, "organisation/accounts", formData, &accountResponse)
	if err != nil {
		return nil, nil, err
	}

	return &accountResponse.Data, &accountResponse.Links, nil
}
