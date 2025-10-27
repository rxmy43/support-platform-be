package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/http/middleware"
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

var decoder = schema.NewDecoder()

func init() {
	decoder.IgnoreUnknownKeys(true)
}

func (h *SupportHandler) PaymentCallback(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Failed to read request body:", err)
		response.ToJSON(w, r, apperror.BadRequest("cannot read body", apperror.CodeInvalidRequestJSONFormat))
		return
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := r.ParseForm(); err != nil {
		log.Println("ParseForm error:", err)
		response.ToJSON(w, r, apperror.BadRequest("invalid form payload", apperror.CodeInvalidRequestJSONFormat))
		return
	}

	var pc support.PaymentCallbackRequest
	if err := decoder.Decode(&pc, r.PostForm); err != nil {
		log.Println("ParseForm error:", err)
		response.ToJSON(w, r, apperror.BadRequest("invalid form payload", apperror.CodeUnknown))
		return
	}

	status, appErr := h.supportService.PaymentCallback(r.Context(), pc)
	if appErr != nil {
		response.ToJSON(w, r, err)
		return
	}

	response.ToJSON(w, r, status)

}

func (h *SupportHandler) GetBestSupporters(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	var cursor *uint

	if cursorStr != "" {
		parsed, err := strconv.ParseUint(cursorStr, 10, 64)
		if err != nil {
			response.ToJSON(w, r, apperror.BadRequest("invalid cursor", apperror.CodeUnknown))
			return
		}
		temp := uint(parsed)
		cursor = &temp
	}

	userID := middleware.GetUserID(r.Context())
	userRole := middleware.GetUserRole(r.Context())

	if userRole != "creator" {
		response.ToJSON(w, r, apperror.Unauthorized("invalid role", apperror.CodeUnauthorizedOperation))
		return
	}

	bestSupports, nextCursor, err := h.supportService.GetSupporters(r.Context(), cursor, *userID)
	if err != nil {
		response.ToJSON(w, r, err)
		return
	}

	data := make([]any, len(bestSupports))
	for i, b := range bestSupports {
		data[i] = b
	}

	resp := response.SuccessPaginateResponse{
		Status:     response.StatusSuccess,
		Data:       data,
		NextCursor: nextCursor,
	}

	response.ToJSON(w, r, resp)
}

func (h *SupportHandler) GetFanSpending(w http.ResponseWriter, r *http.Request) {
	userID := *middleware.GetUserID(r.Context())
	userRole := middleware.GetUserRole(r.Context())

	if userRole != "fan" {
		response.ToJSON(w, r, apperror.Unauthorized("invalid role", apperror.CodeUnauthorizedOperation))
		return
	}

	spending, err := h.supportService.GetFanSpending(r.Context(), userID)
	if err != nil {
		response.ToJSON(w, r, err)
		return
	}

	response.ToJSON(w, r, spending)
}

func (h *SupportHandler) GetFanSpendingHistory(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	var cursor *uint

	if cursorStr != "" {
		parsed, err := strconv.ParseUint(cursorStr, 10, 64)
		if err != nil {
			response.ToJSON(w, r, apperror.BadRequest("invalid cursor", apperror.CodeUnknown))
			return
		}

		temp := uint(parsed)
		cursor = &temp
	}

	var userID *uint
	if userRole := middleware.GetUserRole(r.Context()); userRole != "" && userRole == "fan" {
		userID = middleware.GetUserID(r.Context())
	}

	histories, nextCursor, err := h.supportService.GetFanSupportHistory(r.Context(), cursor, userID)
	if err != nil {
		response.ToJSON(w, r, err)
		return
	}

	data := make([]any, len(histories))
	for i, h := range histories {
		data[i] = h
	}

	resp := response.SuccessPaginateResponse{
		Status:     response.StatusSuccess,
		Data:       data,
		NextCursor: nextCursor,
	}

	response.ToJSON(w, r, resp)
}
