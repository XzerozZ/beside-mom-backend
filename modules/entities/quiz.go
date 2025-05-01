package entities

import "time"

type Quiz struct {
	ID          int       `json:"quiz_id" gorm:"primaryKey;autoIncrement"`
	Question    string    `json:"question" gorm:"not null"`
	Description string    `json:"desc" gorm:"not null"`
	Solution    string    `json:"solution" gorm:"not null"`
	Suggestion  string    `json:"suggestion" gorm:"not null"`
	Banner      string    `json:"banner" gorm:"not null"`
	CategoryID  int       `json:"category_id" gorm:"not null"`
	PeriodID    int       `json:"period_id" gorm:"not null"`
	Category    Category  `json:"category" gorm:"foreignKey:CategoryID;references:ID;"`
	Period      Period    `json:"period" gorm:"foreignKey:PeriodID;references:ID;"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
