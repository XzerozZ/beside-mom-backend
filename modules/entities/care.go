package entities

import "time"

type Care struct {
	ID          string    `json:"c_id" gorm:"primaryKey"`
	Type        string    `json:"type" gorm:"not null"`
	Title       string    `json:"title" gorm:"title"`
	Description string    `json:"desc"`
	Banner      string    `json:"banner" gorm:"not null"`
	UserID      string    `json:"user_id" gorm:"not null"`
	Assets      []Asset   `json:"assets" gorm:"many2many:care_assets;"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
