package migration

import (
	"github.com/uptrace/bun/migrate"
)

var migrations = migrate.NewMigrations()

func New() *migrate.Migrations {
	return migrations
}

func init() {
	// Discover Go migrations (optional)
	migrations.DiscoverCaller()
}
