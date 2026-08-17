package service

import (
	"context"

	"github.com/1260124186-cc/maintenance-work-order-service/internal/domain"
)

func (s *WorkOrderService) DailySummary(ctx context.Context, date string) (summary domain.DailySummary, err error) {
	cursor, err := s.repository.OpenAudit(ctx)
	if err != nil {
		return domain.DailySummary{}, err
	}
	defer func() {
		if closeErr := cursor.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	return domain.Summarize(date, cursor.WorkOrders()), nil
}
