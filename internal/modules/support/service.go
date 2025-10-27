package support

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/config"
	"github.com/rxmy43/support-platform/internal/http/middleware"
	"github.com/rxmy43/support-platform/internal/modules/balance"
	"github.com/rxmy43/support-platform/internal/modules/user"
	"github.com/rxmy43/support-platform/internal/socket"
	"github.com/shopspring/decimal"
)

type SupportService struct {
	supportRepo *SupportRepo
	userRepo    *user.UserRepo
	balanceRepo *balance.BalanceRepo
	hub         *socket.Hub
}

func NewSupportService(supportRepo *SupportRepo, userRepo *user.UserRepo, balanceRepo *balance.BalanceRepo, hub *socket.Hub) *SupportService {
	return &SupportService{
		supportRepo: supportRepo,
		userRepo:    userRepo,
		balanceRepo: balanceRepo,
		hub:         hub,
	}
}

func (s *SupportService) generateSignature(merchantCode string, timestamp int64, merchantKey string) string {
	raw := fmt.Sprintf("%s%d%s", merchantCode, timestamp, merchantKey)
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

func (s *SupportService) generateTimestamp() (int64, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return 0, err
	}

	now := time.Now().In(loc)
	timestampMillis := now.UnixNano() / int64(time.Millisecond)

	return timestampMillis, nil
}

func (s *SupportService) generateSupportID(timestamp int64, creatorID, fanID uint) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	randomHex := hex.EncodeToString(b)
	return fmt.Sprintf("SUPPORT/%d/%s/%d%d", timestamp, randomHex, creatorID, fanID)
}

