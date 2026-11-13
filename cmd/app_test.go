package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/cache/cachetest"
	"github.com/Template7/backend/internal/config"
	"github.com/Template7/backend/internal/db/dbtest"
	"github.com/Template7/backend/internal/route/handler"
	"github.com/Template7/backend/internal/route/middleware"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/backend/internal/user"
	"github.com/Template7/backend/internal/wallet"
	"github.com/Template7/common/logger"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
	"github.com/gin-gonic/gin"
)

// testApp wires the real router and services on in-memory SQLite, real casbin
// and an in-memory cache.
type testApp struct {
	t      *testing.T
	router http.Handler
	users  user.Service
}

func newTestApp(t *testing.T) *testApp {
	t.Helper()
	gin.SetMode(gin.TestMode)
	log := logger.New("error", "console", "test")

	cfg := &config.Config{Env: "test"}
	cfg.Auth.RbacModelPath = "../config/rbac_model.conf"
	cfg.Auth.JwtSecret = "e2e-secret-0123456789abcdef0123456789"

	client, core := dbtest.New(t)
	authSvc := auth.New(client, core, log, cfg)
	userSvc := user.New(authSvc, client, cachetest.New(), log)
	walletSvc := wallet.New(client, log)

	app := NewApp(
		handler.NewAuthController(authSvc, log),
		handler.NewUserController(userSvc, log),
		handler.NewWalletController(walletSvc, log),
		middleware.New(userSvc, authSvc, log),
		cfg, log)
	return &testApp{t: t, router: app.SetupRoutes(), users: userSvc}
}

type resp struct {
	status int
	Code   int             `json:"code"`
	Data   json.RawMessage `json:"data"`
}

func (a *testApp) do(method, path, token string, body interface{}) resp {
	a.t.Helper()
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	w := httptest.NewRecorder()
	a.router.ServeHTTP(w, req)

	var r resp
	if err := json.Unmarshal(w.Body.Bytes(), &r); err != nil {
		a.t.Fatalf("%s %s: body is not JSON: %q", method, path, w.Body.String())
	}
	r.status = w.Code
	return r
}

func (a *testApp) expect(r resp, status int, code int) {
	a.t.Helper()
	if r.status != status || r.Code != code {
		a.t.Fatalf("got %d / code %d (%s), want %d / code %d", r.status, r.Code, r.Data, status, code)
	}
}

func (a *testApp) login(username, password string) string {
	a.t.Helper()
	r := a.do(http.MethodPost, "/api/v1/login/native", "", map[string]string{"username": username, "password": password})
	a.expect(r, http.StatusOK, 3000)
	var d struct{ Token string }
	_ = json.Unmarshal(r.Data, &d)
	return d.Token
}

// createUser creates (and optionally activates) a user through the admin API.
func (a *testApp) createUser(adminToken, username string, activate bool) {
	a.t.Helper()
	r := a.do(http.MethodPost, "/admin/v1/user", adminToken, map[string]string{
		"username": username, "password": "pass-" + username, "role": "user",
	})
	a.expect(r, http.StatusOK, 3000)
	var d struct{ UserId, ActivationCode string }
	_ = json.Unmarshal(r.Data, &d)

	if activate {
		r = a.do(http.MethodPost, "/admin/v1/users/"+d.UserId+"/activate", adminToken, map[string]string{"activationCode": d.ActivationCode})
		a.expect(r, http.StatusOK, 3000)
	}
}

func (a *testApp) walletId(token string) string {
	a.t.Helper()
	r := a.do(http.MethodGet, "/api/v1/user/wallets", token, nil)
	a.expect(r, http.StatusOK, 3000)
	var ws []struct{ Id string }
	_ = json.Unmarshal(r.Data, &ws)
	if len(ws) != 1 {
		a.t.Fatalf("wallets = %s", r.Data)
	}
	return ws[0].Id
}

func (a *testApp) usd(token, walletId string) string {
	a.t.Helper()
	r := a.do(http.MethodGet, "/api/v1/wallets/"+walletId, token, nil)
	a.expect(r, http.StatusOK, 3000)
	var w struct {
		Balances []struct{ Currency, Amount string }
	}
	_ = json.Unmarshal(r.Data, &w)
	for _, b := range w.Balances {
		if b.Currency == "usd" {
			return b.Amount
		}
	}
	a.t.Fatalf("no usd balance in %s", r.Data)
	return ""
}

func code(e *t7Error.Error) int { return int(e.Code) }

