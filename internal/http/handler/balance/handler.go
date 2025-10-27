package balance

import (
	"net/http"

	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/http/middleware"
	"github.com/rxmy43/support-platform/internal/http/response"
	"github.com/rxmy43/support-platform/internal/modules/balance"
)

type BalanceHandler struct {
	balanceService *balance.BalanceService
}

func NewBalanceHandler(balanceService *balance.BalanceService) *BalanceHandler {
	return &BalanceHandler{
		balanceService: balanceService,
	}
}

func (h *BalanceHandler) GetCreatorBalance(w http.ResponseWriter, r *http.Request) {
	userID := *middleware.GetUserID(r.Context())
	userRole := middleware.GetUserRole(r.Context())

	if userRole != "creator" {
		response.ToJSON(w, r, apperror.Unauthorized("invalid role", apperror.CodeUnauthorizedOperation))
		return
	}

	amount, err := h.balanceService.GetCreatorBalance(r.Context(), userID)
	if err != nil {
		response.ToJSON(w, r, err)
		return
	}

	response.ToJSON(w, r, amount)
}
