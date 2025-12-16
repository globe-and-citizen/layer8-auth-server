package responsedto

type ClientRegisterPrecheck struct {
	Salt           string `json:"salt"`
	IterationCount int    `json:"iterationCount"`
}

type ClientLoginPrecheck struct {
	Salt      string `json:"salt"`
	IterCount int    `json:"iter_count"`
	Nonce     string `json:"nonce"`
}

type ClientLogin struct {
	ServerSignature string `json:"server_signature"`
	Token           string `json:"token"`
}

type ClientProfile struct {
	ID              string `json:"id"`
	Secret          string `json:"secret"`
	Name            string `json:"name"`
	RedirectURI     string `json:"redirect_uri"`
	BackendURI      string `json:"backend_uri"`
	X509Certificate string `json:"x509_certificate"`
}
