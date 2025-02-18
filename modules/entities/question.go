package entities

import "time"

type Question struct {
	ID        string    `json:"Q_id" gorm:"primaryKey"`
	Question  string    `json:"question" gorm:"not null"`
	Answer    string    `json:"answer" gorm:"not null"`
	UserID    string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
