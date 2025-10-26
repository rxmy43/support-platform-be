package support

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/http/response"
	"github.com/rxmy43/support-platform/internal/modules/support"
)

type SupportHandler struct {
	supportService *support.SupportService
}

func NewSupportHandler(supportService *support.SupportService) *SupportHandler {
	return &SupportHandler{
		supportService: supportService,
	}
}

func (h *SupportHandler) Donate(w http.ResponseWriter, r *http.Request) {
	var req support.DonationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ToJSON(w, r, apperror.BadRequest("invalida request json format", apperror.CodeInvalidRequestJSONFormat))
		return
	}

	paymentURL, err := h.supportService.Donate(r.Context(), req)
	if err != nil {
		response.ToJSON(w, r, err)
		return
	}

	response.ToJSON(w, r, map[string]string{"payment_url": paymentURL})
}

func (h *SupportHandler) PaymentCallback(w http.ResponseWriter, r *http.Request) {
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.ToJSON(w, r, apperror.BadRequest("invalid request json format", apperror.CodeInvalidRequestJSONFormat))
		fmt.Printf("%+v\n", payload)
		return
	}

	response.ToJSON(w, r, payload)

}
