package requestdto

type UserRegisterPrecheck struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}

type UserRegister struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	PublicKey []byte `json:"public_key" validate:"required"`
	StoredKey string `json:"stored_key" validate:"required"`
	ServerKey string `json:"server_key" validate:"required"`
}

type UserLoginPrecheck struct {
	Username string `json:"username" validate:"required"`
	CNonce   string `json:"c_nonce" validate:"required"`
}

type UserLogin struct {
	Username    string `json:"username" validate:"required"`
	Nonce       string `json:"nonce" validate:"required"`
	CNonce      string `json:"c_nonce" validate:"required"`
	ClientProof string `json:"client_proof" validate:"required"`
}

type UserMetadataUpdate struct {
	DisplayName string `json:"display_name" validate:"required"`
	Color       string `json:"color" validate:"required"`
	Bio         string `json:"bio" validate:"required"`
}

type UserEmailVerify struct {
	Email string `json:"email" validate:"required,email"`
}

type CheckEmailVerificationCode struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}
