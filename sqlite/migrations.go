package sqlite

import (
	"fmt"

	"src.sqlkite.com/utils/log"
)

type Migrate func(conn Conn) error

type Migration struct {
	Version uint16
	Migrate Migrate
}

func MigrateAll(conn Conn, migrations []Migration) error {
	latestVersion, err := GetCurrentMigrationVersion(conn)
	if err != nil {
		return err
	}

	log.Info("migration_check_start").String("storage", "sqlite").Int("installed_version", latestVersion).Log()
	for _, migration := range migrations {
		version := int(migration.Version)
		if version <= latestVersion {
			continue
		}

		err := conn.Transaction(func() error {
			if err := migration.Migrate(conn); err != nil {
				return fmt.Errorf("Failed to run sqlite migration #%d - %w", version, err)
			}

			return conn.Exec(`
				insert into sqlkite_migrations (version) values (?1)
			`, version)
		})

		if err != nil {
			log.Error("migration_fail").Int("version", version).Err(err).Log()
			return err
		}
		log.Info("migration_applied").Int("version", version).Log()
	}
	log.Info("migration_check_end").Log()

	return nil
}

func GetCurrentMigrationVersion(conn Conn) (int, error) {
	exists, err := conn.TableExists("sqlkite_migrations")
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, conn.Exec(`
			create table sqlkite_migrations (
				version integer not null
			)
		`)
	}

	return Scalar[int](conn, `
		select max(version)
		from sqlkite_migrations
	`)
}
