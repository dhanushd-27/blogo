package helper

import (
	"github.com/dhanushd-27/blog_go/models"
	"gorm.io/gorm"
)

func MigrateTables(db *gorm.DB) {
	models.MigrateUsers(db)
	models.MigrateBlogs(db)
}
