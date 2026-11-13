package user

import (
	"context"
	"errors"
	"testing"

	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/cache/cachetest"
	"github.com/Template7/backend/internal/db/dbtest"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/models"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
	"gorm.io/gorm"
)

var ctx = context.Background()

// fakeRoles keeps roles in memory; failAdd makes AddUserRole fail.
type fakeRoles struct {
	auth.Auth
	roles   map[string]authV1.Role
	failAdd bool
}

func (f *fakeRoles) AddUserRole(_ context.Context, userId string, role authV1.Role) error {
	if f.failAdd {
		return errors.New("casbin down")
	}
	f.roles[userId] = role
	return nil
}

func (f *fakeRoles) DeleteUserRoles(_ context.Context, userId string) error {
	delete(f.roles, userId)
	return nil
}

func (f *fakeRoles) GetUserRole(_ context.Context, userId string) (authV1.Role, error) {
	r, ok := f.roles[userId]
	if !ok {
		return 0, t7Error.UserHasNoRole
	}
	return r, nil
}

type fixture struct {
	svc   Service
	core  *gorm.DB
	roles *fakeRoles
	cache *cachetest.Memory
}

func setup(t *testing.T) fixture {
	t.Helper()
	client, core := dbtest.New(t)
	f := fixture{core: core, roles: &fakeRoles{roles: map[string]authV1.Role{}}, cache: cachetest.New()}
	f.svc = New(f.roles, client, f.cache, logger.New("error", "console", "test"))
	return f
}

func newUserReq(name string) *userV1.CreateUserRequest {
	return &userV1.CreateUserRequest{Username: name, Password: "pass-" + name, Role: authV1.Role_user, Nickname: name}
}

func wantCode(t *testing.T, err error, want *t7Error.Error) {
	t.Helper()
	got, ok := t7Error.ToT7Error(err)
	if !ok || got.Code != want.Code {
		t.Fatalf("want %s, got %v", want.Message, err)
	}
}

func TestCreateUser_CreatesUserWalletRoleAndCode(t *testing.T) {
	f := setup(t)

	userId, code, err := f.svc.CreateUser(ctx, newUserReq("alice"))
	if err != nil || userId == "" || code == "" {
		t.Fatalf("create: id=%q code=%q err=%v", userId, code, err)
	}
	if status, _ := f.svc.GetStatus(ctx, userId); status != authV1.AccountStatus_initialized {
		t.Fatalf("status = %v, want initialized", status)
	}
	if n := dbtest.Count(t, f.core, &models.Wallet{}); n != 1 {
		t.Fatalf("wallets = %d, want 1", n)
	}
	if n := dbtest.Count(t, f.core, &models.Balance{}); n != int64(len(defaultCurrencies)) {
		t.Fatalf("balances = %d", n)
	}
	if f.roles.roles[userId] != authV1.Role_user {
		t.Fatal("role not assigned")
	}
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	f := setup(t)
	if _, _, err := f.svc.CreateUser(ctx, newUserReq("alice")); err != nil {
		t.Fatal(err)
	}
	_, _, err := f.svc.CreateUser(ctx, newUserReq("alice"))
	wantCode(t, err, t7Error.UserAlreadyExist)
	if n := dbtest.Count(t, f.core, &models.Wallet{}); n != 1 {
		t.Fatalf("wallets = %d, want 1 (no partial second user)", n)
	}
}

func TestCreateUser_RoleFailureRollsBackUser(t *testing.T) {
	f := setup(t)
	f.roles.failAdd = true

	if _, _, err := f.svc.CreateUser(ctx, newUserReq("alice")); err == nil {
		t.Fatal("want error")
	}
	if n := dbtest.Count(t, f.core, &models.User{}); n != 0 {
		t.Fatalf("users = %d, want 0", n)
	}
	if n := dbtest.Count(t, f.core, &models.Wallet{}); n != 0 {
		t.Fatalf("wallets = %d, want 0", n)
	}
}

func TestActivateUser_CodeIsOneTime(t *testing.T) {
	f := setup(t)
	userId, code, _ := f.svc.CreateUser(ctx, newUserReq("alice"))

	wantCode(t, f.svc.ActivateUser(ctx, userId, "wrong"), t7Error.IncorrectVerifyCode)
	wantCode(t, f.svc.ActivateUser(ctx, userId, ""), t7Error.IncorrectVerifyCode)

	if err := f.svc.ActivateUser(ctx, userId, code); err != nil {
		t.Fatal(err)
	}
	if status, _ := f.svc.GetStatus(ctx, userId); status != authV1.AccountStatus_activated {
		t.Fatalf("status = %v", status)
	}
	if f.cache.Has(userId) {
		t.Fatal("code must be deleted after use")
	}
	wantCode(t, f.svc.ActivateUser(ctx, userId, code), t7Error.IncorrectVerifyCode)
}

func TestDeleteUser(t *testing.T) {
	f := setup(t)
	userId, _, _ := f.svc.CreateUser(ctx, newUserReq("alice"))

	if err := f.svc.DeleteUser(ctx, userId); err != nil {
		t.Fatal(err)
	}
	if _, ok := f.roles.roles[userId]; ok {
		t.Fatal("role must be removed")
	}
	if f.cache.Has(userId) {
		t.Fatal("activation code must be removed")
	}
	_, err := f.svc.GetStatus(ctx, userId)
	wantCode(t, err, t7Error.UserNotfound)
	wantCode(t, f.svc.DeleteUser(ctx, userId), t7Error.UserNotfound)
}

func TestGetInfoAndUpdate(t *testing.T) {
	f := setup(t)
	userId, _, _ := f.svc.CreateUser(ctx, newUserReq("alice"))

	if err := f.svc.UpdateInfo(ctx, userId, models.UserInfo{Nickname: "Ali"}); err != nil {
		t.Fatal(err)
	}
	info, err := f.svc.GetInfo(ctx, userId)
	if err != nil || info.Nickname != "Ali" || info.Role != authV1.Role_user {
		t.Fatalf("info = %+v, %v", info, err)
	}
	wantCode(t, f.svc.UpdateInfo(ctx, "nobody", models.UserInfo{Nickname: "x"}), t7Error.UserNotfound)
}
