package entities

import "time"

type Likes struct {
	UserID    string    `json:"Q_id" gorm:"primaryKey"`
	VideoID   string    `json:"-" gorm:"primaryKey"`
	User      User      `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Video     Video     `gorm:"foreignKey:VideoID;references:ID"`
	CreatedAt time.Time `json:"created_at"`
}
