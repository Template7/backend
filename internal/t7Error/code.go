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
	}

	PasswordIncorrect = &Error{
		Code:    codePasswordIncorrect,
		Message: "password incorrect",
		Type:    typeAuth,
	}

	UnAuthorized = &Error{
		Code:    codeUnAuthorized,
		Message: "token unauthorized",
		Type:    typeAuth,
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
	}
	DbOperationFail = &Error{
		Code:    codeDbOperationFail,
		Message: "db operation fail",
		Type:    typeDb,
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
	}
	HttpUnexpectedResponseCode = &Error{
		Code:    codeHttpUnexpectedResponseCode,
		Message: "http unexpected response code",
		Type:    typeInvalidData,
	}
	HashFail = &Error{
		Code:    codeHashFail,
		Message: "Hash fail",
		Type:    typeAuth,
	}

	UserNotfound = &Error{
		Code:    codeUserNotfound,
		Message: "user not found",
		Type:    typeInvalidData,
	}
	UserAlreadyExist = &Error{
		Code:    codeUserAlreadyExist,
		Message: "user already exist",
		Type:    typeInvalidData,
	}
	SignInFail = &Error{
		Code:    codeSignInFail,
		Message: "Invalid user name or password",
		Type:    typeInvalidData,
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
	}

	TokenSignFail = &Error{
		Code:    codeTokenSignFail,
		Message: "token sign fail",
		Type:    typeAuth,
	}
	TokenAssertionFail = &Error{
		Code:    codeTokenAssertionFail,
		Message: "token assertion fail",
		Type:    typeInvalidData,
	}
	InvalidToken = &Error{
		Code:    codeInvalidToken,
		Message: "invalid token",
		Type:    typeAuth,
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
	}
)
