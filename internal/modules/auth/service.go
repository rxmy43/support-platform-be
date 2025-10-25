package auth

import (
	"context"
	"database/sql"

	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/modules/user"
)

type AuthService struct {
	userRepo *user.UserRepo
}

func NewAuthService(userRepo *user.UserRepo) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) GenerateOTP(ctx context.Context, req GenerateOTPRequest) (string, *apperror.AppError) {
	var fieldErrs []apperror.FieldError

	if req.Phone == "" {
		fieldErrs = append(fieldErrs, apperror.NewFieldError("phone", apperror.CodeFieldRequired))
	}

	existing, err := s.userRepo.FindOneByPhone(ctx, req.Phone)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", apperror.InternalServer("failed fetch user by phone number").WithCause(err)
		}

		fieldErrs = append(fieldErrs, apperror.NewFieldError("phone", apperror.CodePhoneInvalid))
	}

	if len(fieldErrs) > 0 {
		return "", apperror.ValidationError("login validation error", fieldErrs)
	}

	otp := generateOTP()
	saveOTP(existing.Phone, otp)
	println("OTP for", existing.Phone, "is", otp)

	return otp, nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, req VerifyOTPRequest) (*UserResponse, *apperror.AppError) {
	val, ok := otpStore.Load(req.Phone)
	if !ok || val.(string) != req.OTP {
		return nil, apperror.BadRequest("invalid otp", apperror.CodeInvalidCredentials)
	}

	user, err := s.userRepo.FindOneByPhone(ctx, req.Phone)
	if err != nil && err == sql.ErrNoRows {
		return nil, apperror.NotFound("user not found", apperror.CodePhoneInvalid).WithNotFoundField("phone")
	}

	otpStore.Delete(req.Phone)
	return &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Phone: user.Phone,
		Role:  user.Role,
	}, nil
}
