package responsedto

type (
	UserRegisterPrecheck struct {
		Salt           string `json:"salt"`
		IterationCount int    `json:"iterationCount"`
	}

	UserLoginPrecheck struct {
		Salt      string `json:"salt"`
		IterCount int    `json:"iter_count"`
		Nonce     string `json:"nonce"`
	}

	UserLogin struct {
		ServerSignature string `json:"server_signature"`
		Token           string `json:"token"`
	}

	UserGetTelegramSessionID struct {
		SessionID string `json:"session_id"`
	}
)
