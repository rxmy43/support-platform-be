package auth

type GenerateOTPRequest struct {
	Phone string `json:"string"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}
