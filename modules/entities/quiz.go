package entities

type Quiz struct {
	ID   string `json:"asset_id" gorm:"primaryKey"`
	Link string `json:"link" gorm:"not null"`
}