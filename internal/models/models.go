package models

import "github.com/uptrace/bun"

// All Models
var AllModels = []any{
	(*User)(nil),
	(*Comment)(nil),
	(*Student)(nil),
	(*Job)(nil),
	(*Car)(nil),
}

// table : users
type User struct {
	bun.BaseModel `bun:"table:user"`

	ID       int64      `bun:",pk,autoincrement"`
	Name     string     `bun:",notnull"`
	Comments []*Comment `bun:"rel:has-many,join:id=user_id"`
}

type Comment struct {
	bun.BaseModel `bun:"table:comments"`

	ID      int64  `bun:",pk,autoincrement"`
	UserID  int64  `bun:",notnull"`
	Content string `bun:",notnull"`

	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

type Student struct {
	bun.BaseModel `bun:"table:student"`

	ID      int64  `bun:",pk,autoincrement"`
	UserID  int64  `bun:",notnull"`
	Content string `bun:",notnull"`

	User *User `bun:"rel:belongs-to,join:user_id=id"`
}
type Car struct {
	bun.BaseModel `bun:"table:job"`

	ID      int64  `bun:",pk,autoincrement"`
	UserID  int64  `bun:",notnull"`
	Content string `bun:",notnull"`

	User *User `bun:"rel:belongs-to,join:user_id=id"`
}
type Job struct {
	bun.BaseModel `bun:"table:car"`

	ID      int64  `bun:",pk,autoincrement"`
	UserID  int64  `bun:",notnull"`
	Content string `bun:",notnull"`

	User *User `bun:"rel:belongs-to,join:user_id=id"`
}
