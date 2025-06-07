package models

type Blog struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Content   string `gorm:"not null"`
	CreatedBy User   `gorm:"foreignKey:CreatedBy"`
}
