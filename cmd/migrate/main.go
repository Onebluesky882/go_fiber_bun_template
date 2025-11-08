package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/onebluesky882/go_fiber_bun_template/internal/database"
	"github.com/onebluesky882/go_fiber_bun_template/internal/migration"
	"github.com/onebluesky882/go_fiber_bun_template/internal/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
)

func main() {
	ctx := context.Background()
	db := database.New().GetDB()

	// เปิด debug log สำหรับ SQL
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	// ✅ สร้าง SQL ไฟล์อัตโนมัติ
	// generateSQL(db)

	// ✅ สร้างตารางจริงใน DB
	createAllTables(ctx, db)

	// ✅ CLI app สำหรับ migrations
	m := migrate.NewMigrator(db, migration.New(), migrate.WithMarkAppliedOnSuccess(true))

	app := &cli.App{
		Name:  "migrate",
		Usage: "database migrations",
		Commands: []*cli.Command{
			newMigrationCmd(m),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// Dynamic table creation
func createAllTables(ctx context.Context, db *bun.DB) {
	for _, model := range models.AllModels {
		if _, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx); err != nil {
			log.Fatalf("❌ Create table failed: %v", err)
		}
	}
	fmt.Println("✅ All tables created successfully.")
}

// CLI command
func newMigrationCmd(m *migrate.Migrator) *cli.Command {
	return &cli.Command{
		Name:  "migrate",
		Usage: "database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migrations table",
				Action: func(c *cli.Context) error {
					return m.Init(c.Context)
				},
			},
			{
				Name:  "up",
				Usage: "run up migrations",
				Action: func(c *cli.Context) error {
					if err := m.Lock(c.Context); err != nil {
						return err
					}
					defer m.Unlock(c.Context)

					group, err := m.Migrate(c.Context)
					if err != nil {
						return err
					}
					if group.IsZero() {
						fmt.Println("database is up-to-date")
					} else {
						fmt.Printf("migrated to %s\n", group)
					}
					return nil
				},
			},
			{
				Name:  "down",
				Usage: "rollback last migration group",
				Action: func(c *cli.Context) error {
					if err := m.Lock(c.Context); err != nil {
						return err
					}
					defer m.Unlock(c.Context)

					group, err := m.Rollback(c.Context)
					if err != nil {
						return err
					}
					if group.IsZero() {
						fmt.Println("there are no groups to rollback")
					} else {
						fmt.Printf("rolled back to %s\n", group)
					}
					return nil
				},
			},
			{
				Name:  "create",
				Usage: "create up and down SQL migrations",
				Action: func(c *cli.Context) error {
					name := strings.Join(c.Args().Slice(), "_")
					if name == "" {
						return fmt.Errorf("please provide a migration name")
					}

					files, err := m.CreateTxSQLMigrations(c.Context, name)
					if err != nil {
						return err
					}

					for _, f := range files {
						fmt.Printf("created migration: %s (%s)\n", f.Name, f.Path)
					}
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migration status",
				Action: func(c *cli.Context) error {
					ms, err := m.MigrationsWithStatus(c.Context)
					if err != nil {
						return err
					}
					fmt.Println("All migrations:", ms)
					fmt.Println("Unapplied migrations:", ms.Unapplied())
					fmt.Println("Last migration group:", ms.LastGroup())
					return nil
				},
			},
		},
	}
}

// generate .up.sql และ .down.sql อัตโนมัติ
func getTableName(m any) string {
	t := reflect.TypeOf(m)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	tableName := t.Name() // fallback

	if baseField, ok := t.FieldByName("BaseModel"); ok {
		tag := baseField.Tag.Get("bun")
		for _, part := range strings.Split(tag, ",") {
			if strings.HasPrefix(part, "table:") {
				tableName = part[len("table:"):]
			}
		}
	}

	return tableName
}
