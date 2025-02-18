package entities

import "time"

type Video struct {
	ID        string    `json:"video_id" gorm:"primaryKey"`
	Link      string    `json:"video_link" gorm:"not null"`
	View      int       `json:"video_view" gorm:"default:0"`
	UserID    string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
