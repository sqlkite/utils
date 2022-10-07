package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Migrate func(tx pgx.Tx) error

type Migration struct {
	Version uint16
	Migrate Migrate
}

func MigrateAll(db DB, migrations []Migration) error {
	latestVersion, err := getCurrentVersion(db)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		version := int(migration.Version)
		if version <= latestVersion {
			continue
		}

		err := db.Transaction(func(tx pgx.Tx) error {
			if err := migration.Migrate(tx); err != nil {
				return fmt.Errorf("Failed to run pg migration #%d - %w", version, err)
			}

			_, err = tx.Exec(context.Background(), `
				insert into gobl_migrations (version) values ($1)
			`, version)

			return err
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func getCurrentVersion(db DB) (int, error) {
	exists, err := db.TableExists("gobl_migrations")
	if err != nil {
		return 0, err
	}

	if !exists {
		_, err := db.Exec(context.Background(), `
			create table gobl_migrations (
				version integer not null
			)
		`)
		return 0, err
	}

	return Scalar[int](db, `
		select max(version)
		from gobl_migrations
	`)
}
