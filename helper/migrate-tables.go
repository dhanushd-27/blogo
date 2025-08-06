package helper

import (
	"log"

	"github.com/dhanushd-27/blog_go/models"
	"gorm.io/gorm"
)

func MigrateTables(db *gorm.DB) {
	if err := models.MigrateUsers(db); err != nil {
		log.Fatal("User migration failed:", err)
	}
	if err := models.MigrateBlogs(db); err != nil {
		log.Fatal("Blog migration failed:", err)
	}
}
