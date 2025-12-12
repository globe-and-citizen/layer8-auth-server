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
