package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/onebluesky882/go_fiber_bun_template/internal/database"
	"github.com/onebluesky882/go_fiber_bun_template/internal/migration"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
)

func main() {
	dbService := database.New()
	db := dbService.GetDB()

	// เปิด bundebug เพื่อ log SQL
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	app := &cli.App{
		Name:  "migrate",
		Usage: "database migrations",
		Commands: []*cli.Command{
			newMigrationCmd(
				migrate.NewMigrator(db, migration.New(), migrate.WithMarkAppliedOnSuccess(true)),
			),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func newMigrationCmd(m *migrate.Migrator) *cli.Command {
	return &cli.Command{
		Name:  "migrate",
		Usage: "database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migrations table",
				Action: func(ctx *cli.Context) error {
					return m.Init(ctx.Context)
				},
			},
			{
				Name:  "up",
				Usage: "run up migrations",
				Action: func(ctx *cli.Context) error {
					if err := m.Lock(ctx.Context); err != nil {
						return err
					}
					defer m.Unlock(ctx.Context)

					group, err := m.Migrate(ctx.Context)
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
				Action: func(ctx *cli.Context) error {
					if err := m.Lock(ctx.Context); err != nil {
						return err
					}
					defer m.Unlock(ctx.Context)

					group, err := m.Rollback(ctx.Context)
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
				Action: func(ctx *cli.Context) error {
					name := strings.Join(ctx.Args().Slice(), "_")
					if name == "" {
						return fmt.Errorf("please provide a migration name")
					}

					files, err := m.CreateTxSQLMigrations(ctx.Context, name)
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
				Action: func(ctx *cli.Context) error {
					ms, err := m.MigrationsWithStatus(ctx.Context)
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
