package sqlite

import "fmt"

type Migrate func(conn Conn) error

type Migration struct {
	Version uint16
	Migrate Migrate
}

func MigrateAll(conn Conn, migrations []Migration) error {
	latestVersion, err := getCurrentVersion(conn)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		version := int(migration.Version)
		if version <= latestVersion {
			continue
		}

		err := conn.Transaction(func() error {
			if err := migration.Migrate(conn); err != nil {
				return fmt.Errorf("Failed to run migration #%d - %w", version, err)
			}

			return conn.Exec(`
				insert into gobl_migrations (version) values (?1)
			`, version)
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func getCurrentVersion(conn Conn) (int, error) {
	exists, err := conn.TableExists("gobl_migrations")
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, conn.Exec(`
			create table gobl_migrations (
				version integer not null
			)
		`)
	}

	return Scalar[int](conn, `
		select max(version)
		from gobl_migrations
	`)
}
