package auth

import (
	"context"
	"testing"
	"time"

	"github.com/Template7/backend/internal/db"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/models"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret-0123456789abcdef0123456789"

type fakeDb struct{ db.Client }

func (fakeDb) GetUserById(_ context.Context, userId string) (models.User, error) {
	return models.User{Id: userId, Status: authV1.AccountStatus_activated}, nil
}

func newTestService(secret string) *service {
	return &service{
		jwtSecret: []byte(secret),
		db:        fakeDb{},
		log:       logger.New("error", "console", "test"),
	}
}

func TestToken_RoundTrip(t *testing.T) {
	s := newTestService(testSecret)

	tk, err := s.genUserToken(context.Background(), "user-1", authV1.Role_user)
	if err != nil {
		t.Fatalf("gen: %v", err)
	}
	claims, err := s.ParseToken(context.Background(), tk)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if claims.UserId != "user-1" || claims.Role != authV1.Role_user.String() || claims.Status != int(authV1.AccountStatus_activated) {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestToken_RejectsOtherSecret(t *testing.T) {
	tk, err := newTestService("another-secret-0123456789abcdef012345").genUserToken(context.Background(), "user-1", authV1.Role_admin)
	if err != nil {
		t.Fatalf("gen: %v", err)
	}
	if _, err := newTestService(testSecret).ParseToken(context.Background(), tk); err == nil {
		t.Fatal("token signed with a different secret must be rejected")
	}
}

func TestToken_RejectsNoneAlg(t *testing.T) {
	claims := UserTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))},
		UserId:           "user-1",
		Role:             authV1.Role_admin.String(),
	}
	tk, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if _, err := newTestService(testSecret).ParseToken(context.Background(), tk); err == nil {
		t.Fatal("alg=none token must be rejected")
	}
}

func TestToken_RejectsExpired(t *testing.T) {
	claims := UserTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute))},
		UserId:           "user-1",
	}
	tk, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testSecret))
	if _, err := newTestService(testSecret).ParseToken(context.Background(), tk); err == nil {
		t.Fatal("expired token must be rejected")
	}
}

func TestToken_RejectsMissingExpiry(t *testing.T) {
	tk, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, UserTokenClaims{UserId: "user-1"}).SignedString([]byte(testSecret))
	if _, err := newTestService(testSecret).ParseToken(context.Background(), tk); err == nil {
		t.Fatal("token without exp must be rejected")
	}
}
