package persistent

import (
	"database/sql"
	"fmt"
	"log"
	"phone-directory/internal/repository"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func standardToGormConnection(sqlDB repository.DBTx) (*gorm.DB, error) {

	log.Println("Converting sql.DB to gorm.DB")
	db, ok := sqlDB.(*sql.DB)
	if !ok {
		return nil, fmt.Errorf("sqlDB is not a *sql.DB")
	}
	var err error
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm connection: %w", err)
	}
	return gormDB, nil
}
