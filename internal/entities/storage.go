package entities

import (
	"github.com/jmoiron/sqlx"
)

type Link struct {
	Id           int    `json:"id,omitempty" db:"id"`
	ShortLink    string `json:"short_link,omitempty" db:"short_link"`
	OriginalLink string `json:"original_link" db:"original_link"`
}

var (
	Db        *sqlx.DB
	TableName = "links"
)
