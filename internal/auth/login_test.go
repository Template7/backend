package auth

import (
	"context"
	"testing"

	"github.com/Template7/backend/internal/cache/cachetest"
	"github.com/Template7/backend/internal/config"
	"github.com/Template7/backend/internal/db/dbtest"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/models"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
)

var ctx = context.Background()

const (
	maxPerUser = 3
	maxPerIp   = 5
)

// newRealService builds the auth service on SQLite with real casbin, and a user
// "alice" / "right-password" with the user role.
func newRealService(t *testing.T) Auth {
	t.Helper()
	client, core := dbtest.New(t)

	cfg := &config.Config{}
	cfg.Auth.RbacModelPath = "../../config/rbac_model.conf"
	cfg.Auth.JwtSecret = testSecret
	cfg.Auth.Login.MaxFailuresPerUser = maxPerUser
	cfg.Auth.Login.MaxFailuresPerIp = maxPerIp

	s := New(client, core, cachetest.New(), logger.New("error", "console", "test"), cfg)

	hash, err := HashPassword("right-password")
	if err != nil {
		t.Fatal(err)
	}
	if err := core.Create(&models.User{Id: "u-alice", Username: "alice", Password: hash, Status: authV1.AccountStatus_activated}).Error; err != nil {
		t.Fatal(err)
	}
	if err := s.AddUserRole(ctx, "u-alice", authV1.Role_user); err != nil {
		t.Fatal(err)
	}
	return s
}

func wantCode(t *testing.T, err error, want *t7Error.Error) {
	t.Helper()
	got, ok := t7Error.ToT7Error(err)
	if !ok || got.Code != want.Code {
		t.Fatalf("want %s, got %v", want.Message, err)
	}
}

func TestLogin_Success(t *testing.T) {
	s := newRealService(t)
	tk, err := s.Login(ctx, "alice", "right-password", "10.0.0.1")
	if err != nil || tk == "" {
		t.Fatalf("login: %q, %v", tk, err)
	}
}

func TestLogin_UsernameLockedAfterMaxFailures(t *testing.T) {
	s := newRealService(t)

	for i := 0; i < maxPerUser; i++ {
		_, err := s.Login(ctx, "alice", "wrong", "10.0.0.1")
		wantCode(t, err, t7Error.SignInFail)
	}
	// locked now, even with the right password and from another IP
	_, err := s.Login(ctx, "alice", "right-password", "10.0.0.2")
	wantCode(t, err, t7Error.TooManyAttempts)
}

func TestLogin_SuccessResetsUsernameCount(t *testing.T) {
	s := newRealService(t)

	for i := 0; i < maxPerUser-1; i++ {
		_, _ = s.Login(ctx, "alice", "wrong", "10.0.0.1")
	}
	if _, err := s.Login(ctx, "alice", "right-password", "10.0.0.1"); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < maxPerUser-1; i++ {
		_, _ = s.Login(ctx, "alice", "wrong", "10.0.0.1")
	}
	if _, err := s.Login(ctx, "alice", "right-password", "10.0.0.1"); err != nil {
		t.Fatalf("count should have been reset by the earlier success: %v", err)
	}
}

func TestLogin_IpLockedAcrossUsernames(t *testing.T) {
	s := newRealService(t)

	// spraying many usernames from one IP; unknown users count too
	for i := 0; i < maxPerIp; i++ {
		_, err := s.Login(ctx, "nobody-"+string(rune('a'+i)), "x", "10.0.0.9")
		wantCode(t, err, t7Error.SignInFail)
	}
	_, err := s.Login(ctx, "alice", "right-password", "10.0.0.9")
	wantCode(t, err, t7Error.TooManyAttempts)

	// other IPs are not affected
	if _, err := s.Login(ctx, "alice", "right-password", "10.0.0.10"); err != nil {
		t.Fatal(err)
	}
}

func TestSyncRolePolicies(t *testing.T) {
	s := newRealService(t)
	user := authV1.Role_user.String()

	if err := s.SyncRolePolicies(ctx, authV1.Role_user, []Rule{
		{Path: "/api/v1/a", Method: "GET"},
		{Path: "/api/v1/items/:id", Method: "POST"},
	}); err != nil {
		t.Fatal(err)
	}
	if !s.CheckPermission(ctx, user, "/api/v1/a", "GET") || !s.CheckPermission(ctx, user, "/api/v1/items/42", "POST") {
		t.Fatal("synced rules must be allowed")
	}

	// second sync: /a removed, /b added, unchanged rule kept; running it twice is harmless
	for i := 0; i < 2; i++ {
		if err := s.SyncRolePolicies(ctx, authV1.Role_user, []Rule{
			{Path: "/api/v1/items/:id", Method: "POST"},
			{Path: "/api/v1/b", Method: "GET"},
		}); err != nil {
			t.Fatal(err)
		}
	}
	if s.CheckPermission(ctx, user, "/api/v1/a", "GET") {
		t.Fatal("stale rule must be removed")
	}
	if !s.CheckPermission(ctx, user, "/api/v1/b", "GET") || !s.CheckPermission(ctx, user, "/api/v1/items/42", "POST") {
		t.Fatal("current rules must be allowed")
	}
	if s.CheckPermission(ctx, user, "/api/v1/items/42", "GET") {
		t.Fatal("method must match")
	}
	if !s.CheckPermission(ctx, authV1.Role_admin.String(), "/anything", "DELETE") {
		t.Fatal("admin passes every check")
	}
}
