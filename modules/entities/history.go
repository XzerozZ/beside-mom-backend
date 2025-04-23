package entities

import "time"

type History struct {
	ID             string    `json:"H_id" gorm:"primaryKey"`
	QuizID         int       `json:"quiz_id" gorm:"not null"`
	Answer         bool      `json:"answer" gorm:"not null"`
	Status         bool      `json:"status" gorm:"not null"`
	Solution       string    `json:"solution_status" gorm:"not null"`
	EvaluatedTimes int       `json:"evaluate_times" gorm:"not null"`
	Times          int       `json:"done_times" gorm:"not null"`
	KidID          string    `json:"-" gorm:"not null"`
	Quiz           Quiz      `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
