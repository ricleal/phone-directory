package persistent

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
	_ "gorm.io/driver/mysql"
)

// Storage is struct that holds the database connection.
type Storage struct {
	db *Database
}

// NewStorage returns a new Handler with a database connection.
func NewStorage() (*Storage, error) {

	db, err := NewDatabase()
	if err != nil {
		return nil, fmt.Errorf("error creating database: %w", err)
	}
	if err := db.Migrate(); err != nil {
		return nil, fmt.Errorf("error migrating database: %w", err)
	}
	log.Info().Msg("migrations ran successfully")
	return &Storage{
		db: db,
	}, nil
}

// Close closes the database connection.
func (s *Storage) Close() error {
	return s.db.Close()
}

// DB returns the database connection.
func (s *Storage) DB() *sql.DB {
	return s.db.DB()
}

func (s *Storage) Ping() error {
	return s.db.Ping()
}
