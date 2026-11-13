package t7Error

import "net/http"

const (
	codeUnAuthorized Code = iota + 1024
	codeInvalidBody
	codeDbConnectionFail
	codeDbOperationFail
	codeInvalidDocumentId
	codeHttpOperationFail
	codeHttpUnexpectedResponseCode
	codeHashFail
	codeUserNotfound
	codeUserAlreadyExist
	codeSignInFail
	codeVerifyCodeExpired
	codeIncorrectVerifyCode
	codeRedisOperationFail
	codeTokenSignFail
	codeTokenAssertionFail
	codeInvalidToken
	codeDecodeFail
	codeWalletNotFound
	codePasswordIncorrect
	codeUserHasNoRole
	codeBlockedAccount
	codeUnknown
	codeInsufficientBalance
	codeForbidden
	codeAccountNotActivated
	codeTooManyAttempts
)

const (
	typeAuth Type = iota + 32
	typeInvalidData
	typeDb
	typeRedis
	typeNetwork
	typeOther
)

var (
	Unknown = &Error{
		Code:    codeUnknown,
		Message: "unknown",
		Type:    typeOther,
		status:  http.StatusInternalServerError,
	}

	BlockedAccount = &Error{
		Code:    codeBlockedAccount,
		Message: "blocked account",
		Type:    typeAuth,
		status:  http.StatusForbidden,
	}

	UserHasNoRole = &Error{
		Code:    codeUserHasNoRole,
		Message: "user has no role",
		Type:    typeAuth,
		status:  http.StatusForbidden,
	}

	PasswordIncorrect = &Error{
		Code:    codePasswordIncorrect,
		Message: "password incorrect",
		Type:    typeAuth,
		status:  http.StatusUnauthorized,
	}

	UnAuthorized = &Error{
		Code:    codeUnAuthorized,
		Message: "token unauthorized",
		Type:    typeAuth,
		status:  http.StatusUnauthorized,
	}

	InvalidBody = &Error{
		Code:    codeInvalidBody,
		Message: "invalid body",
		Type:    typeInvalidData,
	}

	DbConnectionFail = &Error{
		Code:    codeDbConnectionFail,
		Message: "db connection fail",
		Type:    typeDb,
		status:  http.StatusInternalServerError,
	}
	DbOperationFail = &Error{
		Code:    codeDbOperationFail,
		Message: "db operation fail",
		Type:    typeDb,
		status:  http.StatusInternalServerError,
	}
	InvalidDocumentId = &Error{
		Code:    codeInvalidDocumentId,
		Message: "invalid document id",
		Type:    typeInvalidData,
	}

	HttpOperationFail = &Error{
		Code:    codeHttpOperationFail,
		Message: "http operation fail",
		Type:    typeNetwork,
		status:  http.StatusBadGateway,
	}
	HttpUnexpectedResponseCode = &Error{
		Code:    codeHttpUnexpectedResponseCode,
		Message: "http unexpected response code",
		Type:    typeInvalidData,
		status:  http.StatusBadGateway,
	}
	HashFail = &Error{
		Code:    codeHashFail,
		Message: "Hash fail",
		Type:    typeAuth,
		status:  http.StatusInternalServerError,
	}

	UserNotfound = &Error{
		Code:    codeUserNotfound,
		Message: "user not found",
		Type:    typeInvalidData,
		status:  http.StatusNotFound,
	}
	UserAlreadyExist = &Error{
		Code:    codeUserAlreadyExist,
		Message: "user already exist",
		Type:    typeInvalidData,
		status:  http.StatusConflict,
	}
	SignInFail = &Error{
		Code:    codeSignInFail,
		Message: "Invalid user name or password",
		Type:    typeInvalidData,
		status:  http.StatusUnauthorized,
	}

	VerifyCodeExpired = &Error{
		Code:    codeVerifyCodeExpired,
		Message: "verify code expired",
		Type:    typeInvalidData,
	}
	IncorrectVerifyCode = &Error{
		Code:    codeIncorrectVerifyCode,
		Message: "incorrect verify code",
		Type:    typeInvalidData,
	}

	RedisOperationFail = &Error{
		Code:    codeRedisOperationFail,
		Message: "redis operation fail",
		Type:    typeRedis,
		status:  http.StatusInternalServerError,
	}

	TokenSignFail = &Error{
		Code:    codeTokenSignFail,
		Message: "token sign fail",
		Type:    typeAuth,
		status:  http.StatusInternalServerError,
	}
	TokenAssertionFail = &Error{
		Code:    codeTokenAssertionFail,
		Message: "token assertion fail",
		Type:    typeInvalidData,
		status:  http.StatusUnauthorized,
	}
	InvalidToken = &Error{
		Code:    codeInvalidToken,
		Message: "invalid token",
		Type:    typeAuth,
		status:  http.StatusUnauthorized,
	}

	DecodeFail = &Error{
		Code:    codeDecodeFail,
		Message: "decode fail",
		Type:    typeInvalidData,
	}

	WalletNotFound = &Error{
		Code:    codeWalletNotFound,
		Message: "wallet not found",
		Type:    typeDb,
		status:  http.StatusNotFound,
	}

	InsufficientBalance = &Error{
		Code:    codeInsufficientBalance,
		Message: "insufficient balance",
		Type:    typeInvalidData,
	}

	Forbidden = &Error{
		Code:    codeForbidden,
		Message: "permission denied",
		Type:    typeAuth,
		status:  http.StatusForbidden,
	}

	AccountNotActivated = &Error{
		Code:    codeAccountNotActivated,
		Message: "account not activated",
		Type:    typeAuth,
		status:  http.StatusForbidden,
	}

	TooManyAttempts = &Error{
		Code:    codeTooManyAttempts,
		Message: "too many failed attempts, try again later",
		Type:    typeAuth,
		status:  http.StatusTooManyRequests,
	}
)
