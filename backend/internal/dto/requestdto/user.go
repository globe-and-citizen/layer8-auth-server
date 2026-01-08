package requestdto

import "globe-and-citizen/layer8/auth-server/backend/pkg/scram"

type UserRegisterPrecheck struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}

type UserRegister struct {
	PublicKey                        []byte `json:"public_key" validate:"required"`
	scram.ClientRegisterFinalMessage `json:",inline"`
}

type UserLoginPrecheck struct {
	scram.ClientLoginFirstMessage `json:",inline"`
}

type UserLogin struct {
	CNonce                        string `json:"c_nonce" validate:"required"` // fixme ClientNonce shouldn't be here, needs to be saved from precheck message
	scram.ClientLoginFinalMessage `json:",inline"`
}

type UserResetPasswordPrecheck struct {
	Username string `json:"username" validate:"required"`
}

type UserResetPassword struct {
	Signature                        []byte `json:"signature" validate:"required"`
	scram.ClientRegisterFinalMessage `json:",inline"`
}

type UserMetadataUpdate struct {
	DisplayName string `json:"display_name"`
	Color       string `json:"color"`
	Bio         string `json:"bio"`
}

type UserEmailVerify struct {
	Email string `json:"email" validate:"required,email"`
}

type UserCheckEmailVerificationCode struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type UserCheckPhoneNumberVerificationCode struct {
	VerificationCode string `json:"verification_code"`
}
