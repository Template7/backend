package auth

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Template7/backend/internal/cache"
	"github.com/Template7/backend/internal/db"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/models"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"gorm.io/gorm"
)

type userDb struct {
	db.Client
	users  map[string]models.User // by username
	status map[string]authV1.AccountStatus
}

func (d *userDb) GetUser(_ context.Context, username string) (models.User, error) {
	u, ok := d.users[username]
	if !ok {
		return models.User{}, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (d *userDb) SetUserStatus(_ context.Context, userId string, status authV1.AccountStatus) error {
	d.status[userId] = status
	return nil
}

type memCache struct {
	cache.Interface
	codes map[string]string
}

func (c *memCache) GetUserActivationCode(_ context.Context, userId string) (string, error) {
	code, ok := c.codes[userId]
	if !ok {
		return "", errors.New("redis: nil")
	}
	return code, nil
}

func (c *memCache) DeleteUserActivationCode(_ context.Context, userId string) error {
	delete(c.codes, userId)
	return nil
}

func newUserTestService(t *testing.T) (*service, *userDb, *memCache) {
	t.Helper()
	hash, err := hashedPassword("correct-password")
	if err != nil {
		t.Fatal(err)
	}
	d := &userDb{
		users:  map[string]models.User{"alice": {Id: "u-alice", Username: "alice", Password: hash}},
		status: map[string]authV1.AccountStatus{},
	}
	c := &memCache{codes: map[string]string{"u-alice": "the-code"}}
	return &service{db: d, cache: c, jwtSecret: []byte(testSecret), log: logger.New("error", "console", "test")}, d, c
}

func TestLogin_UnknownUserAndWrongPasswordLookTheSame(t *testing.T) {
	s, _, _ := newUserTestService(t)

	_, errUnknown := s.Login(context.Background(), "nobody", "whatever")
	_, errWrong := s.Login(context.Background(), "alice", "wrong-password")

	for name, err := range map[string]error{"unknown user": errUnknown, "wrong password": errWrong} {
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok || t7Err.Code != t7Error.SignInFail.Code || t7Err.GetStatus() != http.StatusUnauthorized {
			t.Errorf("%s: want SignInFail/401, got %v", name, err)
		}
	}
}

func TestActivateUser_CodeIsOneTime(t *testing.T) {
	s, d, c := newUserTestService(t)

	if !s.ActivateUser(context.Background(), "u-alice", "the-code") {
		t.Fatal("first activation should succeed")
	}
	if d.status["u-alice"] != authV1.AccountStatus_activated {
		t.Fatal("status not updated")
	}
	if _, ok := c.codes["u-alice"]; ok {
		t.Fatal("code must be deleted after use")
	}
	if s.ActivateUser(context.Background(), "u-alice", "the-code") {
		t.Fatal("code must not work twice")
	}
}

func TestActivateUser_RejectsWrongOrEmptyCode(t *testing.T) {
	s, d, _ := newUserTestService(t)

	for _, code := range []string{"wrong", ""} {
		if s.ActivateUser(context.Background(), "u-alice", code) {
			t.Fatalf("code %q must be rejected", code)
		}
	}
	if _, changed := d.status["u-alice"]; changed {
		t.Fatal("status must not change")
	}
}

func TestBcryptCostIsAtLeastDefault(t *testing.T) {
	if bcryptCost < 10 {
		t.Fatalf("bcrypt cost %d is below 10", bcryptCost)
	}
}
