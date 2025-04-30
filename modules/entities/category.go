package entities

type Category struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Category string `json:"category" gorm:"not null"`
}
