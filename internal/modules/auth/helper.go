package auth

import (
	"fmt"
	"math/rand"
	"sync"
)

var otpStore = sync.Map{}

func generateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func saveOTP(phone, otp string) {
	otpStore.Store(phone, otp)
}
