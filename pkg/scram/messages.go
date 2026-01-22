package scram

type ClientRegisterFirstMessage struct {
	Username string `json:"username" validate:"required"`
}

type ServerRegisterFirstMessage struct {
	Salt           string `json:"salt" validate:"required,base64,min=22,max=24"`
	IterationCount int    `json:"iteration_count" validate:"required,min=1"`
}

type ClientRegisterFinalMessage struct {
	Username  string `json:"username" validate:"required"`
	StoredKey string `json:"stored_key" validate:"required"`
	ServerKey string `json:"server_key" validate:"required"`
}

type ClientLoginFirstMessage struct {
	Username    string `json:"username" validate:"required"`
	ClientNonce string `json:"c_nonce" validate:"required,base64,min=22,max=24"`
}

type ServerLoginFirstMessage struct {
	Salt           string `json:"salt" validate:"required,base64,min=22,max=24"`
	IterationCount int    `json:"iteration_count" validate:"required,min=1"`
	Nonce          string `json:"nonce" validate:"required,base64,min=44,max=48"`
}

type ClientLoginFinalMessage struct {
	Username       string `json:"username" validate:"required"`
	ChannelBinding string `json:"c" validate:"required,base64"`
	Nonce          string `json:"nonce" validate:"required,base64"`
	ClientProof    string `json:"client_proof" validate:"required,base64"`
}

type ServerLoginFinalMessage struct {
	ServerSignature string `json:"verifier" validate:"required,base64"`
}
