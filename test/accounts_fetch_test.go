package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/agatticelli/form3-client-go/form3"
	"github.com/google/uuid"
)

func Test_FetchAccount(t *testing.T) {
	cases := []struct {
		name           string
		accountID      string
		expectAPIError bool
		checker        func(t *testing.T, fetchedAccount *form3.Account) error
	}{
		{
			name:      "Fetch account successfully",
			accountID: existingAccountID,
			checker:   fetchAccountExistsChecker,
		},
		{
			name:           "Fetch account with invalid ID",
			accountID:      "invalid",
			expectAPIError: true,
		},
		{
			name:           "Fetch account with non-existing ID",
			accountID:      uuid.NewString(),
			expectAPIError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()

			client := initClient(t)

			// seed database for testing conflicts
			accountsTableSeeder(t, client)

			// fetch test case account
			account, _, err := client.Account.Fetch(context.Background(), tc.accountID)

			// check errors
			if tc.expectAPIError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if _, ok := err.(*form3.Form3APIError); !ok {
					t.Fatalf("expected Form3APIError but got %T", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// check results
			if tc.checker != nil {
				if err := tc.checker(t, account); err != nil {
					t.Fatalf("error checking test case: %v", err)
				}
			}
		})
	}
}

func fetchAccountExistsChecker(t *testing.T, fetchedAccount *form3.Account) error {
	t.Helper()

	if fetchedAccount.ID != existingAccountID {
		return fmt.Errorf("account ID does not match: %s != %s", fetchedAccount.ID, existingAccountID)
	}

	if fetchedAccount.OrganisationID == "" {
		return fmt.Errorf("account organisation ID is empty")
	}

	if fetchedAccount.Attributes == nil {
		return fmt.Errorf("account attributes is nil")
	}

	if fetchedAccount.Version == nil {
		return fmt.Errorf("account version is empty")
	}

	if fetchedAccount.CreatedOn == "" {
		return fmt.Errorf("account created on is empty")
	}

	return nil
}
