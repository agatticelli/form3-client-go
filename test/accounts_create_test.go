package integration

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/agatticelli/form3-client-go/form3"
	"github.com/google/uuid"
)

func Test_CreateAccount(t *testing.T) {
	cases := []struct {
		name            string
		accountToCreate *form3.CreateAccountData
		checker         func(*testing.T, *form3.Client, *form3.CreateAccountData) error
		expectAPIError  bool
	}{
		{
			name:            "Create account successfully",
			accountToCreate: accountFactory(uuid.NewString(), "FR"),
			checker:         createAccountExistsChecker,
		},
		{
			name:            "Create account with invalid ID",
			accountToCreate: accountFactory("invalid", "FR"),
			expectAPIError:  true,
		},
		{
			name:            "Create account with missing required field",
			accountToCreate: accountFactory(uuid.NewString(), ""),
			expectAPIError:  true,
		},
		{
			name:            "Create account with already existing ID",
			accountToCreate: accountFactory(existingAccountID, "FR"),
			expectAPIError:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()

			client := initClient(t)

			// seed database for testing conflicts
			accountsTableSeeder(t, client)

			// create test case account
			_, _, err := client.Account.Create(context.Background(), tc.accountToCreate.ID, tc.accountToCreate.OrganisationID, tc.accountToCreate.Attributes)

			// check errors
			if tc.expectAPIError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if errors.Is(err, &form3.Form3APIError{}) {
					t.Fatalf("expected Form3APIError but got %T", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// check results
			if tc.checker != nil {
				if err := tc.checker(t, client, tc.accountToCreate); err != nil {
					t.Fatalf("error checking test case: %v", err)
				}
			}
		})
	}
}

func createAccountExistsChecker(t *testing.T, client *form3.Client, createdAccount *form3.CreateAccountData) error {
	t.Helper()

	account, _, err := client.Account.Fetch(context.Background(), createdAccount.ID)
	if err != nil {
		return fmt.Errorf("error fetching account: %v", err)
	}

	if account.ID != createdAccount.ID {
		return fmt.Errorf("account ID does not match: %s != %s", account.ID, createdAccount.ID)
	}

	if account.OrganisationID != createdAccount.OrganisationID {
		return fmt.Errorf("account OrganisationID does not match: %s != %s", account.OrganisationID, createdAccount.OrganisationID)
	}

	return nil
}
