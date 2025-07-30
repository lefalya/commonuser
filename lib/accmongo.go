package lib

import "github.com/lefalya/pageflow"

type AccountMongo struct {
	*pageflow.MongoItem `bson:",inline" json:",inline"`
	*Base               `bson:",inline" json:",inline"`
}
