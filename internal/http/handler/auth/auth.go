package auth

import (
	"encoding/json"
	"net/http"

	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/http/response"
	"github.com/rxmy43/support-platform/internal/modules/auth"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) GenerateOTP(w http.ResponseWriter, r *http.Request) {
	var req auth.GenerateOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ToJSON(w, r, apperror.BadRequest("invalid json format", apperror.CodeInvalidRequestJSONFormat))
		return
	}

	otp, err := h.authService.GenerateOTP(r.Context(), req)
	if err != nil {
		response.ToJSON(w, r, err)
		return
	}

	resp := response.SuccessResponse{
		Status:  "success",
		Message: "generated otp",
		Data:    otp,
	}

	response.ToJSON(w, r, resp)
}

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req auth.VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ToJSON(w, r, apperror.BadRequest("invalid json format", apperror.CodeInvalidRequestJSONFormat))
		return
	}

	user, err := h.authService.VerifyOTP(r.Context(), req)
	if err != nil {
		response.ToJSON(w, r, err)
		return
	}

	response.ToJSON(w, r, user)
}
