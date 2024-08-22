package main 
import (
	"log"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/logger"
)
var DB *gorm.DB

func ConnectDB() error{
	// connect to db
	DB, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	DB.Logger = logger.Default.LogMode(logger.Info)
	
	// apply migrations 
	log.Println("applying migrations ...")
	err = DB.AutoMigrate(&Note{})
	if err != nil {
		return err
	}
	log.Println("🚀️ connected successfully to the database ..")
	return nil 
}