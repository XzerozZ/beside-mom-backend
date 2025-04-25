package entities

import "time"

type Growth struct {
	ID        string    `json:"G_id" gorm:"primaryKey"`
	Length    float64   `json:"length" gorm:"not null"`
	Weight    float64   `json:"weight" gorm:"not null"`
	Months    int       `json:"months" gorm:"not null"`
	KidID     string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
