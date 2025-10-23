package middleware

import (
	"context"
	"net/http"

	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/http/response"
	"github.com/rxmy43/support-platform/internal/modules/auth"
)

func AuthMiddleware(authManager *auth.AuthManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				response.ToJSON(w, r, apperror.Unauthorized("invalid token", apperror.CodeTokenInvalid))
			}

			if session, exists := authManager.GetSession(cookie.Value); exists {
				ctx := context.WithValue(r.Context(), "user_id", session.UserID)
				ctx = context.WithValue(ctx, "user_role", session.Role)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			response.ToJSON(w, r, apperror.Unauthorized("unauthorized", apperror.CodeUnauthorizedOperation))
		})
	}
}
