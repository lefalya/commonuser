package google_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/lefalya/commonuser/lib/google"
	"github.com/lefalya/pageflow"

	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

var (
	db   *sqlx.DB
	mock sqlxmock.Sqlmock
	err  error
)

func init() {

	db, mock, err = sqlxmock.Newx()
	if err != nil {
		panic(fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	}
}

func TestCreate(t *testing.T) {

	googleManager := google.NewGoogleManagerSQL(db, nil, "google")

	googleAccountData := google.GoogleSQL{}

	pageflow.InitSQLItem(&googleAccountData)

	googleAccountData.GoogleBase.SetUserUUID("test-user-uuid")
	googleAccountData.GoogleBase.SetAccessToken("test-access-token")
	googleAccountData.GoogleBase.SetTokenType("test-token-type")
	googleAccountData.GoogleBase.SetRefreshToken("test-refresh-token")
	googleAccountData.GoogleBase.SetExpiry(nil)
	googleAccountData.GoogleBase.SetExpiresIn(3600)
	googleAccountData.GoogleBase.SetRaw("test-raw-data")
	googleAccountData.GoogleBase.SetExpiryDelta(0)
	googleAccountData.GoogleBase.SetScopes(map[string]bool{"scope1": true, "scope2": true})
	googleAccountData.GoogleBase.SetEmail("test@example.com")

	// Set up mock expectations for database transaction
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO google").WillReturnResult(sqlxmock.NewResult(1, 1))
	mock.ExpectCommit()

	sqlTransaction, err := db.Beginx()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	id, err := googleManager.Insert(sqlTransaction, googleAccountData)
	if err != nil {
		t.Fatalf("Failed to insert Google account: %v", err)
	}

	// Commit the transaction
	err = sqlTransaction.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	// expectedID := int64(1)
	if id != 1 {
		t.Fatalf("Expected ID to be 1, got %d", id)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	googleAccountDataJson, _ := json.Marshal(googleAccountData)
	fmt.Println("Google Account Data:", string(googleAccountDataJson))
}
