package consts

import "time"

const SaltSize = 32
const SecretSize = 16

const EmailSendTimeout = time.Second * 10

const MiddlewareKeyUserUserID = "user_id"
const MiddlewareKeyClientUsername = "username"
const MiddlewareKeyClientClientID = "client_id"
