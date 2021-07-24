package util

func ParseBearerToken(token string) string {
	return token[7:]
}
