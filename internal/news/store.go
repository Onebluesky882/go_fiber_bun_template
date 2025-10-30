package news

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Store struct {
	db bun.IDB
}

func NewStore(db bun.IDB) *Store {
	return &Store{
		db: db,
	}
}

// create news

func (s Store) Create(ctx context.Context, news Record) (createNews Record, err error) {
	err = s.db.NewInsert().Model(&news).Returning("*").Scan(ctx, &createNews)
	if err != nil {
		return createNews, err
	}
	return createNews, err
}

//. find  by id

func (s Store) FindByID(ctx context.Context, id uuid.UUID) (news Record, err error) {
	err = s.db.NewInsert().Model(&news).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return news, err
	}
	return news, err
}

// Find All
func (s Store) FindAll(ctx context.Context) (news []Record, err error) {
	err = s.db.NewInsert().Model(&news).Scan(ctx, &news)
	if err != nil {
		return news, err
	}
	return news, err
}

// delete

func (s Store) DeleteNews(ctx context.Context, id uuid.UUID) error {
	err := s.db.NewDelete().Model(&Record{}).Where("id = ?", id)
	if err != nil {
		return nil
	}
	return nil
}

// update by id

func (s Store) UpdateByID(ctx context.Context, id uuid.UUID, news Record) error {
	_, err := s.db.NewUpdate().Model(&news).Where("id = ?", id).Returning("*").Exec(ctx)
	if err != nil {
		return err
	}
	return err
}
