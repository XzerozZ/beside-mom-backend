package entities

import "time"

type Appointment struct {
	ID          string    `json:"a_id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Date        time.Time `json:"date" gorm:"not null"`
	StartTime   time.Time `json:"start_time" gorm:"not null"`
	Building    string    `json:"building" gorm:"not null"`
	Requirement string    `json:"requirement"`
	Doctor      string    `json:"doctor" gorm:"not null"`
	Status      int       `json:"status" gorm:"not null"`
	UserID      string    `json:"user_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User        User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
	CreatedAt   time.Time `json:"created_at"`
}
