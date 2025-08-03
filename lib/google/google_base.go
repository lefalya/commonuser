package google

import (
	"time"

	"github.com/lefalya/item"
)

type GoogleBase struct {
	UserUUID     string          `json:"user_uuid" bson:"user_uuid" db:"user_uuid"`
	AccessToken  string          `json:"access_token" bson:"access_token" db:"access_token"`
	TokenType    string          `json:"token_type" bson:"token_type" db:"token_type"`
	RefreshToken string          `json:"refresh_token" bson:"refresh_token" db:"refresh_token"`
	Expiry       *time.Time      `json:"expiry" bson:"expiry" db:"expiry"`
	ExpiresIn    int64           `json:"expires_in" bson:"expires_in" db:"expires_in"`
	Raw          string          `json:"raw" bson:"raw" db:"raw"`
	ExpiryDelta  time.Duration   `json:"expiry_delta" bson:"expiry_delta" db:"expiry_delta"`
	Scopes       map[string]bool `json:"scopes" bson:"scopes"`
	Email        string          `json:"email" bson:"email"`
}

func (b *GoogleBase) SetUserUUID(uuid string) {
	b.UserUUID = uuid
}

func (b *GoogleBase) SetAccessToken(token string) {
	b.AccessToken = token
}

func (b *GoogleBase) SetTokenType(tokenType string) {
	b.TokenType = tokenType
}

func (b *GoogleBase) SetRefreshToken(token string) {
	b.RefreshToken = token
}

func (b *GoogleBase) SetExpiry(expiry *time.Time) {
	b.Expiry = expiry
}

func (b *GoogleBase) SetExpiresIn(expiresIn int64) {
	b.ExpiresIn = expiresIn
}

func (b *GoogleBase) SetRaw(raw string) {
	b.Raw = raw
}

func (b *GoogleBase) SetExpiryDelta(expiryDelta time.Duration) {
	b.ExpiryDelta = expiryDelta
}

func (b *GoogleBase) SetScopes(scopes map[string]bool) {
	b.Scopes = scopes
}

func (b *GoogleBase) SetEmail(email string) {
	b.Email = email
}

type Google struct {
	*item.Foundation
	GoogleBase
}
