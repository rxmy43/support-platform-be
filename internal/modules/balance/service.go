package balance

import (
	"context"

	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/modules/user"
)

type BalanceService struct {
	balanceRepo *BalanceRepo
	userRepo    *user.UserRepo
}

func NewBalanceService(balanceRepo *BalanceRepo, userRepo *user.UserRepo) *BalanceService {
	return &BalanceService{
		balanceRepo: balanceRepo,
		userRepo:    userRepo,
	}
}

func (s *BalanceService) GetCreatorBalance(ctx context.Context, creatorID uint) (int64, *apperror.AppError) {
	user, err := s.userRepo.FindByID(ctx, creatorID)
	if err != nil {
		return 0, apperror.Unauthorized("invalid user", apperror.CodeUnauthorizedOperation)
	}

	if user.Role != "creator" {
		return 0, apperror.Unauthorized("only creator has balance", apperror.CodeUnauthorizedOperation)
	}

	amount, err := s.balanceRepo.GetBalanceAmountByUserID(ctx, user.ID)
	if err != nil {
		return 0, apperror.InternalServer("failed get balance").WithCause(err)
	}

	return amount, nil
}
