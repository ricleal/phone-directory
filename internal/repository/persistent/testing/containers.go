package testing

import (
	"context"
	"fmt"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func initContainer(ctx context.Context) (*mysql.MySQLContainer, error) {
	dbname := os.Getenv("MYSQL_DATABASE")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	initScript := os.Getenv("SQL_INIT_FILE")

	container, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:5.7"),
		mysql.WithDatabase(dbname),
		mysql.WithUsername(user),
		mysql.WithPassword(password),
		mysql.WithScripts(initScript),
		// testcontainers.WithWaitStrategy(
		// 	wait.ForLog("MySQL Community Server (GPL)").WithOccurrence(1),
		// ),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to run mysql container: %w", err)
	}

	connectionString, err := container.ConnectionString(ctx, "charset=utf8&parseTime=True&loc=Local&multiStatements=true")
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	// os.Setenv("DB_CONNECTION_STRING", connectionString)
	os.Setenv("DB_DSN", connectionString)
	log.Debug().Str("connection_string", connectionString).Msg("mysql test container running")
	return container, nil
}

// SetupContainer sets up a mysql container and runs migrations.
func SetupContainer(ctx context.Context) (*mysql.MySQLContainer, error) {
	container, err := initContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup container: %w", err)
	}
	return container, nil
}

// TeardownContainer terminates the mysql container.
func TeardownContainer(ctx context.Context, container *mysql.MySQLContainer) error {
	if container != nil {
		if err := container.Terminate(ctx); err != nil {
			return fmt.Errorf("failed to terminate container: %w", err)
		}
	}
	return nil
}
