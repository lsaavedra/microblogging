package testhelpers

import (
	"os/exec"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"microblogging/db"
)

var (
	once     sync.Once
	instance *db.Database
)

func NewDbForTest() *db.Database {
	testDbName := "test_db"
	ClearTestDB(testDbName)
	once.Do(func() {
		gormDB, err := gorm.Open(sqlite.Open(testDbName), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		instance = &db.Database{DB: gormDB}
		instance.Migrate()
	})

	return instance
}

func ClearTestDB(testDbName string) {
	exec.Command("rm", "-f", testDbName)
}
