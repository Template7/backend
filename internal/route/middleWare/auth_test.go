package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/common/logger"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/gin-gonic/gin"
)

type fakeAuth struct {
	auth.Auth
	status authV1.AccountStatus
}

func (f fakeAuth) ParseToken(_ context.Context, _ string) (*auth.UserTokenClaims, error) {
	return &auth.UserTokenClaims{UserId: "user-1", Role: authV1.Role_user.String(), Status: int(f.status)}, nil
}

func serve(status authV1.AccountStatus) (code int, reachedHandler bool) {
	gin.SetMode(gin.TestMode)
	m := New(nil, fakeAuth{status: status}, logger.New("error", "console", "test"))

	r := gin.New()
	r.GET("/x", m.AuthToken, func(c *gin.Context) {
		reachedHandler = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Code, reachedHandler
}

func TestAuthToken_BlockedAccountIsStopped(t *testing.T) {
	code, reached := serve(authV1.AccountStatus_blocked)
	if reached {
		t.Fatal("blocked account must not reach the handler")
	}
	if code == http.StatusOK {
		t.Fatalf("want error status, got %d", code)
	}
}

func TestAuthToken_ActiveAccountPasses(t *testing.T) {
	code, reached := serve(authV1.AccountStatus_activated)
	if !reached || code != http.StatusOK {
		t.Fatalf("active account should pass, code=%d reached=%v", code, reached)
	}
}
