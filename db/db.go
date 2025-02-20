package db

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"microblogging/internal/model"

	"microblogging/internal/config"
)

type Database struct {
	*gorm.DB
}

var (
	once     sync.Once
	instance *Database
)

// New creates a new database connection. This method will panic if
// there is an error connecting to the database. This is intentional because
// the application cannot run without a database connection. This method functions
// as a singleton, so it will only create a connection once.
// If you need the connection at runtime, use GetDB instead.
func Init(env *config.Config) *Database {
	createInstance := func() {
		instance = postgresDB(env)
		instance.Migrate()
	}

	once.Do(createInstance)

	return instance
}

func postgresDB(env *config.Config) *Database {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		env.PostgresHost,
		env.PgUser,
		env.PgPassword,
		env.PgDatabase,
		env.PostgresPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &Database{db}
}

// Will migrate all models of the application to the database. This method
// should be called only once, when the application starts.
func (d *Database) Migrate() {
	fmt.Println("Migrating models")
	models := []interface{}{
		&model.User{},
		&model.Tweet{},
		&model.Follow{},
	}

	for _, model := range models {
		if err := d.AutoMigrate(model); err != nil {
			fmt.Printf("Error migrating model %T: %v\n", model, err)
		}
	}
	fmt.Println("Migrating models: DONE")
}

func (d Database) Close() {
	sqlDB, err := d.DB.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.Close()
}
