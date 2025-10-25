package support

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/config"
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

type CallbackPayload struct {
	MerchantCode     string `json:"merchantCode"`
	Amount           string `json:"amount"`
	MerchantOrderId  string `json:"merchantOrderId"`
	ProductDetail    string `json:"productDetail"`
	AdditionalParam  string `json:"additionalParam"`
	PaymentCode      string `json:"paymentCode"`
	ResultCode       string `json:"resultCode"`
	MerchantUserId   string `json:"merchantUserId"`
	Reference        string `json:"reference"`
	Signature        string `json:"signature"`
	PublisherOrderId string `json:"publisherOrderId"`
	SettlementDate   string `json:"settlementDate"`
	IssuerCode       string `json:"issuerCode"`
}

func verifySignature(payload CallbackPayload) bool {
	apiKey := config.Load().DuitkuAPIKey
	raw := payload.MerchantCode + payload.Amount + payload.MerchantOrderId + apiKey
	hash := md5.Sum([]byte(raw))
	expectedSig := hex.EncodeToString(hash[:])
	return expectedSig == payload.Signature
}

func (h *SupportHandler) PaymentCallback(w http.ResponseWriter, r *http.Request) {
	var payload CallbackPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.ToJSON(w, r, apperror.BadRequest("invalid request json format", apperror.CodeInvalidRequestJSONFormat))
		return
	}

	if !verifySignature(payload) {
		response.ToJSON(w, r, apperror.Unauthorized("invalid signature", apperror.CodeUnknown))
		return
	}

	// Convert amount string to int
	amount, _ := strconv.Atoi(payload.Amount)

	// Simulasikan lookup donasi berdasarkan MerchantOrderId
	fmt.Printf("Received callback for order %s (result=%s)\n", payload.MerchantOrderId, payload.ResultCode)

	if payload.ResultCode == "00" {
		// Success → Update status donasi + tambah saldo ke creator
		fmt.Printf("✅ Donation success: %d added to creator balance\n", amount)
		// update DB: donations.status='success', creators.balance += amount
	} else {
		// Failed → Update status donasi jadi failed
		fmt.Printf("❌ Donation failed: order %s\n", payload.MerchantOrderId)
		// update DB: donations.status='failed'
	}
}
