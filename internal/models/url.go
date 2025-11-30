package models

// URL represents a shortened URL mapping in the database
type URL struct {
	ID            uint   `gorm:"primaryKey"`
	ShortCode     string `gorm:"uniqueIndex;not null;size:20"`
	OriginalURL   string `gorm:"not null;size:2048"`
	NormalizedURL string `gorm:"uniqueIndex;not null;size:2048"`
}

// TableName specifies the table name for the URL model
func (URL) TableName() string {
	return "urls"
}
