package support

type DonationRequest struct {
	Amount    int  `json:"amount"`
	CreatorID uint `json:"creator_id"`
}

type PaymentCallbackRequest struct {
	MerchantCode     string `json:"merchantCode"`
	Amount           int    `json:"amount"`
	MerchantOrderID  string `json:"merchantOrderId"`
	ProductDetail    string `json:"productDetail"`
	PaymentCode      string `json:"paymentCode"`
	ResultCode       string `json:"resultCode"`
	MerchantUserID   string `json:"merchantUserId"`
	Reference        string `json:"reference"`
	Signature        string `json:"signature"`
	PublisherOrderId string `json:"publisherOrderId"`
	SPUserHash       string `json:"spUserHash"`
	SettlementDate   string `json:"settlementDate"`
	IssuerCode       string `json:"issuerCode"`
}
