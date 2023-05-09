package entity

import "time"

type License struct {
	Id           int       `json:"id" gorm:"column:id"`
	LicenseKey   string    `json:"license_key" gorm:"column:license_key"`
	HardwareCode string    `json:"hardware_code" gorm:"column:hardware_code"`
	Enable       bool      `json:"enable" gorm:"column:enable"`
	StartAt      time.Time `json:"start_at" gorm:"column:start_at"`
	EndAt        time.Time `json:"end_at" gorm:"column:end_at"`
	Owner        string    `json:"owner" gorm:"column:owner"`
}

func (License) TableName() string {
	return "license"
}
