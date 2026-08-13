package gormModels

import "time"

type UserMetadata struct {
	ID                    uint       `gorm:"column:id; primaryKey; not null"`
	DisplayName           string     `gorm:"column:display_name"`
	DisplayNameUpdatedAt  *time.Time `gorm:"column:display_name_updated_at"`
	Color                 string     `gorm:"column:color"`
	ColorUpdatedAt        *time.Time `gorm:"column:color_updated_at"`
	Bio                   string     `gorm:"column:bio"`
	BioUpdatedAt          *time.Time `gorm:"column:bio_updated_at"`
	IsEmailVerified       bool       `gorm:"column:is_email_verified; default:false"`
	EmailVerifiedAt       *time.Time `gorm:"column:email_verified_at"`
	IsPhoneNumberVerified bool       `gorm:"column:is_phone_number_verified; default:false"`
	Location              string     `gorm:"column:location"`
	PhoneVerifiedAt       *time.Time `gorm:"column:phone_number_verified_at"`
	CreatedAt             time.Time  `gorm:"column:created_at; default:CURRENT_TIMESTAMP"`
	UpdatedAt             time.Time  `gorm:"column:updated_at; default:CURRENT_TIMESTAMP; autoUpdateTime"`
}

func (UserMetadata) TableName() string {
	return "user_metadata"
}
