package model

import "database/sql"

type Admin struct {
	ID       int `json:"id"`
	Admin_id int `json:"admin_id"`
	Email string `json:"email"`
	Date     sql.NullTime
}
