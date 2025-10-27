package support

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
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

	appUrl := cfg.AppURL

	// Construct Request Payload
	payload := map[string]interface{}{
		"paymentAmount":   req.Amount,
		"merchantOrderId": supportID,
		"productDetails":  fmt.Sprintf("%s provided support of IDR %d to %s", fan.Name, req.Amount, creator.Name),
		"email":           fmt.Sprintf("%s@gmail.com", strings.ReplaceAll(fan.Phone, "+", "")),
		"callbackUrl":     fmt.Sprintf("%s/api/payment/callback", appUrl),
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

func (s *SupportService) verifySignature(merchantCode, amount, merchantOrderID, merchantKey, signature string) bool {
	raw := merchantCode + amount + merchantOrderID + merchantKey
	hash := md5.Sum([]byte(raw))
	genSig := hex.EncodeToString(hash[:])
	return genSig == signature
}

func (s *SupportService) PaymentCallback(ctx context.Context, req PaymentCallbackRequest) (string, *apperror.AppError) {
	const (
		IGNORED string = "IGNORED"
		SUCCESS string = "SUCCESS"
	)

	log.Println("=== PaymentCallback START ===")
	log.Println("Request received:", req)

	support, err := s.supportRepo.GetPaymentTimestamp(ctx, req.Reference, req.MerchantOrderID)
	if err != nil {
		log.Println("Failed to get payment timestamp:", err)
		return "", apperror.InternalServer("failed get payment timestamp").WithCause(err)
	}

	cfg := config.Load()
	merchantKey := cfg.Duitku.MerchantKey

	if !s.verifySignature(req.MerchantCode, req.Amount, req.MerchantOrderID, merchantKey, req.Signature) {
		return "", apperror.Forbidden("invalid signature", apperror.CodeUnknown)
	}
	log.Println("Signature verification PASSED")

	if req.ResultCode != "00" {
		log.Println("ResultCode not 00, ignoring transaction")
		return IGNORED, nil
	}
	log.Println("ResultCode is 00, processing transaction")

	tx, err := s.supportRepo.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return "", apperror.InternalServer("failed begin transaction").WithCause(err)
	}
	defer tx.Rollback()
	log.Println("Transaction started")

	amountInt, err := strconv.ParseInt(req.Amount, 10, 64)
	if err != nil {
		log.Println("Invalid amount:", req.Amount)
		return "", apperror.BadRequest("invalid amount", apperror.CodeUnknown)
	}

	query := `
    UPDATE balances
    SET amount = amount + $1
    WHERE user_id = $2
`
	res, err := tx.ExecContext(ctx, query, amountInt, support.CreatorID)
	if err != nil {
		log.Println("Failed updating balance:", err)
		return "", apperror.InternalServer("failed updating balance").WithCause(err)
	}

	affected, _ := res.RowsAffected()
	log.Println("Balance update affected rows:", affected)

	query = `
		UPDATE supports
		SET status = 'paid'
		WHERE reference_code = $1
		AND support_id = $2
	`
	res, err = tx.ExecContext(ctx, query, req.Reference, req.MerchantOrderID)
	if err != nil {
		log.Println("Failed updating support status:", err)
		return "", apperror.InternalServer("failed updating payment status").WithCause(err)
	}
	affected, _ = res.RowsAffected()
	log.Println("Support status update affected rows:", affected)

	if err := tx.Commit(); err != nil {
		log.Println("Failed committing transaction:", err)
		return "", apperror.InternalServer("failed committing transaction").WithCause(err)
	}
	log.Println("Transaction committed")

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

	log.Println("Broadcasting event to creator:", support.CreatorID)
	s.hub.BroadcastToCreator(support.CreatorID, msg)
	log.Println("=== PaymentCallback END ===")
	return SUCCESS, nil
}

func (s *SupportService) GetSupporters(ctx context.Context, cursor *uint, creatorID uint) ([]BestSupporters, *uint, *apperror.AppError) {
	user, err := s.userRepo.FindByID(ctx, creatorID)
	if err != nil {
		return nil, nil, apperror.Unauthorized("Unauthorized", apperror.CodeUnauthorizedOperation)
	}

	if user.Role != "creator" {
		return nil, nil, apperror.Unauthorized("Unauthorized", apperror.CodeUnauthorizedOperation)
	}

	bestSupports, nextCursor, err := s.supportRepo.GetCreatorSupporters(ctx, cursor, user.ID)
	if err != nil {
		return nil, nil, apperror.InternalServer("failed to get creator's best supporters").WithCause(err)
	}

	if len(bestSupports) == 0 {
		return []BestSupporters{}, nil, nil
	}

	return bestSupports, nextCursor, nil
}

func (s *SupportService) GetFanSpending(ctx context.Context, fanID uint) (int64, *apperror.AppError) {
	user, err := s.userRepo.FindByID(ctx, fanID)
	if err != nil {
		return 0, apperror.Unauthorized("invalid user", apperror.CodeUnauthorizedOperation)
	}

	if user.Role != "fan" {
		return 0, apperror.Unauthorized("invalid role", apperror.CodeUnauthorizedOperation)
	}

	amount, err := s.supportRepo.GetFanSpendingAmount(ctx, user.ID)
	if err != nil {
		return 0, apperror.InternalServer("failed get fan spending").WithCause(err)
	}

	return amount, nil
}

func (s *SupportService) GetFanSupportHistory(ctx context.Context, cursor, fanID *uint) ([]FanSupportHistory, *uint, *apperror.AppError) {
	user, err := s.userRepo.FindByID(ctx, *fanID)
	if err != nil {
		return nil, nil, apperror.Unauthorized("invalid user", apperror.CodeUnauthorizedOperation)
	}

	if user.Role != "fan" {
		return nil, nil, apperror.Unauthorized("invalid role", apperror.CodeUnauthorizedOperation)
	}

	histories, nextCursor, err := s.supportRepo.GetFanSupportHistory(ctx, cursor, fanID)
	if err != nil {
		return nil, nil, apperror.InternalServer("failed get fan spending history").WithCause(err)
	}

	return histories, nextCursor, nil
}
