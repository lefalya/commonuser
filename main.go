package commonuser

import (
	"encoding/json"
	"github.com/lefalya/item"
	"github.com/lefalya/pageflow"
	"time"
)

// Original author: Muammar Zikri Aksana
// The following Account and AssociatedAccount types and their JSON marshaling/unmarshaling
// implementations were originally written by Muammar Zikri Aksana.
type AssociatedAccount struct {
	Name     string `json:"name,omitempty" db:"-"`
	Email    string `json:"email,omitempty" db:"-"`
	Uuid     string `json:"uuid,omitempty" db:"-"`
	Sub      string `json:"sub,omitempty" db:"-"`
	Provider string `json:"provider,omitempty" db:"-"`
}

type Base struct {
	Sub               string              `json:"sub,omitempty" db:"-"`
	Name              string              `json:"name,omitempty" db:"name"`
	Username          string              `json:"username,omitempty" db:"username"`
	Password          string              `json:"-" db:"password"`
	Email             string              `json:"email,omitempty" db:"email"`
	Avatar            string              `json:"avatar,omitempty" db:"avatar"`
	AssociatedAccount []AssociatedAccount `json:"associatedAccount,omitempty" db:"-"`
	Suspended         bool                `json:"suspended,omitempty" db:"suspended"`
}

type Account struct {
	*item.Foundation
	Base
}

func (c *Account) UnmarshalJSON(data []byte) error {
	type Alias Account
	aux := &struct {
		CreatedAt *int64 `json:"createdAt,omitempty"`
		UpdatedAt *int64 `json:"updatedAt,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if aux.CreatedAt != nil {
		c.CreatedAt = time.Unix(0, *aux.CreatedAt*int64(time.Millisecond))
	}
	if aux.UpdatedAt != nil {
		c.UpdatedAt = time.Unix(0, *aux.UpdatedAt*int64(time.Millisecond))
	}

	return nil
}

func (c Account) MarshalJSON() ([]byte, error) {
	type Alias Account
	output := struct {
		Alias
		CreatedAt *int64 `json:"createdAt,omitempty"`
		UpdatedAt *int64 `json:"updatedAt,omitempty"`
	}{
		Alias: (Alias)(c),
	}

	if !c.CreatedAt.IsZero() {
		createdAt := c.CreatedAt.UnixNano() / int64(time.Millisecond)
		output.CreatedAt = &createdAt
	}
	if !c.UpdatedAt.IsZero() {
		updatedAt := c.UpdatedAt.UnixNano() / int64(time.Millisecond)
		output.UpdatedAt = &updatedAt
	}

	return json.Marshal(&output)
}

type AccountMongo struct {
	*pageflow.MongoItem `bson:",inline" json:",inline"`
	*Base               `bson:",inline" json:",inline"`
}

// TODO: implement AccountSQL
