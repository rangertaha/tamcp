package db

// Thing is a generic entity with a name and description.
type Thing struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:255;not null" json:"name"`
}
