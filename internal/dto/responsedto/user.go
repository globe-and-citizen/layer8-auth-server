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
	Username            string    `json:"username"`
	DisplayName         string    `json:"display_name"`
	Bio                 string    `json:"bio"`
	Color               string    `json:"color"`
	EmailVerified       bool      `json:"email_verified"`
	LastEmailVerifiedAt time.Time `json:"last_email_verified_at"`
	PhoneNumberVerified bool      `json:"phone_number_verified"`
	LastPhoneVerifiedAt time.Time `json:"last_phone_verified_at"`
	PhoneLocation       string    `json:"phone_location"`
}

type UserGetTelegramSessionID struct {
	SessionID string `json:"session_id"`
}
