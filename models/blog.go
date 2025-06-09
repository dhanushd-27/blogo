package models

import "gorm.io/gorm"

type Blog struct {
	ID      uint   `gorm:"primaryKey"`
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
	User    User   `gorm:"foreignKey:Email"`
	Email   string
}

func MigrateBlogs(db *gorm.DB) error {
	err := db.AutoMigrate(&Blog{})
	return err
}
