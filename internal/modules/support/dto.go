package support

import "time"

type DonationRequest struct {
	Amount    int  `json:"amount"`
	CreatorID uint `json:"creator_id"`
}

type PaymentCallbackRequest struct {
	MerchantCode     string `schema:"merchantCode"`
	Amount           string `schema:"amount"`
	MerchantOrderID  string `schema:"merchantOrderId"`
	ProductDetail    string `schema:"productDetail"`
	AdditionalParam  string `schema:"additionalParam"`
	ResultCode       string `schema:"resultCode"`
	PaymentCode      string `schema:"paymentCode"`
	MerchantUserID   string `schema:"merchantUserId"`
	Reference        string `schema:"reference"`
	Signature        string `schema:"signature"`
	PublisherOrderID string `schema:"publisherOrderId"`
	SettlementDate   string `schema:"settlementDate"`
	VaNumber         string `schema:"vaNumber"`
	SourceAccount    string `schema:"sourceAccount"`
}

type BestSupporters struct {
	ID      uint      `json:"id" db:"id"`
	FanName string    `json:"fan_name" db:"fan_name"`
	Amount  int64     `json:"amount" db:"amount"`
	SentAt  time.Time `json:"sent_at" db:"sent_at"`
}

type FanSupportHistory struct {
	ID          uint      `json:"id" db:"id"`
	CreatorName string    `json:"creator_name" db:"creator_name"`
	Amount      int64     `json:"amount" db:"amount"`
	SentAt      time.Time `json:"sent_at" db:"sent_at"`
}
