package entities

import "time"

type Kid struct {
	ID          string    `json:"u_id" gorm:"primaryKey"`
	Firstname   string    `json:"fname" gorm:"not null"`
	Lastname    string    `json:"lname" gorm:"not null"`
	Username    string    `json:"uname" gorm:"not null"`
	Sex         string    `json:"sex" gorm:"not null"`
	BirthDate   time.Time `json:"birth_date" gorm:"not null;type:date"`
	BloodType   string    `json:"blood_type" gorm:"not null"`
	RHType      string    `json:"rh_type"`
	BirthWeight float64   `json:"weight" gorm:"not null"`
	BirthLength float64   `json:"length" gorm:"not null"`
	Note        string    `json:"note" gorm:"not null"`
	ImageLink   string    `json:"image_link"`
	UserID      string    `json:"user_id" gorm:"not null"`
	Growth      []Growth  `json:"growth" gorm:"foreignKey:KidID"`
	User        User      `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
