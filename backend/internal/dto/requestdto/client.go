package requestdto

type ClientRegisterPrecheck struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}

type ClientRegister struct {
	Name        string `json:"name" validate:"required"`
	RedirectURI string `json:"redirect_uri" validate:"required"`
	BackendURI  string `json:"backend_uri" validate:"required"`
	Username    string `json:"username" validate:"required,min=3,max=50"`
	StoredKey   string `json:"stored_key" validate:"required"`
	ServerKey   string `json:"server_key" validate:"required"`
}

type ClientLoginPrecheck struct {
	Username string `json:"username" validate:"required"`
	CNonce   string `json:"c_nonce" validate:"required"`
}

type ClientLogin struct {
	Username    string `json:"username" validate:"required"`
	Nonce       string `json:"nonce" validate:"required"`
	CNonce      string `json:"c_nonce" validate:"required"`
	ClientProof string `json:"client_proof" validate:"required"`
}

type ClientCheckBackendURI struct {
	BackendURI string `json:"backend_uri" validate:"required"`
}

type ClientUploadNTorCertificate struct {
	Certificate string `json:"certificate" validate:"required"`
}