func TestEndToEnd(t *testing.T) {
	a := newTestApp(t)
	ctx := context.Background()

	// bootstrap an admin directly through the service
	adminId, actCode, err := a.users.CreateUser(ctx, &userV1.CreateUserRequest{Username: "admin", Password: "pass-admin", Role: authV1.Role_admin})
	if err != nil {
		t.Fatal(err)
	}
	if err := a.users.ActivateUser(ctx, adminId, actCode); err != nil {
		t.Fatal(err)
	}
	admin := a.login("admin", "pass-admin")

	a.createUser(admin, "alice", true)
	a.createUser(admin, "bob", true)
	a.createUser(admin, "carol", false) // never activated

	alice := a.login("alice", "pass-alice")
	bob := a.login("bob", "pass-bob")
	carol := a.login("carol", "pass-carol")
	aliceWallet, bobWallet := a.walletId(alice), a.walletId(bob)

	t.Run("login failures look the same", func(t *testing.T) {
		a.expect(a.do(http.MethodPost, "/api/v1/login/native", "", map[string]string{"username": "alice", "password": "wrong"}), http.StatusUnauthorized, code(t7Error.SignInFail))
		a.expect(a.do(http.MethodPost, "/api/v1/login/native", "", map[string]string{"username": "nobody", "password": "x"}), http.StatusUnauthorized, code(t7Error.SignInFail))
	})

	t.Run("no or bad token is 401", func(t *testing.T) {
		a.expect(a.do(http.MethodGet, "/api/v1/user/info", "", nil), http.StatusUnauthorized, code(t7Error.InvalidToken))
		a.expect(a.do(http.MethodGet, "/api/v1/user/info", "garbage", nil), http.StatusUnauthorized, code(t7Error.InvalidToken))
	})

	t.Run("user on admin route is 403", func(t *testing.T) {
		a.expect(a.do(http.MethodPost, "/admin/v1/user", alice, map[string]string{"username": "x", "password": "y", "role": "user"}), http.StatusForbidden, code(t7Error.Forbidden))
	})

	t.Run("money flow", func(t *testing.T) {
		a.expect(a.do(http.MethodPost, "/api/v1/wallets/"+aliceWallet+"/deposit", alice, map[string]interface{}{"currency": "usd", "amount": 100}), http.StatusOK, 3000)
		a.expect(a.do(http.MethodPost, "/api/v1/wallets/"+aliceWallet+"/withdraw", alice, map[string]interface{}{"currency": "usd", "amount": 30}), http.StatusOK, 3000)
		a.expect(a.do(http.MethodPost, "/api/v1/transfer", alice, map[string]interface{}{
			"fromWalletId": aliceWallet, "toWalletId": bobWallet, "currency": "usd", "amount": 50,
		}), http.StatusOK, 3000)

		if got := a.usd(alice, aliceWallet); got != "20" {
			t.Fatalf("alice usd = %s, want 20", got)
		}
		if got := a.usd(bob, bobWallet); got != "50" {
			t.Fatalf("bob usd = %s, want 50", got)
		}

		r := a.do(http.MethodGet, "/api/v1/wallets/"+aliceWallet+"/history", alice, nil)
		a.expect(r, http.StatusOK, 3000)
		var h []struct{ Direction string }
		_ = json.Unmarshal(r.Data, &h)
		if len(h) != 3 {
			t.Fatalf("history = %s", r.Data)
		}
	})

	t.Run("insufficient balance is 400", func(t *testing.T) {
		a.expect(a.do(http.MethodPost, "/api/v1/wallets/"+aliceWallet+"/withdraw", alice, map[string]interface{}{"currency": "usd", "amount": 1000}),
			http.StatusBadRequest, code(t7Error.InsufficientBalance))
	})

	t.Run("someone else's wallet is 404", func(t *testing.T) {
		nf := code(t7Error.WalletNotFound)
		a.expect(a.do(http.MethodGet, "/api/v1/wallets/"+bobWallet, alice, nil), http.StatusNotFound, nf)
		a.expect(a.do(http.MethodPost, "/api/v1/wallets/"+bobWallet+"/withdraw", alice, map[string]interface{}{"currency": "usd", "amount": 1}), http.StatusNotFound, nf)
		a.expect(a.do(http.MethodGet, "/api/v1/wallets/"+bobWallet+"/history", alice, nil), http.StatusNotFound, nf)
		a.expect(a.do(http.MethodPost, "/api/v1/transfer", alice, map[string]interface{}{
			"fromWalletId": bobWallet, "toWalletId": aliceWallet, "currency": "usd", "amount": 1,
		}), http.StatusNotFound, nf)
		if got := a.usd(bob, bobWallet); got != "50" {
			t.Fatalf("bob usd changed to %s", got)
		}
	})

	t.Run("not activated is 403 for money routes", func(t *testing.T) {
		carolWallet := a.walletId(carol) // initialized users may read
		a.expect(a.do(http.MethodPost, "/api/v1/wallets/"+carolWallet+"/deposit", carol, map[string]interface{}{"currency": "usd", "amount": 1}),
			http.StatusForbidden, code(t7Error.AccountNotActivated))
	})

	t.Run("deleted user's token stops working", func(t *testing.T) {
		r := a.do(http.MethodGet, "/api/v1/user/info", bob, nil)
		a.expect(r, http.StatusOK, 3000)
		var info struct{ UserId string }
		_ = json.Unmarshal(r.Data, &info)

		a.expect(a.do(http.MethodDelete, "/admin/v1/users/"+info.UserId, admin, nil), http.StatusOK, 3000)
		a.expect(a.do(http.MethodGet, "/api/v1/user/info", bob, nil), http.StatusUnauthorized, code(t7Error.InvalidToken))
		a.expect(a.do(http.MethodDelete, "/admin/v1/users/"+info.UserId, admin, nil), http.StatusNotFound, code(t7Error.UserNotfound))
	})
}
