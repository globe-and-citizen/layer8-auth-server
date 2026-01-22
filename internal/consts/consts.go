package consts

import "time"

const SaltSize = 32
const SecretSize = 16

const EmailSendTimeout = time.Second * 10

const MiddlewareKeyUserUsername = "user_username"
const MiddlewareKeyUserUserID = "user_id"
const MiddlewareKeyClientUsername = "client_username"
const MiddlewareKeyClientClientID = "client_id"
const MiddlewareKeyOAuthScopes = "oauth_scopes"
