package lib

import (
	"database/sql"
	"github.com/lefalya/pageflow"
	"github.com/redis/go-redis/v9"
)

type ResetPasswordManagerSQL struct {
	base       *pageflow.Base[AccountSQL]
	db         *sql.DB
	entityName string
}

func (ar *ResetPasswordManagerSQL) FindRequest(email string) (*AccountSQL, error) {
	query := "SELECT * FROM " + ar.entityName + "ResetPassword WHERE email = $1"
	row := ar.db.QueryRow(query, email)

}

func (ar *ResetPasswordManagerSQL) CreateRequest(email string) error {
	return nil
}

func (ar *ResetPasswordManagerSQL) Reset(email string, password string) error {
	return nil
}

func NewResetPasswordSQL(db *sql.DB, redis *redis.Client, entityName string) *ResetPasswordManagerSQL {
	base := pageflow.NewBase[AccountSQL](redis, entityName+":%s")
	return &ResetPasswordManagerSQL{
		base: base,
		db:   db,
	}
}
