package entities

import "time"

type Kid struct {
	ID          string    `json:"u_id" gorm:"primaryKey" `
	Firstname   string    `json:"fname" gorm:"not null" `
	Lastname    string    `json:"lname" gorm:"not null"`
	Username    string    `json:"uname" gorm:"not null"`
	Sex         string    `json:"sex" gorm:"not null"`
	BirthDate   string    `json:"birth_date" gorm:"not null"`
	BloodType   string    `json:"blood_type" gorm:"not null"`
	BirthWeight float64   `json:"birth_weight" gorm:"not null"`
	BirthLength float64   `json:"birth_length" gorm:"not null"`
	Note        string    `json:"note" gorm:"not null"`
	ImageLink   string    `json:"image_link"`
	UserID      string    `json:"user_id" gorm:"unique;foreignKey:UserID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
