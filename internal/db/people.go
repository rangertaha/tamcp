package db

// Geography is a geographical location with a name and coordinates.
type Person struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:255;not null" json:"name"`
}
