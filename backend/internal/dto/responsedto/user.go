package responsedto

type UserRegisterPrecheck struct {
	Salt           string `json:"salt"`
	IterationCount int    `json:"iterationCount"`
}

type UserLoginPrecheck struct {
	Salt      string `json:"salt"`
	IterCount int    `json:"iter_count"`
	Nonce     string `json:"nonce"`
}

type UserLogin struct {
	ServerSignature string `json:"server_signature"`
	Token           string `json:"token"`
}

type UserProfile struct {
	Username            string `json:"username"`
	DisplayName         string `json:"display_name"`
	Bio                 string `json:"bio"`
	Color               string `json:"color"`
	EmailVerified       bool   `json:"email_verified"`
	PhoneNumberVerified bool   `json:"phone_number_verified"`
}

type UserGetTelegramSessionID struct {
	SessionID string `json:"session_id"`
}

type UserResetPasswordPrecheck struct {
	Salt           string `json:"salt"`
	IterationCount int    `json:"iterationCount"`
}
