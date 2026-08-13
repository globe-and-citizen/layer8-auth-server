package responsedto

import (
	"globe-and-citizen/layer8/auth-server/pkg/scram"
	"time"
)

type UserRegisterPrecheck struct {
	scram.ServerRegisterFirstMessage `json:",inline"`
}

type UserLoginPrecheck struct {
	scram.ServerLoginFirstMessage `json:",inline"`
}

type UserLogin struct {
	scram.ServerLoginFinalMessage `json:",inline"`
	Token                         string `json:"token"`
}

type UserResetPasswordPrecheck struct {
	scram.ServerRegisterFirstMessage `json:",inline"`
}

type UserProfile struct {
	Username              string     `json:"username"`
	DisplayName           string     `json:"display_name"`
	DisplayNameUpdatedAt  *time.Time `json:"display_name_updated_at"`
	Bio                   string     `json:"bio"`
	BioUpdatedAt          *time.Time `json:"bio_updated_at"`
	Color                 string     `json:"color"`
	ColorUpdatedAt        *time.Time `json:"color_updated_at"`
	EmailVerified         bool       `json:"email_verified"`
	EmailVerifiedAt       *time.Time `json:"email_verified_at"`
	PhoneNumberVerified   bool       `json:"phone_number_verified"`
	PhoneNumberVerifiedAt *time.Time `json:"phone_number_verified_at"`
	PhoneLocation         string     `json:"phone_location"`
}

type UserGetTelegramSessionID struct {
	SessionID string `json:"session_id"`
}
