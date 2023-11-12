package persistent

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	log.Logger = log.With().Caller().Logger()

	dsn := os.Getenv("DB_DSN")
	//TODO: remove this: creds leaking
	log.Debug().Str("dsn", dsn).Msg("connecting to database")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return &Database{
		db: db,
	}, nil
}

func (s *Database) Close() error {
	log.Debug().Msg("closing database connection")
	return s.db.Close()
}

func (s *Database) Migrate() error {
	s.db.SetMaxOpenConns(25)
	s.db.SetMaxIdleConns(25)
	s.db.SetConnMaxLifetime(5 * time.Minute)
	driver, err := mysql.WithInstance(s.db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"mysql", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func (s *Database) DB() *sql.DB {
	return s.db
}

func (s *Database) Ping() error {
	return s.db.Ping()
}
