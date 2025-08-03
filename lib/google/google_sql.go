package google

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lefalya/pageflow"
	"github.com/redis/go-redis/v9"
)

var (
	GOOGLE_NOT_FOUND = errors.New("google not found")
)

type GoogleSQL struct {
	*pageflow.SQLItem `bson:",inline" json:",inline"`
	*GoogleBase       `bson:",inline" json:",inline"`
}

func NewGoogleSQL() *GoogleSQL {
	google := &GoogleSQL{
		GoogleBase: &GoogleBase{},
	}

	pageflow.InitSQLItem(google)
	return google
}

type GoogleManagerSQL struct {
	db         *sqlx.DB
	base       *pageflow.Base[GoogleSQL]
	entityName string
}

func (gsql *GoogleManagerSQL) SetEntityName(entityName string) {
	gsql.entityName = entityName
}

func (gsql *GoogleManagerSQL) Insert(sqlTransaction *sqlx.Tx, google GoogleSQL) (int64, error) {

	var (
		rawSqlFields    = []string{}
		rawSqlDataTypes = []string{}
		rawSqlValues    = []interface{}{}
	)

	queryBuilder := func(sqlField string, sqlValue any) {
		rawSqlFields = append(rawSqlFields, sqlField)
		rawSqlDataTypes = append(rawSqlDataTypes, "?")
		switch v := sqlValue.(type) {
		case []string:
			// Convert string array to PostgreSQL array format
			if len(v) > 0 {
				rawSqlValues = append(rawSqlValues, "{"+strings.Join(v, ",")+"}")
			} else {
				rawSqlValues = append(rawSqlValues, "{}")
			}
		case map[string]bool:
			// Convert map to array of keys where value is true
			var keys []string
			for key, value := range v {
				if value {
					keys = append(keys, key)
				}
			}
			if len(keys) > 0 {
				rawSqlValues = append(rawSqlValues, "{"+strings.Join(keys, ",")+"}")
			} else {
				rawSqlValues = append(rawSqlValues, "{}")
			}
		default:
			rawSqlValues = append(rawSqlValues, v)
		}

	}

	// uuid
	queryBuilder("uuid", google.UUID)

	// randId
	queryBuilder("rand_id", google.RandId)

	// user_uuid
	queryBuilder("user_uuid", google.UserUUID)

	// access_token
	queryBuilder("access_token", google.AccessToken)

	// token_type
	queryBuilder("token_type", google.TokenType)

	// refresh_token
	queryBuilder("refresh_token", google.RefreshToken)

	// expiry
	queryBuilder("expiry", google.Expiry)

	// expires_in
	queryBuilder("expires_in", google.ExpiresIn)

	// raw
	queryBuilder("raw", google.Raw)

	// expiry_delta
	queryBuilder("expiry_delta", google.ExpiryDelta)

	// scopes
	queryBuilder("scopes", google.Scopes)

	// email
	queryBuilder("email", google.Email)

	// created_at
	queryBuilder("created_at", google.CreatedAt)

	// Build the SQL query
	rawSqlFieldsJoin := strings.Join(rawSqlFields, ",")
	rawSqlDataTypesJoin := strings.Join(rawSqlDataTypes, ",")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", gsql.entityName, rawSqlFieldsJoin, rawSqlDataTypesJoin)

	results, err := sqlTransaction.Exec(query, rawSqlValues...)
	if err != nil {
		return 0, err
	}

	lastID, err := results.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (gsql *GoogleManagerSQL) Update(sqlTransaction *sqlx.Tx, google GoogleSQL) error {

	var (
		rawSqlFields = []string{}
		rawSqlValues = []interface{}{}
	)

	queryBuilder := func(sqlField string, sqlValue any) {
		rawSqlFields = append(rawSqlFields, fmt.Sprintf("%s = ?", sqlField))
		rawSqlValues = append(rawSqlValues, sqlValue)

		switch v := sqlValue.(type) {
		case []string:
			// Convert string array to PostgreSQL array format
			if len(v) > 0 {
				rawSqlValues = append(rawSqlValues, fmt.Sprintf("%s = ?", "{"+strings.Join(v, ",")+"}"))
			} else {
				rawSqlValues = append(rawSqlValues, fmt.Sprintf("%s = ?", "{}"))
			}
		case map[string]bool:
			// Convert map to array of keys where value is true
			var keys []string
			for key, value := range v {
				if value {
					keys = append(keys, key)
				}
			}
			if len(keys) > 0 {
				rawSqlValues = append(rawSqlValues, fmt.Sprintf("%s = ?", "{"+strings.Join(keys, ",")+"}"))
			} else {
				rawSqlValues = append(rawSqlValues, fmt.Sprintf("%s = ?", "{}"))
			}
		default:
			rawSqlValues = append(rawSqlValues, v)
		}
	}

	// uuid
	queryBuilder("uuid", google.UUID)

	// randId
	queryBuilder("rand_id", google.RandId)

	// user_uuid
	queryBuilder("user_uuid", google.UserUUID)

	// access_token
	queryBuilder("access_token", google.AccessToken)

	// token_type
	queryBuilder("token_type", google.TokenType)

	// refresh_token
	queryBuilder("refresh_token", google.RefreshToken)

	// expiry
	queryBuilder("expiry", google.Expiry)

	// expires_in
	queryBuilder("expires_in", google.ExpiresIn)

	// raw
	queryBuilder("raw", google.Raw)

	// expiry_delta
	queryBuilder("expiry_delta", google.ExpiryDelta)

	// scopes
	queryBuilder("scopes", google.Scopes)

	// email
	queryBuilder("email", google.Email)

	// updated_at
	queryBuilder("updated_at", google.UpdatedAt)

	rawSqlFieldsJoin := strings.Join(rawSqlFields, ", ")

	query := fmt.Sprintf("UPDATE %s SET (%s) = (%s) WHERE id = ?", gsql.entityName, rawSqlFieldsJoin, rawSqlFieldsJoin)

	_, err := sqlTransaction.Exec(query, rawSqlValues...)
	if err != nil {
		return err
	}

	return nil
}

func (gsql *GoogleManagerSQL) Delete(sqlTransaction *sqlx.Tx, google GoogleSQL) error {
	query := "DELETE FROM " + gsql.entityName + " WHERE uuid = ?"
	_, err := sqlTransaction.Exec(query, google.UUID)
	if err != nil {
		return err
	}
	return nil
}

func (gsql *GoogleManagerSQL) FindByUserUUID(userUUID string) (*GoogleSQL, error) {
	query := fmt.Sprintf(`
		SELECT 
			uuid,
			rand_id,
			user_uuid,
			access_token,
			token_type,
			refresh_token,
			expiry,
			expires_in,
			raw,
			expiry_delta,
			scopes,
			email,
			created_at,
			updated_at
		FROM %s
		WHERE user_uuid = ?`, gsql.entityName)

	var google GoogleSQL
	err := gsql.db.QueryRowx(query, userUUID).StructScan(&google)
	if err != nil {
		return nil, err
	}

	return &google, nil
}

func (gsql *GoogleManagerSQL) SeedByUserUUID(userUUID string) error {
	google, err := gsql.FindByUserUUID(userUUID)
	if err != nil {
		return err
	}

	if google == nil {
		return GOOGLE_NOT_FOUND
	}

	err = gsql.base.Set(*google, userUUID)

	return nil
}

func (gsql *GoogleManagerSQL) FindByUUID(uuid string) (*GoogleSQL, error) {
	query := fmt.Sprintf(`
		SELECT 
			uuid,
			rand_id,
			user_uuid,
			access_token,
			token_type,
			refresh_token,
			expiry,
			expires_in,
			raw,
			expiry_delta,
			scopes,
			email,
			created_at,
			updated_at
		FROM %s
		WHERE uuid = ?`, gsql.entityName)

	var google GoogleSQL
	err := gsql.db.QueryRowx(query, uuid).StructScan(&google)
	if err != nil {
		return nil, err
	}

	return &google, nil
}

func (gsql *GoogleManagerSQL) SeedByUUID(uuid string) error {
	google, err := gsql.FindByUUID(uuid)
	if err != nil {
		return err
	}

	if google == nil {
		return GOOGLE_NOT_FOUND
	}

	errSetAcc := gsql.base.Set(*google, uuid)
	if errSetAcc != nil {
		return errSetAcc
	}
	return nil
}

func (gsql *GoogleManagerSQL) FindByEmail(email string) (*GoogleSQL, error) {
	query := fmt.Sprintf(`
		SELECT 
			uuid,
			rand_id,
			user_uuid,
			access_token,
			token_type,
			refresh_token,
			expiry,
			expires_in,
			raw,
			expiry_delta,
			scopes,
			email,
			created_at,
			updated_at
		FROM %s
		WHERE email = ?`, gsql.entityName)

	var google GoogleSQL
	err := gsql.db.QueryRowx(query, email).StructScan(&google)
	if err != nil {
		return nil, err
	}

	return &google, nil
}

func (gsql *GoogleManagerSQL) SeedByEmail(email string) error {
	google, err := gsql.FindByEmail(email)
	if err != nil {
		return err
	}

	if google == nil {
		return GOOGLE_NOT_FOUND
	}

	errSetAcc := gsql.base.Set(*google, email)
	if errSetAcc != nil {
		return errSetAcc
	}
	return nil
}

func NewGoogleManagerSQL(db *sqlx.DB, redis *redis.Client, entityName string) *GoogleManagerSQL {
	base := pageflow.NewBase[GoogleSQL](redis, entityName+":%s")

	return &GoogleManagerSQL{
		db:         db,
		base:       base,
		entityName: entityName,
	}
}

type GoogleFetchers struct {
	base *pageflow.Base[GoogleSQL]
}

func (gf *GoogleFetchers) FetchByUserUUID(userUUID string) (*GoogleSQL, error) {
	google, err := gf.base.Get(userUUID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return &google, nil
}

func (gf *GoogleFetchers) FetchByUUID(uuid string) (*GoogleSQL, error) {
	google, err := gf.base.Get(uuid)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return &google, nil
}
func (gf *GoogleFetchers) FetchByEmail(email string) (*GoogleSQL, error) {
	google, err := gf.base.Get(email)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return &google, nil
}
func (gf *GoogleFetchers) FetchByRandId(randId string) (*GoogleSQL, error) {
	google, err := gf.base.Get(randId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return &google, nil
}
