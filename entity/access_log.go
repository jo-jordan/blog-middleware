package entity

import "time"

type AccessLog struct {
	ID         uint64
	RealIP     uint      `gorm:"column:real_ip"`
	Time       time.Time `gorm:"column:time"`
	Type       uint      `gorm:"column:type"`
	ResourceID uint64    `gorm:"column:res_id"`
}

type Tabler interface {
	TableName() string
}

func (AccessLog) TableName() string {
	return "access_log"
}
