package pgdb

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// ConnectDatabase returns the database instance
// After connecting to postgres with required envionment variables.
// Also AutoMigrate the Database on First Run
func ConnectDatabase(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	return db
}

func CloseConnection(db *gorm.DB) {
	sqlDB, _ := db.DB()
	sqlDB.Close()
}
