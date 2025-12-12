package responsedto

type Response struct {
	IsSuccess bool        `json:"is_success"`
	Message   string      `json:"message"`
	Error     interface{} `json:"errors"`
	Data      interface{} `json:"data"`
}

type (
	LoginClientResponse struct {
		ServerSignature string `json:"server_signature"`
		Token           string `json:"token"`
	}

	ProfileResponse struct {
		Username            string `json:"username"`
		DisplayName         string `json:"display_name"`
		Bio                 string `json:"bio"`
		Color               string `json:"color"`
		EmailVerified       bool   `json:"email_verified"`
		PhoneNumberVerified bool   `json:"phone_number_verified"`
	}

	ClientResponse struct {
		ID              string `json:"id"`
		Secret          string `json:"secret"`
		Name            string `json:"name"`
		RedirectURI     string `json:"redirect_uri"`
		BackendURI      string `json:"backend_uri"`
		X509Certificate string `json:"x509_certificate"`
	}

	RegisterClientPrecheckResponse struct {
		Salt           string `json:"salt"`
		IterationCount int    `json:"iterationCount"`
	}

	ResetPasswordPrecheckResponse struct {
		Salt           string `json:"salt"`
		IterationCount int    `json:"iterationCount"`
	}

	ClientUnpaidAmountResponse struct {
		UnpaidAmount int `json:"unpaid_amount"`
	}
)
