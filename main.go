package commonuser

import (
	"database/sql"
	"github.com/lefalya/commonuser/lib"
	"github.com/redis/go-redis/v9"
)

func NewAccountManagerSQL(db *sql.DB, redis *redis.Client, entityName string) *lib.AccountManagerSQL {
	return lib.NewAccountManagerSQL(db, redis, entityName)
}

func NewUpdateEmailManagerSQL(db *sql.DB, entityName string) *lib.UpdateEmailManagerSQL {
	return lib.NewUpdateEmailManagerSQL(db, entityName)
}

func NewResetPasswordSQL(db *sql.DB, redis *redis.Client, entityName string) *lib.ResetPasswordManagerSQL {
	return lib.NewResetPasswordSQL(db, redis, entityName)
}

func NewJWTHandler(jwtSecret string, jwtTokenIssuer string, jwtTokenLifeSpan int) *lib.JWTHandler {
	return lib.NewJWTHandler(jwtSecret, jwtTokenIssuer, jwtTokenLifeSpan)
}

func CreateResetPasswordTableSQL(db *sql.DB, entityName string) error {
	tableName := entityName + "ResetPassword"
	query := `CREATE TABLE IF NOT EXISTS ` + tableName + ` (
		email VARCHAR(255) UNIQUE NOT NULL,
		uuid VARCHAR(255) UNIQUE,
		accountuuid VARCHAR(255) UNIQUE,
		createdat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		token VARCHAR(255) UNIQUE
	)`

	_, err := db.Exec(query)
	return err
}

func CreateAccountTableSQL(db *sql.DB, entityName string) error {
	query := `CREATE TABLE IF NOT EXISTS ` + entityName + ` (
		uuid VARCHAR(255) UNIQUE NOT NULL,
		randId VARCHAR(255) UNIQUE,
		createdat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		name VARCHAR(255),
		username VARCHAR(255) UNIQUE,
		password VARCHAR(255),
		email VARCHAR(255) UNIQUE,
		avatar VARCHAR(255),
		suspended BOOLEAN DEFAULT FALSE
	)`

	_, err := db.Exec(query)
	return err
}

func CreateUpdateEmailTableSQL(db *sql.DB, entityName string) error {
	tableName := entityName + "UpdateEmail"
	query := `CREATE TABLE IF NOT EXISTS ` + tableName + ` (
		uuid VARCHAR(255) UNIQUE NOT NULL,
		randId VARCHAR(255) UNIQUE,
		createdat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updatedat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		accountuuid VARCHAR(255) UNIQUE,
		previousemailaddress VARCHAR(255),
		newemailaddress VARCHAR(255) UNIQUE,
		resettoken VARCHAR(255)
	)`

	_, err := db.Exec(query)
	return err
}
