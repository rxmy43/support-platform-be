package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/config"
)

type PaginationMeta struct {
	Page       int `json:"page,omitempty"`
	Limit      int `json:"per_page,omitempty"`
	TotalItems int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

type ResponseStatus string

const (
	StatusSuccess ResponseStatus = "success"
	StatusError   ResponseStatus = "error"
)

type SuccessResponse struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
	Data    any            `json:"data,omitempty"`
}

type SuccessPaginateResponse struct {
	Status     ResponseStatus `json:"status"`
	Data       []any          `json:"data"`
	NextCursor *uint          `json:"next_cursor"`
}

type ErrorResponse struct {
	Status        ResponseStatus        `json:"status"`
	Code          apperror.ErrorCode    `json:"code"`
	Message       string                `json:"message"`
	NotFoundField string                `json:"not_found_field,omitempty"`
	FieldErrors   []apperror.FieldError `json:"field_errors,omitempty"`
	Details       string                `json:"details,omitempty"`
}

func ToJSON(w http.ResponseWriter, r *http.Request, payload any) {
	w.Header().Set("Content-Type", "application/json")

	if err, ok := payload.(*apperror.AppError); ok {
		statusCode := err.HTTPStatus()
		w.WriteHeader(statusCode)

		resp := ErrorResponse{
			Status:  StatusError,
			Code:    err.Code,
			Message: err.Message,
		}

		if statusCode == http.StatusUnprocessableEntity && len(err.FieldErrors) > 0 {
			hasInternal := false
			message := ""
			expect := ""
			for _, f := range err.FieldErrors {
				if f.IsInternalError {
					hasInternal = true
					message = f.Message
					expect = f.Expect
					break
				}
			}

			if hasInternal {
				w.WriteHeader(500)
				if config.Load().Env == "development" {
					resp.Message = message
					resp.Details = fmt.Sprintf("expect : %s", expect)
				} else {
					resp.Message = "Internal Server Error, please contact support."
				}
			} else {
				w.WriteHeader(statusCode)
				resp.FieldErrors = err.FieldErrors
			}
		}

		// Add not found field if present
		if statusCode == http.StatusNotFound && err.NotFoundField != "" {
			resp.NotFoundField = err.NotFoundField
		}

		// Add error details for internal server errors (useful for debugging)
		if statusCode >= 500 && err.Error() != err.Message {
			if config.Load().Env == "development" {
				resp.Message = err.Message
				resp.Details = err.Error()
			} else {
				resp.Message = "Internal Server Error, please contact support."
			}
		}

		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	// Handle regular Go errors by converting them to AppError
	if err, ok := payload.(error); ok {
		appErr := apperror.New(
			apperror.CodeInternalServerError,
			500,
			"An unexpected error occurred",
		).WithCause(err)

		w.WriteHeader(appErr.HTTPStatus())
		resp := ErrorResponse{
			Status:  StatusError,
			Code:    appErr.Code,
			Message: appErr.Message,
		}

		// Include error details for debugging
		if config.Load().Env == "development" {
			resp.Details = err.Error()
		}

		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	// ==== SUCCESS HANDLING ====
	w.WriteHeader(http.StatusOK)

	switch v := payload.(type) {
	case string:
		resp := SuccessResponse{
			Status:  StatusSuccess,
			Message: v,
		}
		_ = json.NewEncoder(w).Encode(resp)

	case SuccessResponse, SuccessPaginateResponse:
		// Direct struct
		_ = json.NewEncoder(w).Encode(v)
	case nil:
		resp := SuccessResponse{
			Status:  StatusSuccess,
			Message: "Operation completed successfully",
		}
		_ = json.NewEncoder(w).Encode(resp)

	default:
		// Success with data
		resp := SuccessResponse{
			Status:  "success",
			Message: "Operation completed successfully",
			Data:    v,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}
}
