package service

import (
	"context"

	"github.com/1260124186-cc/maintenance-work-order-service/internal/domain"
)

type Repository interface {
	ListAssets(ctx context.Context) ([]domain.Asset, error)
	GetAsset(ctx context.Context, id string) (*domain.Asset, bool, error)
	SaveWorkOrder(ctx context.Context, order domain.WorkOrder) error
	GetWorkOrder(ctx context.Context, id string) (*domain.WorkOrder, bool, error)
	UpdateWorkOrder(ctx context.Context, order domain.WorkOrder) error
	ListWorkOrders(ctx context.Context) ([]domain.WorkOrder, error)
	OpenAudit(ctx context.Context) (domain.AuditCursor, error)
}
