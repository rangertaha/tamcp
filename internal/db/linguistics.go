package db

import "time"

// Linguistics is a linguistic entity with a name and description.
type Language struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"size:255;not null" json:"description"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
