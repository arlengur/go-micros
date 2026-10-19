package migrator

import (
	"database/sql"
	"log"

	migrate "github.com/rubenv/sql-migrate"
	_ "modernc.org/sqlite"
)

type migrator struct {
	db         *sql.DB
	migrations *migrate.FileMigrationSource
}

func NewMigrator() *migrator {
	db, err := sql.Open("sqlite", "./tasks.db")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return &migrator{
		db: db,
		migrations: &migrate.FileMigrationSource{
			Dir: "migrations",
		},
	}
}

func (m *migrator) Up() error {
	count, err := migrate.Exec(m.db, "sqlite3", m.migrations, migrate.Up)
	if err != nil {
		return err
	}
	log.Printf("Applied %d migrations", count)

	return nil
}
