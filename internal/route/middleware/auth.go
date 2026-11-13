package middleware

import (
	"github.com/Template7/backend/internal/route/response"
	"github.com/Template7/backend/internal/t7Error"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/gin-gonic/gin"
)

// Permission checks the caller's role against the RBAC policy for this path.
func (m *Controller) Permission(c *gin.Context) {
	role, ok := currentRole(c)
	if !ok {
		response.Error(c, t7Error.InvalidToken.WithDetail("no role in context"))
		return
	}
	if !m.authSvc.CheckPermission(c, role, c.Request.URL.Path, c.Request.Method) {
		response.Error(c, t7Error.Forbidden)
		return
	}
	c.Next()
}

// AuthToken validates the bearer token and puts the user id, role and current
// account status on the context.
func (m *Controller) AuthToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		response.Error(c, t7Error.InvalidToken.WithDetail("empty token"))
		return
	}
	claims, err := m.authSvc.ParseToken(c, token)
	if err != nil {
		response.Error(c, err)
		return
	}

	// account status can change after the token was issued (activated, blocked, deleted),
	// so read the current value instead of trusting the claim
	status := m.authSvc.GetUserStatus(c, claims.UserId)
	if status < 0 {
		response.Error(c, t7Error.InvalidToken.WithDetail("user not found"))
		return
	}
	if status == authV1.AccountStatus_blocked {
		response.Error(c, t7Error.BlockedAccount)
		return
	}

	c.Set(UserId, claims.UserId)
	c.Set(Role, claims.Role)
	c.Set(Status, int(status))
	c.Next()
}

// AuthUserWallet allows the request only if the :walletId belongs to the caller.
// A wallet the caller doesn't own is reported as not found, so ids can't be probed.
func (m *Controller) AuthUserWallet(c *gin.Context) {
	walletId := c.Param("walletId")
	if walletId == "" {
		response.BadRequest(c, "empty wallet id")
		return
	}
	userId, ok := CurrentUserId(c)
	if !ok {
		response.Error(c, t7Error.InvalidToken.WithDetail("no user id in context"))
		return
	}

	for _, uw := range m.userSvc.GetUserWallets(c, userId) {
		if uw.Id == walletId {
			c.Next()
			return
		}
	}
	response.Error(c, t7Error.WalletNotFound)
}

// CheckAccountStatusActivated allows only activated accounts.
func (m *Controller) CheckAccountStatusActivated(c *gin.Context) {
	m.requireStatus(c, func(s int) bool { return s == int(authV1.AccountStatus_activated) })
}

// CheckAccountStatusInitialized allows initialized and activated accounts.
func (m *Controller) CheckAccountStatusInitialized(c *gin.Context) {
	m.requireStatus(c, func(s int) bool { return s >= int(authV1.AccountStatus_initialized) })
}

func (m *Controller) requireStatus(c *gin.Context, allowed func(int) bool) {
	status, ok := currentStatus(c)
	if !ok {
		response.Error(c, t7Error.InvalidToken.WithDetail("no account status in context"))
		return
	}
	if !allowed(status) {
		response.Error(c, t7Error.AccountNotActivated)
		return
	}
	c.Next()
}
