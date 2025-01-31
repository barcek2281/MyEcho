package model

import (
	"database/sql"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Post struct {
	Id      int    `json:"id"`
	User_id int    `json:"user_id"`
	Content string `json:"content"`
	Author  string `json:"author"` // такой колонки нету в модели

	Created_at sql.NullTime
}

func (p *Post) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Content, validation.NilOrNotEmpty, validation.Length(1, 300)),
	)
}

func (p *Post) ConverDateToString() string {
	return p.Created_at.Time.String()
}
