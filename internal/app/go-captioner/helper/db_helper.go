package helper

import (
	"github.com/sam-laister/tiktok-creator/ent"
)

func GetDB() (*ent.Client, error) {
	client, err := ent.Open("sqlite3", "file:app.db?cache=shared&_fk=1")
	if err != nil {
		return nil, err
	}
	return client, nil
}
