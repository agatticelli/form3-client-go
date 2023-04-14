package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/agatticelli/form3-client-go/form3"
	"github.com/google/uuid"
)

const existingAccountID = "d290f1ee-6c54-4b01-90e6-d701748f0851"

func accountFactory(ID, country string) *form3.CreateAccountData {
	return &form3.CreateAccountData{
		ID:             ID,
		OrganisationID: uuid.New().String(),
		Attributes: &form3.CreateAccountAttributes{
			BankID:        "400300",
			BankIDCode:    "GBDSC",
			Bic:           "NWBKGB22",
			Name:          []string{"Samantha Holder"},
			Country:       country,
			AccountNumber: form3.ToPointer("31926819"),
		},
	}
}

func accountsTableSeeder(t *testing.T, client *form3.Client) {
	t.Helper()

	// First, let's delete all records from previous failed tests
	truncateAccountsTable(t, client)

	createAccountData := accountFactory(existingAccountID, "GB")
	_, _, err := client.CreateAccount(context.Background(), createAccountData.ID, createAccountData.OrganisationID, createAccountData.Attributes)
	if err != nil {
		t.Fatalf("error seeding account: %v", err)
	}
}

func truncateAccountsTable(t *testing.T, client *form3.Client) {
	t.Helper()

	var accountList form3.Form3BodyResponse[[]form3.Account]
	err := client.Do(context.Background(), http.MethodGet, "organisation/accounts", nil, &accountList)
	if err != nil {
		t.Fatal(err)
	}

	for _, account := range accountList.Data {
		client.DeleteAccount(context.Background(), account.ID, *account.Version)
	}
}
