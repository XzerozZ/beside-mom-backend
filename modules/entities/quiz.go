package entities

import "time"

type Quiz struct {
	ID          int       `json:"quiz_id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" gorm:"not null"`
	Question    string    `json:"question" gorm:"not null"`
	Description string    `json:"desc" gorm:"not null"`
	Solution    string    `json:"solution" gorm:"not null"`
	Suggestion  string    `json:"suggestion" gorm:"not null"`
	Banner      string    `json:"banner" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
