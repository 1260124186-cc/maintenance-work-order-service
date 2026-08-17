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
	orders := cursor.WorkOrders()
	if len(orders) == 0 {
		return domain.Summarize(date, orders), nil
	}
	defer cursor.Close()
	return domain.Summarize(date, orders), nil
}
