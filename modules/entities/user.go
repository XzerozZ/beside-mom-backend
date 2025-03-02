package entities

import "time"

type User struct {
	ID        string    `json:"u_id" gorm:"primaryKey"`
	Firstname string    `json:"fname"`
	Lastname  string    `json:"lname"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-"`
	ImageLink string    `json:"image_link"`
	RoleID    int       `json:"-" gorm:"not null"`
	Role      Role      `json:"role" gorm:"foreignKey:RoleID"`
	Kid       []Kid     `json:"kids,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
