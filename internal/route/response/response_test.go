package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Template7/backend/api/types"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
	logger.New("error", "console", "test")
}

func run(t *testing.T, h gin.HandlerFunc) (*httptest.ResponseRecorder, types.HttpRespBase, bool) {
	t.Helper()
	nextCalled := false
	r := gin.New()
	r.GET("/x", func(c *gin.Context) {
		c.Request.Header.Set(HeaderRequestId, "req-1")
		h(c)
	}, func(c *gin.Context) { nextCalled = true })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))

	var body types.HttpRespBase
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("body is not JSON: %q", w.Body.String())
	}
	return w, body, nextCalled
}

func TestError_UsesStatusAndCodeOfT7Error(t *testing.T) {
	cases := []struct {
		err    *t7Error.Error
		status int
	}{
		{t7Error.InvalidBody, http.StatusBadRequest},
		{t7Error.InvalidToken, http.StatusUnauthorized},
		{t7Error.Forbidden, http.StatusForbidden},
		{t7Error.WalletNotFound, http.StatusNotFound},
		{t7Error.UserAlreadyExist, http.StatusConflict},
		{t7Error.DbOperationFail, http.StatusInternalServerError},
	}
	for _, tc := range cases {
		w, body, next := run(t, func(c *gin.Context) { Error(c, tc.err.WithDetail("internal detail")) })
		if w.Code != tc.status || body.Code != int(tc.err.Code) || body.Message != tc.err.Message {
			t.Errorf("%s: got %d %+v, want %d code %d", tc.err.Message, w.Code, body, tc.status, tc.err.Code)
		}
		if body.RequestId != "req-1" {
			t.Errorf("%s: request id = %q", tc.err.Message, body.RequestId)
		}
		if next {
			t.Errorf("%s: Error must abort the chain", tc.err.Message)
		}
	}
}

func TestError_UnknownErrorIs500(t *testing.T) {
	w, body, _ := run(t, func(c *gin.Context) { Error(c, errors.New("boom")) })
	if w.Code != http.StatusInternalServerError || body.Code != int(t7Error.Unknown.Code) {
		t.Fatalf("got %d %+v", w.Code, body)
	}
}

func TestError_DetailIsNotLeaked(t *testing.T) {
	w, _, _ := run(t, func(c *gin.Context) { Error(c, t7Error.DbOperationFail.WithDetail("SELECT secret FROM x")) })
	if got := w.Body.String(); strings.Contains(got, "secret") {
		t.Fatalf("detail leaked into response: %s", got)
	}
}
