package user

import (
	"context"

	"github.com/onebluesky882/go_fiber_bun_template/internal/models/sql"
	"github.com/uptrace/bun"
)

// sql User
type User struct {
	sql.User
}

type Repository struct {
	db *bun.DB
}

// constructor
func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context) ([]User, error) {
	var user []User
	err := r.db.NewSelect().Model(&user).Scan(ctx)
	return user, err
}

func (r *Repository) GetUserById(ctx context.Context, id int) (User, error) {
	var user User
	err := r.db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	return user, err
}
