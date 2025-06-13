package models

import "time"

// Product represents a product stored in the database.
type Product struct {
    ID        string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Name      string     `gorm:"type:varchar(255);not null"`
    SKU       string     `gorm:"type:varchar(50);uniqueIndex;not null"`
    Price     float64    `gorm:"not null"`
    Stock     int        `gorm:"not null"`
    CreatedAt time.Time  `gorm:"autoCreateTime"`
    UpdatedAt time.Time  `gorm:"autoUpdateTime"`
    DeletedAt *time.Time `gorm:"index"`
}
