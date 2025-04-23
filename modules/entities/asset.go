package entities

type Asset struct {
	ID   string `json:"asset_id" gorm:"primaryKey"`
	Link string `json:"link" gorm:"not null"`
}
