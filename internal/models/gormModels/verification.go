package gormModels

import "time"

type EmailVerificationData struct {
	ID               uint      `gorm:"primaryKey; autoIncrement; not null" json:"id"`
	UserId           uint      `gorm:"column:user_id; unique; not null" json:"user_id"`
	Salt             string    `gorm:"column:salt; not null" json:"salt"`
	Email            string    `gorm:"column:email; not null" json:"email"`
	VerificationCode string    `gorm:"column:verification_code; not null" json:"verification_code"`
	ExpiresAt        time.Time `gorm:"column:expires_at; not null" json:"expires_at"`
}

func (EmailVerificationData) TableName() string {
	return "email_verification_data"
}

type PhoneNumberVerificationData struct {
	ID               uint      `gorm:"primaryKey; autoIncrement; not null" json:"id"`
	UserId           uint      `gorm:"column:user_id; unique; not null" json:"user_id"`
	Salt             string    `gorm:"column:salt; not null" json:"salt"`
	PhoneNumber      string    `gorm:"column:phone_number; not null" json:"phone_number"`
	VerificationCode string    `gorm:"column:verification_code; not null" json:"verification_code"`
	ExpiresAt        time.Time `gorm:"column:expires_at; not null" json:"expires_at"`
}

func (PhoneNumberVerificationData) TableName() string {
	return "phone_number_verification_data"
}
