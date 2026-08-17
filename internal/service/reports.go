package service

import (
	"context"

	"github.com/1260124186-cc/maintenance-work-order-service/internal/domain"
)

func (s *WorkOrderService) DailySummary(ctx context.Context, date string) (domain.DailySummary, error) {
	cursor, err := s.repository.OpenAudit(ctx)
	if err != nil {
		return domain.DailySummary{}, err
	}
	summary := domain.Summarize(date, cursor.WorkOrders())
	if closeErr := cursor.Close(); closeErr != nil {
		return domain.DailySummary{}, closeErr
	}
	return summary, nil
}
