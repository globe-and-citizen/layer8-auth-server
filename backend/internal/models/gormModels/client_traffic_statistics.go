package gormModels

import "time"

type ClientTrafficStatistics struct {
	ID                         uint      `gorm:"primaryKey; autoIncrement; not null" json:"id"`
	ClientId                   string    `gorm:"column:client_id; unique; not null" json:"client_id"`
	TotalUsageBytes            int       `gorm:"column:total_usage_bytes" json:"total_usage_bytes"`
	UnpaidAmount               float64   `gorm:"column:unpaid_amount" json:"unpaid_amount"`
	LastTrafficUpdateTimestamp time.Time `gorm:"column:last_traffic_update_timestamp" json:"last_traffic_update_timestamp"`
}

func (ClientTrafficStatistics) TableName() string {
	return "client_traffic_statistics"
}
