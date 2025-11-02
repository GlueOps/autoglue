package httpmiddleware

import (
	"net/http"

	"github.com/glueops/autoglue/internal/utils"
)

func RequireRole(minRole string) func(http.Handler) http.Handler {
	// order: owner > admin > member
	rank := map[string]int{
		"role:member":    1,
		"role:admin":     2,
		"role:owner":     3,
		"org:machine":    2,
		"org:machine:ro": 1,
	}
	need := map[string]bool{
		"member": true, "admin": true, "owner": true,
	}
	if !need[minRole] {
		minRole = "member"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles, ok := RolesFrom(r.Context())
			if !ok || len(roles) == 0 {
				utils.WriteError(w, http.StatusForbidden, "forbidden", "no roles in context")
				return
			}
			max := 0
			for _, ro := range roles {
				if rank[ro] > max {
					max = rank[ro]
				}
			}
			if max < rank["role:"+minRole] {
				utils.WriteError(w, http.StatusForbidden, "forbidden", "insufficient role")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
