package requestdto

type AuthorizationToken struct {
	ClientUUID        string `json:"client_oauth_uuid" validate:"required"`
	ClientSecret      string `json:"client_oauth_secret" validate:"required"`
	AuthorizationCode string `json:"authorization_code" validate:"required"`
}

type ZkMetadata struct {
	ClientUUID   string `json:"client_oauth_uuid" validate:"required"`
	ClientSecret string `json:"client_oauth_secret" validate:"required"`
}
