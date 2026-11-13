package auth

import (
	"context"

	"github.com/Template7/backend/internal/db"
	"github.com/Template7/backend/internal/t7Error"
)

// Login checks the credentials and returns a signed token.
//
// Failed attempts are counted per username and per client IP. Once either count
// reaches its limit, logins are refused with TooManyAttempts until the window
// ends, even with the right password. A successful login clears the username's
// count. If the cache is unavailable, logins are not limited (fail open).
func (s *service) Login(ctx context.Context, username string, password string, clientIp string) (string, error) {
	keys := s.failureKeys(username, clientIp)
	if s.tooManyFailures(ctx, keys) {
		return "", t7Error.TooManyAttempts
	}

	token, err := s.checkCredentials(ctx, username, password)
	if err != nil {
		if t7Err, ok := t7Error.ToT7Error(err); ok && t7Err.Code == t7Error.SignInFail.Code {
			s.countFailure(ctx, keys)
		}
		return "", err
	}

	if err := s.cache.ResetLoginFailures(ctx, keys.user); err != nil {
		s.log.WithContext(ctx).WithError(err).Warn("fail to reset login failures")
	}
	return token, nil
}

func (s *service) checkCredentials(ctx context.Context, username string, password string) (string, error) {
	user, err := s.db.GetUser(ctx, username)
	if db.IsNotFound(err) {
		// same answer and similar timing as a wrong password, so usernames can't be probed
		checkPasswordHash(password, dummyPasswordHash)
		return "", t7Error.SignInFail
	}
	if err != nil {
		return "", t7Error.DbOperationFail.WithDetail(err.Error())
	}
	if !checkPasswordHash(password, user.Password) {
		return "", t7Error.SignInFail
	}

	role, err := s.GetUserRole(ctx, user.Id)
	if err != nil {
		return "", err
	}
	return s.genUserToken(user.Id, role, user.Status)
}

type failureKeys struct {
	user string
	ip   string // empty when the client IP is unknown
}

func (s *service) failureKeys(username string, clientIp string) failureKeys {
	k := failureKeys{user: "user:" + username}
	if clientIp != "" {
		k.ip = "ip:" + clientIp
	}
	return k
}

func (s *service) tooManyFailures(ctx context.Context, keys failureKeys) bool {
	check := func(key string, limit int64) bool {
		if key == "" {
			return false
		}
		n, err := s.cache.GetLoginFailures(ctx, key)
		if err != nil {
			s.log.WithContext(ctx).WithError(err).Error("fail to read login failures, not limiting")
			return false
		}
		return n >= limit
	}
	return check(keys.user, s.limits.perUser) || check(keys.ip, s.limits.perIp)
}

func (s *service) countFailure(ctx context.Context, keys failureKeys) {
	for _, key := range []string{keys.user, keys.ip} {
		if key == "" {
			continue
		}
		if _, err := s.cache.IncrLoginFailures(ctx, key, s.limits.window); err != nil {
			s.log.WithContext(ctx).WithError(err).Error("fail to count login failure")
		}
	}
}

func (s *service) CheckPermission(ctx context.Context, sub, obj, act string) bool {
	ok, err := s.core.Enforce(sub, obj, act)
	if err != nil {
		s.log.WithContext(ctx).WithError(err).With("sub", sub).With("obj", obj).With("act", act).Error("fail to check permission")
		return false
	}
	return ok
}
