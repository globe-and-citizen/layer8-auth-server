package gormModels

import "time"

type UserMetadata struct {
	ID                    uint      `gorm:"column:id; primaryKey; not null"`
	DisplayName           string    `gorm:"column:display_name"`
	Color                 string    `gorm:"column:color"`
	Bio                   string    `gorm:"column:bio"`
	IsEmailVerified       bool      `gorm:"column:is_email_verified; default:false"`
	IsPhoneNumberVerified bool      `gorm:"column:is_phone_number_verified; default:false"`
	LastEmailVerifiedAt   time.Time `gorm:"column:last_email_verified_at"`
	LastPhoneVerifiedAt   time.Time `gorm:"column:last_phone_verified_at"`
	Location              string    `gorm:"column:location"`
	CreatedAt             time.Time `gorm:"column:created_at; default:CURRENT_TIMESTAMP"`
	UpdatedAt             time.Time `gorm:"column:updated_at; default:CURRENT_TIMESTAMP; autoUpdateTime"`
}

func (UserMetadata) TableName() string {
	return "user_metadata"
}
