package entities

import "time"

type Evaluate struct {
	ID             string     `json:"E_id" gorm:"primaryKey"`
	Status         bool       `json:"status" gorm:"not null"`
	Solution       string     `json:"solution_status" gorm:"not null"`
	EvaluatedTimes int        `json:"evaluate_times" gorm:"not null"`
	PeriodID       int        `json:"period_id" gorm:"not null"`
	KidID          string     `json:"-" gorm:"not null"`
	CompletedAt    *time.Time `json:"completed_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
