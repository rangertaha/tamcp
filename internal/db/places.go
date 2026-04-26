package db

import "time"

// Geography is a geographical location with a name and coordinates.
type Place struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Latitude  float64   `gorm:"not null" json:"latitude"`
	Longitude float64   `gorm:"not null" json:"longitude"`
	UpdatedAt time.Time `json:"updatedAt"`
}
