package model

import (
	"database/sql"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Post struct {
	Id      int
	User_id int
	Content string
	Author  string
	Created_at    sql.NullTime
}

func (p *Post) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Content, validation.NilOrNotEmpty, validation.Length(1, 300)),
	)
}
