package config

const (
	ContextKeyGormTx = "ezutil.gormTx"

	MsgTransactionError = "error processing transaction"
	MsgNoPermission     = "user does not have permission to perform this action"
	MsgAuthMissingToken = "authorization token is missing"
	MsgAuthInvalidToken = "authorization token is invalid"
	MsgAuthUserNotFound = "user is not found"
	MsgInvalidJson      = "JSON is invalid or malformed"
	MsgMissingBody      = "request body is missing or empty"
)
