package auth

import (
	"context"

	"github.com/Template7/backend/internal/db"
	"github.com/Template7/backend/internal/t7Error"
)

func (s *service) Login(ctx context.Context, username string, password string) (string, error) {
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

func (s *service) CheckPermission(ctx context.Context, sub, obj, act string) bool {
	ok, err := s.core.Enforce(sub, obj, act)
	if err != nil {
		s.log.WithContext(ctx).WithError(err).With("sub", sub).With("obj", obj).With("act", act).Error("fail to check permission")
		return false
	}
	return ok
}
