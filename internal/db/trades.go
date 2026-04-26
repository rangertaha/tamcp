package db

import "time"

// Trade is a financial transaction with a symbol, side, quantity, and price.
type Trade struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Source    string    `gorm:"size:32;index" json:"source"`
	Symbol    string    `gorm:"size:32;index" json:"symbol"`
	Side      string    `gorm:"size:8;not null" json:"side"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	UpdatedAt time.Time `json:"updatedAt"`
}
