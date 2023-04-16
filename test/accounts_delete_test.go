package integration

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/agatticelli/form3-client-go/form3"
	"github.com/google/uuid"
)

func Test_DeleteAccount(t *testing.T) {
	cases := []struct {
		name               string
		accountID          string
		version            int64
		expectAPIError     bool
		accountExistsCheck bool
		checker            func(t *testing.T, accountID string) error
	}{
		{
			name:               "Delete account successfully",
			accountID:          existingAccountID,
			version:            0,
			accountExistsCheck: true,
			checker:            deleteAccountExistsChecker,
		},
		{
			name:           "Delete account with invalid ID",
			accountID:      "invalid",
			version:        0,
			expectAPIError: true,
		},
		{
			name:           "Delete account with non-existing ID",
			accountID:      uuid.NewString(),
			version:        0,
			expectAPIError: true,
		},
		{
			name:           "Delete with invalid version",
			accountID:      existingAccountID,
			version:        1,
			expectAPIError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()

			client := initClient(t)

			// seed database for testing conflicts
			accountsTableSeeder(t, client)

			// first check that account exists
			if tc.accountExistsCheck {
				account, _, err := client.Account.Fetch(context.Background(), tc.accountID)
				if err != nil || account.ID != tc.accountID {
					t.Fatalf("account %s should exists", tc.accountID)
				}
			}

			// delete test case account
			err := client.Account.Delete(context.Background(), tc.accountID, tc.version)

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
				if err := tc.checker(t, tc.accountID); err != nil {
					t.Fatalf("error checking test case: %v", err)
				}
			}
		})
	}
}

func deleteAccountExistsChecker(t *testing.T, accountID string) error {
	t.Helper()

	_, _, err := client.Account.Fetch(context.Background(), accountID)
	if err == nil {
		return fmt.Errorf("account should not exists")
	}

	unwrappedErr := errors.Unwrap(err)
	form3Err, ok := unwrappedErr.(*form3.Form3APIError)
	if !ok {
		return fmt.Errorf("expected Form3APIError but got %T", err)
	}

	if form3Err.StatusCode != 404 {
		return fmt.Errorf("expected status code 404 but got %d", form3Err.StatusCode)
	}

	return nil
}