func (s *SupportService) Donate(ctx context.Context, req DonationRequest) (string, *apperror.AppError) {
	// Get and Check fanID
	fanID := middleware.GetUserID(ctx)
	if fanID == nil {
		return "", apperror.Unauthorized("invalid user id", apperror.CodeUnauthorizedOperation)
	}

	fan, err := s.userRepo.FindByID(ctx, *fanID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", apperror.Unauthorized("invalid credentials", apperror.CodeInvalidCredentials)
		}

		return "", apperror.InternalServer("failed checking fan credentials").WithCause(err)
	}

	// Checking fan role
	if fan.Role != "fan" {
		return "", apperror.Forbidden("only fan can donate", apperror.CodeUnknown)
	}

	// Checking amount not zero or negative numbers
	if req.Amount <= 0 {
		return "", apperror.BadRequest("amount cannot be lower than or equal to 0", apperror.CodeNegativeNotAllowed)
	}

	// Check Creator ID
	creator, err := s.userRepo.FindByID(ctx, req.CreatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", apperror.NotFound("creator not found", apperror.CodeResourceNotFound)
		}

		return "", apperror.InternalServer("failed checking creator id").WithCause(err)
	}

	// Checking creator role
	if creator.Role != "creator" {
		return "", apperror.BadRequest("you only allowed to donate to creator", apperror.CodeUnknown)
	}

	// Get Duitku API Config
	cfg := config.Load()
	merchantCode := cfg.Duitku.MerchantCode
	merchantKey := cfg.Duitku.MerchantKey

	// Get Timestamp
	tmstp, e := s.generateTimestamp()
	if e != nil {
		return "", apperror.InternalServer("failed generating timestamp").WithCause(err)
	}

	// Generate Signature
	signature := s.generateSignature(merchantCode, tmstp, merchantKey)

	// Generate Support ID
	supportID := s.generateSupportID(tmstp, creator.ID, fan.ID)

	// Construct Request Payload
	payload := map[string]interface{}{
		"paymentAmount":   req.Amount,
		"merchantOrderId": supportID,
		"productDetails":  fmt.Sprintf("%s provided support of IDR %d to %s", fan.Name, req.Amount, creator.Name),
		"email":           fmt.Sprintf("%s@gmail.com", strings.ReplaceAll(fan.Phone, "+", "")),
		"callbackUrl":     "",
		"returnUrl":       "",
	}

	url := "https://api-sandbox.duitku.com/api/merchant/createInvoice"

	body, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-duitku-signature", signature)

	tmstpStr := strconv.FormatInt(tmstp, 10)
	httpReq.Header.Set("x-duitku-timestamp", tmstpStr)
	httpReq.Header.Set("x-duitku-merchantcode", merchantCode)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", apperror.InternalServer("failed hit Duitku API").WithCause(err)
	}
	defer resp.Body.Close()

	var result struct {
		StatusCode    string `json:"statusCode"`
		StatusMessage string `json:"statusMessage"`
		Reference     string `json:"reference"`
		MerchantCode  string `json:"merchantCode"`
		PaymentURL    string `json:"paymentUrl"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	// Save Support
	newSupport := &Support{
		FanID:            fan.ID,
		CreatorID:        creator.ID,
		Amount:           decimal.NewFromInt(int64(req.Amount)),
		SupportID:        supportID,
		SentAt:           time.Now(),
		ReferenceCode:    result.Reference,
		Status:           "pending",
		PaymentTimestamp: tmstp,
	}

	if err := s.supportRepo.Create(ctx, newSupport); err != nil {
		return "", apperror.InternalServer("failed creating new support record").WithCause(err)
	}

	return result.PaymentURL, nil
}

func (s *SupportService) verifySignature(merchantCode string, timestamp int64, merchantKey string, signature string) bool {
	genSign := s.generateSignature(merchantCode, timestamp, merchantKey)
	return genSign == signature
}

func (s *SupportService) PaymentCallback(ctx context.Context, req PaymentCallbackRequest) (string, *apperror.AppError) {
	const (
		IGNORED string = "IGNORED"
		SUCCESS string = "SUCCESS"
	)

	support, err := s.supportRepo.GetPaymentTimestamp(ctx, req.Reference, req.MerchantOrderID)
	if err != nil {
		return "", apperror.InternalServer("failed get payment timestamp").WithCause(err)
	}

	fmt.Println("SUPPORT CREATOR ID => ", support.CreatorID)
	fmt.Println("SUPPORT CREATOR Name => ", support.CreatorName)
	fmt.Println("SUPPORT Fan ID => ", support.FanID)
	fmt.Println("SUPPORT Fan Name => ", support.FanName)
	fmt.Println("Rereference => ", req.Reference)
	fmt.Println("SupportID => ", req.MerchantOrderID)

	cfg := config.Load()
	merchantKey := cfg.Duitku.MerchantKey

	if !s.verifySignature(req.MerchantCode, support.PaymentTimestamp, merchantKey, req.Signature) {
		return "", apperror.Forbidden("invalid signature", apperror.CodeUnknown)
	}

	if req.ResultCode != "00" {
		return IGNORED, nil
	}

	// Transaction
	tx, err := s.supportRepo.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", apperror.InternalServer("failed begin transaction").WithCause(err)
	}
	defer tx.Rollback()

	// Update Balance Creator
	query := `
		UPDATE balances
		SET amount = amount + $1
		WHERE user_id = $2
	`

	_, err = tx.ExecContext(ctx, query, req.Amount, support.CreatorID)
	if err != nil {
		return "", apperror.InternalServer("failed updating balance").WithCause(err)
	}

	// Update payment status
	query = `
		UPDATE supports
		SET status = 'paid'
		WHERE reference_code = $1
		AND support_id = $2
	`

	_, err = tx.ExecContext(ctx, query, req.Reference, req.MerchantOrderID)
	if err != nil {
		return "", apperror.InternalServer("failed updating payment status").WithCause(err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return "", apperror.InternalServer("failed committing transaction").WithCause(err)
	}

	msg := socket.EventMessage{
		Event: "support_received",
		Data: map[string]interface{}{
			"amount":       req.Amount,
			"reference":    req.Reference,
			"fan_name":     support.FanName,
			"fan_id":       support.FanID,
			"creator_name": support.CreatorName,
			"creator_id":   support.CreatorID,
		},
	}

	s.hub.BroadcastToCreator(support.CreatorID, msg)
	return SUCCESS, nil
}
