package auth

import (
	"context"
	"sort"

	"github.com/Template7/backend/internal/t7Error"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
)

// Rule allows a role to call one route.
type Rule struct {
	Path   string // gin path pattern, e.g. /api/v1/wallets/:walletId
	Method string
}

// SyncRolePolicies makes the stored policies of role exactly rules: missing ones
// are added and stale ones removed. It is idempotent, so every instance can run it
// at startup.
func (s *service) SyncRolePolicies(ctx context.Context, role authV1.Role, rules []Rule) error {
	want := map[[2]string]bool{}
	for _, r := range rules {
		want[[2]string{r.Path, r.Method}] = true
	}

	existing, err := s.core.GetFilteredPolicy(0, role.String())
	if err != nil {
		return t7Error.DbOperationFail.WithDetail(err.Error())
	}
	have := map[[2]string]bool{}
	var stale [][]string
	for _, p := range existing {
		if len(p) < 3 {
			continue
		}
		key := [2]string{p[1], p[2]}
		have[key] = true
		if !want[key] {
			stale = append(stale, p)
		}
	}

	var missing [][]string
	for key := range want {
		if !have[key] {
			missing = append(missing, []string{role.String(), key[0], key[1]})
		}
	}
	sort.Slice(missing, func(i, j int) bool { return missing[i][1]+missing[i][2] < missing[j][1]+missing[j][2] })

	log := s.log.WithContext(ctx).With("role", role.String())
	if len(stale) > 0 {
		if _, err := s.core.RemovePolicies(stale); err != nil {
			return t7Error.DbOperationFail.WithDetail(err.Error())
		}
		log.With("policies", stale).Info("removed stale policies")
	}
	if len(missing) > 0 {
		if _, err := s.core.AddPolicies(missing); err != nil {
			return t7Error.DbOperationFail.WithDetail(err.Error())
		}
		log.With("policies", missing).Info("added policies")
	}
	return nil
}
