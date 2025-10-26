package support

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rxmy43/support-platform/internal/repo"
)

type SupportRepo struct {
	*repo.BaseRepo[Support]
}

func NewSupportRepo(DB *sqlx.DB) *SupportRepo {
	return &SupportRepo{
		BaseRepo: &repo.BaseRepo[Support]{
			DB:        DB,
			TableName: "supports",
		},
	}
}

type supportTimestampCreatorFan struct {
	PaymentTimestamp int64  `db:"payment_timestamp"`
	FanID            uint   `db:"fan_id"`
	CreatorID        uint   `db:"creator_id"`
	FanName          string `db:"fan_name"`
	CreatorName      string `db:"creator_name"`
}

func (r *SupportRepo) GetPaymentTimestamp(ctx context.Context, reference, supportID string) (*supportTimestampCreatorFan, error) {
	var result *supportTimestampCreatorFan

	query := `
		SELECT 
			s.payment_timestamp, 
			s.creator_id, 
			s.fan_id, 
			f.name AS fan_name,
			c.name AS creator_name
		FROM supports s
		JOIN users f ON s.fan_id = f.id
		JOIN users f ON s.creator_id = c.id
		WHERE s.reference_code = $1
		AND s.support_id = $2
	`

	err := r.DB.SelectContext(ctx, &result, query, reference, supportID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
