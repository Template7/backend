package t7Error

type Code string
type Message string

var (
	UnAuthorized = &Error{
		Code:    "A001",
		Message: "token unauthorized",
	}

	InvalidBody = &Error{
		Code:    "B001",
		Message: "invalid body",
	}

	DbConnectionFail = &Error{
		Code:    "D001",
		Message: "db connection fail",
	}
	DbOperationFail = &Error{
		Code:    "D002",
		Message: "db operation fail",
	}
	InvalidDocumentId = &Error{
		Code:    "D003",
		Message: "invalid document id",
	}

	HttpOperationFail = &Error{
		Code:    "H001",
		Message: "http operation fail",
	}
	HttpUnexpectedResponseCode = &Error{
		Code:    "H002",
		Message: "http unexpected response code",
	}

	UserNotfound = &Error{
		Code:    "U001",
		Message: "user not found",
	}
	UserAlreadyExist = &Error{
		Code:    "U002",
		Message: "user already exist",
	}
	MultipleMatchedUser = &Error{
		Code:    "U003",
		Message: "multiple matched users",
	}

	OrderNotFound = &Error{
		Code:    "O001",
		Message: "order not found",
	}
	InvalidOrderId = &Error{
		Code:    "O002",
		Message: "invalid order id",
	}

	VerifyCodeExpired = &Error{
		Code:    "V001",
		Message: "verify code expired",
	}
	IncorrectVerifyCode = &Error{
		Code:    "V002",
		Message: "incorrect verify code",
	}

	RedisOperationFail = &Error{
		Code:    "R001",
		Message: "redis operation fail",
	}

	TokenSignFail = &Error{
		Code:    "T001",
		Message: "token sign fail",
	}

	DecodeFail = &Error{
		Code:    "C001",
		Message: "decode fail",
	}
)
