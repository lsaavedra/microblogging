package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMockPostgresConnection(t *testing.T) (*Database, sqlmock.Sqlmock) {
	sqlMockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cannot create sqlmock: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock",
		DriverName:           "postgres",
		Conn:                 sqlMockDB,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector)
	if err != nil {
		t.Fatalf("cannot open gorm connection: %v", err)
	}

	conn := Database{
		gormDB,
	}

	return &conn, mock
}
