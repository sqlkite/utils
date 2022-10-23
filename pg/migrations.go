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

func MigrateAll(db DB, appName string, migrations []Migration) error {
	latestVersion, err := GetCurrentMigrationVersion(db, appName)
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

			_, err = tx.Exec(context.Background(), `insert into gobl_migrations (app, version) values ($1, $2)`, appName, version)

			if err != nil {
				return fmt.Errorf("pg insert into gobl_migrations - %w", err)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func GetCurrentMigrationVersion(db DB, appName string) (int, error) {
	exists, err := db.TableExists("gobl_migrations")
	if err != nil {
		return 0, err
	}

	if !exists {
		_, err := db.Exec(context.Background(), `
			create table gobl_migrations (
				app text not null,
				version integer not null,
				created timestamptz not null default now(),
				primary key(app, version)
			)
		`)
		if err != nil {
			return 0, fmt.Errorf("pg create gobl_migrations - %w", err)
		}
		return 0, nil
	}

	value, err := Scalar[*int](db, `
		select max(version)
		from gobl_migrations
		where app = $1
	`, appName)

	if err != nil {
		return 0, fmt.Errorf("pg max migration - %w", err)
	}
	if value == nil {
		return 0, nil
	}
	return *value, nil
}
