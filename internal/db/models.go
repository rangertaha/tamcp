package db

import "time"

// Ticker is a tradable instrument loaded by a provider. Composite unique
// index on (provider, symbol) allows the same ticker from multiple sources.
type Ticker struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Source    string    `gorm:"uniqueIndex:idx_ticker;size:32;not null" json:"source"`
	Symbol    string    `gorm:"uniqueIndex:idx_ticker;size:64;not null" json:"symbol"`
	Name      string    `json:"name"`
	Exchange  string    `gorm:"size:32;index" json:"exchange"`
	Active    bool      `gorm:"index" json:"active"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Bar is an OHLCV aggregate for a given instrument and interval.
type Bar struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Source   string    `gorm:"size:32;index" json:"source"`
	Symbol   string    `gorm:"uniqueIndex:idx_bar;size:32;not null" json:"symbol"`
	Interval string    `gorm:"uniqueIndex:idx_bar;size:16;not null" json:"interval"`
	Start    time.Time `gorm:"uniqueIndex:idx_bar;not null" json:"start"`
	Open     float64   `json:"open"`
	High     float64   `json:"high"`
	Low      float64   `json:"low"`
	Close    float64   `json:"close"`
	Volume   float64   `json:"volume"`
}

// Order is a trading order recorded for analysis.
type Order struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Source      string    `gorm:"size:32;index" json:"source"`
	OrderID     string    `gorm:"uniqueIndex;size:64;not null" json:"orderId"`
	Symbol      string    `gorm:"size:32;index;not null" json:"symbol"`
	Side        string    `gorm:"size:8;not null" json:"side"`
	Type        string    `gorm:"size:16;not null" json:"type"`
	Qty         float64   `json:"qty"`
	Price       float64   `json:"price"`
	Status      string    `gorm:"size:16;index" json:"status"`
	SubmittedAt time.Time `json:"submittedAt"`
}
