package user

import (
	"context"

	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

// constructor
func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.NewSelect().Model(&users).Scan(ctx)
	return users, err
}
