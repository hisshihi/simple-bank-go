package repo

import (
	"context"
	"testing"

	"github.com/hisshihi/simple-bank/core/db"
	"gorm.io/gorm"
)

func TestAccountCreate(t *testing.T) {
	tests := []struct {
		name     string
		owner    string
		balance  float64
		currency string
		wantErr  error
	}{
		{
			name:     "success - RUB",
			owner:    "hiss",
			balance:  1000.00,
			currency: db.RUB.String(),
			wantErr:  nil,
		},
		{
			name:     "success - USD",
			owner:    "john",
			balance:  500.50,
			currency: db.USD.String(),
			wantErr:  nil,
		},
		{
			name:     "success - zero balance",
			owner:    "arina",
			balance:  0.00,
			currency: db.EUR.String(),
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB := setupTestDB(t)
			repo := NewAccountRepo(testDB)
			ctx := context.Background()

			// Act
			err := repo.Create(ctx, tt.owner, tt.balance, tt.currency)

			// Assert
			if tt.wantErr != nil {
				assertErrorIs(t, err, tt.wantErr)
			}

			assertNoError(t, err)

			// Verify the account was created
			count, err := gorm.G[db.Account](testDB).Where("owner = ?", tt.owner).Count(ctx, "owner")
			assertNoError(t, err)
			if count != 1 {
				t.Fatalf("expected 1 account, got %d", count)
			}
		})
	}
}

func TestFindAccountByID(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T, db *gorm.DB) uint
		accountID uint
		wantErr   error
	}{
		{
			name: "success - account exists",
			setup: func(t *testing.T, database *gorm.DB) uint {
				account, err := createAccountAndReturnResult(t, database, "hiss", 500.00, db.USD)
				if err != nil {
					t.Fatalf("failed to create account: %v", err)
				}
				return account.ID
			},
			wantErr: nil,
		},
		{
			name: "error - account not found",
			setup: func(t *testing.T, database *gorm.DB) uint {
				return 999 // несуществующий ID
			},
			accountID: 999,
			wantErr:   db.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB := setupTestDB(t)
			repo := NewAccountRepo(testDB)
			ctx := context.Background()

			accountID := tt.accountID
			if tt.setup != nil {
				accountID = tt.setup(t, testDB)
			}

			// Act
			account, err := repo.FindByID(ctx, accountID)

			// Assert
			if tt.wantErr != nil {
				assertErrorIs(t, err, tt.wantErr)
				return
			}

			assertNoError(t, err)
			if account.ID != accountID {
				t.Fatalf("expected account ID %d, got %d", accountID, account.ID)
			}
		})
	}
}

func TestFindAllAccounts(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(t *testing.T, db *gorm.DB)
		accountCount int
		wantErr      error
	}{
		{
			name: "success - multiple accounts",
			setup: func(t *testing.T, database *gorm.DB) {
				_, err := createAccountAndReturnResult(t, database, "hiss", 500.00, db.USD)
				assertNoError(t, err)
				_, err = createAccountAndReturnResult(t, database, "hiss", 500.00, db.USD)
				assertNoError(t, err)
			},
			accountCount: 2,
			wantErr:      nil,
		},
		{
			name: "success - no accounts",
			setup: func(t *testing.T, db *gorm.DB) {
				// ничего не создаем
			},
			accountCount: 0,
			wantErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB := setupTestDB(t)
			repo := NewAccountRepo(testDB)
			ctx := context.Background()

			if tt.setup != nil {
				tt.setup(t, testDB)
			}

			// Act
			accounts, err := repo.FindAllAccounts(ctx)

			// Assert
			if tt.wantErr != nil {
				assertErrorIs(t, err, tt.wantErr)
				return
			}

			assertNoError(t, err)
			if len(accounts) != tt.accountCount {
				t.Fatalf("expected %d accounts, got %d", tt.accountCount, len(accounts))
			}
		})
	}
}

func TestUpdateAccountBalance(t *testing.T) {
	tests := []struct {
		name       string
		accountID  uint
		newBalance float64
		setup      func(t *testing.T, db *gorm.DB) uint
		wantErr    error
	}{
		{
			name:       "success - update balance",
			accountID:  1,
			newBalance: 2000.00,
			setup: func(t *testing.T, database *gorm.DB) uint {
				account, err := createAccountAndReturnResult(t, database, "hiss", 1000.00, db.RUB)
				assertNoError(t, err)
				return account.ID
			},
			wantErr: nil,
		},
		{
			name:       "error - account not found",
			accountID:  999,
			newBalance: 1500.00,
			setup: func(t *testing.T, db *gorm.DB) uint {
				// ничего не создаем
				return 999
			},
			wantErr: db.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB := setupTestDB(t)
			repo := NewAccountRepo(testDB)
			ctx := context.Background()

			accountID := tt.accountID
			if tt.setup != nil {
				accountID = tt.setup(t, testDB)
			}

			// Act
			err := repo.UpdateBalance(ctx, accountID, tt.newBalance)

			// Assert
			if tt.wantErr != nil {
				assertErrorIs(t, err, tt.wantErr)
				return
			}

			assertNoError(t, err)

			// Verify the balance was updated
			updatedAccount, err := repo.FindByID(ctx, accountID)
			assertNoError(t, err)
			if updatedAccount.Balance != tt.newBalance {
				t.Fatalf("expected balance %f, got %f", tt.newBalance, updatedAccount.Balance)
			}
		})
	}
}

func createAccountAndReturnResult(t testing.TB, testDB *gorm.DB, owner string, balance float64, currency db.Currency) (*db.Account, error) {
	t.Helper()

	account := db.Account{Owner: owner, Balance: balance, Currency: currency}

	result := gorm.WithResult()
	err := gorm.G[db.Account](testDB, result).Create(context.Background(), &account)
	return &account, err
}
