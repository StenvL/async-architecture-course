package model

import (
	"github.com/guregu/null"
)

type User struct {
	ID    int         `db:"id" json:"id"`
	Name  null.String `db:"name" json:"name"`
	Email string      `db:"email" json:"email"`
	Role  string      `db:"role" json:"role"`
}
