package entities

import "time"

type Appointment struct {
	ID          string    `json:"a_id" gorm:"primaryKey"`
	Date        time.Time `json:"date" gorm:"not null"`
	StartTime   time.Time `json:"start_time" gorm:"not null"`
	Building    string    `json:"building" gorm:"not null"`
	Requirement string    `json:"requirement"`
	Doctor      string    `json:"doctor" gorm:"not null"`
	Status      int       `json:"status" gorm:"not null"`
	UserID      string    `json:"user_id" gorm:"not null"`
	User        User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
	CreatedAt   time.Time `json:"created_at"`
}
