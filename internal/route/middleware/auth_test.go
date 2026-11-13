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

// fakeAuth issues tokens that always claim "activated"; dbStatus is the
// status currently stored for the user.
type fakeAuth struct {
	auth.Auth
	dbStatus authV1.AccountStatus
	allow    bool
}

func (f fakeAuth) CheckPermission(_ context.Context, _, _, _ string) bool {
	return f.allow
}

func (f fakeAuth) ParseToken(_ context.Context, _ string) (*auth.UserTokenClaims, error) {
	return &auth.UserTokenClaims{UserId: "user-1", Role: authV1.Role_user.String(), Status: int(authV1.AccountStatus_activated)}, nil
}

func (f fakeAuth) GetUserStatus(_ context.Context, _ string) authV1.AccountStatus {
	return f.dbStatus
}

func newTestController(dbStatus authV1.AccountStatus) *Controller {
	return New(nil, fakeAuth{dbStatus: dbStatus}, logger.New("error", "console", "test"))
}

func serve(handlers ...gin.HandlerFunc) (code int, reachedHandler bool) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handlers = append(handlers, func(c *gin.Context) {
		reachedHandler = true
		c.Status(http.StatusOK)
	})
	r.GET("/x", handlers...)

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Code, reachedHandler
}

func TestAuthToken_BlockedAccountIsStopped(t *testing.T) {
	// token still says "activated" but the account was blocked afterwards
	code, reached := serve(newTestController(authV1.AccountStatus_blocked).AuthToken)
	if reached {
		t.Fatal("blocked account must not reach the handler")
	}
	if code == http.StatusOK {
		t.Fatalf("want error status, got %d", code)
	}
}

func TestAuthToken_DeletedUserIsStopped(t *testing.T) {
	code, reached := serve(newTestController(-1).AuthToken)
	if reached || code != http.StatusUnauthorized {
		t.Fatalf("deleted user must get 401, code=%d reached=%v", code, reached)
	}
}

func TestAuthToken_ActiveAccountPasses(t *testing.T) {
	code, reached := serve(newTestController(authV1.AccountStatus_activated).AuthToken)
	if !reached || code != http.StatusOK {
		t.Fatalf("active account should pass, code=%d reached=%v", code, reached)
	}
}

func TestAuthToken_UsesCurrentStatusNotTokenClaim(t *testing.T) {
	// token claims "activated", DB says still "initialized": activated-only routes must refuse
	m := newTestController(authV1.AccountStatus_initialized)
	_, reached := serve(m.AuthToken, m.CheckAccountStatusActivated)
	if reached {
		t.Fatal("status from the token claim must not be trusted")
	}
}

func TestRecover_ReturnsJsonErrorForNonErrorPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	m := newTestController(authV1.AccountStatus_activated)
	r := gin.New()
	r.Use(m.RecoverMiddleware)
	r.GET("/panic", func(c *gin.Context) { panic("boom") })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/panic", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("want 500, got %d", w.Code)
	}
	if w.Body.Len() == 0 {
		t.Fatal("want a JSON error body")
	}
}

func TestRequest_ClientCannotSetRequestId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	m := newTestController(authV1.AccountStatus_activated)
	r := gin.New()
	var got string
	r.GET("/x", m.Request, func(c *gin.Context) { got = c.GetHeader(HeaderRequestId) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set(HeaderRequestId, "client-chosen")
	r.ServeHTTP(httptest.NewRecorder(), req)

	if got == "" || got == "client-chosen" {
		t.Fatalf("request id must be server-generated, got %q", got)
	}
}

func TestPermission_DeniedIs403(t *testing.T) {
	m := New(nil, fakeAuth{dbStatus: authV1.AccountStatus_activated, allow: false}, logger.New("error", "console", "test"))
	code, reached := serve(m.AuthToken, m.Permission)
	if reached || code != http.StatusForbidden {
		t.Fatalf("want 403, code=%d reached=%v", code, reached)
	}
}

func TestCheckAccountStatusActivated_NotActivatedIs403(t *testing.T) {
	m := newTestController(authV1.AccountStatus_initialized)
	code, _ := serve(m.AuthToken, m.CheckAccountStatusActivated)
	if code != http.StatusForbidden {
		t.Fatalf("want 403, got %d", code)
	}
}

func TestBlockedAccountIs403(t *testing.T) {
	code, _ := serve(newTestController(authV1.AccountStatus_blocked).AuthToken)
	if code != http.StatusForbidden {
		t.Fatalf("want 403, got %d", code)
	}
}
