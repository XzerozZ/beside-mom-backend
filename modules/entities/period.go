package entities

type Period struct {
	ID     int    `gorm:"primaryKey;autoIncrement"`
	Period string `json:"period" gorm:"not null"`
}
