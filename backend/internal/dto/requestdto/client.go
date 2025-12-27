package requestdto

import "globe-and-citizen/layer8/auth-server/backend/pkg/scram"

type ClientRegisterPrecheck struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}

type ClientRegister struct {
	Name                             string `json:"name" validate:"required"`
	RedirectURI                      string `json:"redirect_uri" validate:"required"`
	BackendURI                       string `json:"backend_uri" validate:"required"`
	scram.ClientRegisterFinalMessage `json:",inline"`
}

type ClientLoginPrecheck struct {
	scram.ClientLoginFirstMessage `json:",inline"`
}

type ClientLogin struct {
	CNonce                        string `json:"c_nonce" validate:"required"`
	scram.ClientLoginFinalMessage `json:",inline"`
}

type ClientCheckBackendURI struct {
	BackendURI string `json:"backend_uri" validate:"required"`
}

type ClientUploadNTorCertificate struct {
	Certificate string `json:"certificate" validate:"required"`
}
