package migration

import (
	"context"

	models "github.com/onebluesky882/go_fiber_bun_template/internal/models/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

// Migrations เก็บ Go migrations
var Migrations = migrate.NewMigrations()

// RegisterAllModels เพิ่ม Go migrations สำหรับทุกโมเดล
func RegisterAllModels() {
	for _, model := range models.AllModels {
		m := model // capture variable for closure

		Migrations.MustRegister(
			func(ctx context.Context, db *bun.DB) error { // up migration
				_, err := db.NewCreateTable().
					Model(m).
					IfNotExists().
					Exec(ctx)
				return err
			},
			func(ctx context.Context, db *bun.DB) error { // down migration
				_, err := db.NewDropTable().
					Model(m).
					IfExists().
					Exec(ctx)
				return err
			},
		)
	}
}

// New คืนค่า *migrate.Migrations
func New() *migrate.Migrations {
	return Migrations
}
