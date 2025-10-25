package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/http/response"
)

type ctxKey string

const (
	userIDKey   ctxKey = "userID"
	userRoleKey ctxKey = "userRole"
)

func UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		userRole := r.Header.Get("X-User-Role")

		if userID == "" {
			response.ToJSON(w, r, apperror.Unauthorized("Missing Header User ID", apperror.CodeUnauthorizedOperation))
			return
		}

		if userRole == "" {
			response.ToJSON(w, r, apperror.Unauthorized("Missing Header User Role", apperror.CodeUnauthorizedOperation))
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, userRoleKey, userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) *uint {
	if val, ok := ctx.Value(userIDKey).(string); ok {
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil
		}

		temp := uint(parsed)
		return &temp
	}
	return nil
}

func GetUserRole(ctx context.Context) string {
	if val, ok := ctx.Value(userRoleKey).(string); ok {
		return val
	}
	return ""
}
