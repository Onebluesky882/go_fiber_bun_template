package user

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:user"`

	ID   int16  `bun:"id,pk,autoincrement"`
	Name string `bun:"name"`
}
