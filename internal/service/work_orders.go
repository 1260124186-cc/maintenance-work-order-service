package service

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/1260124186-cc/maintenance-work-order-service/internal/domain"
)

type WorkOrderService struct {
	repository Repository
	now        func() time.Time
	nextID     atomic.Uint64
}

func NewWorkOrderService(repository Repository) *WorkOrderService {
	return &WorkOrderService{
		repository: repository,
		now:        time.Now,
	}
}

func (s *WorkOrderService) ListAssets(ctx context.Context) ([]domain.Asset, error) {
	return s.repository.ListAssets(ctx)
}

func (s *WorkOrderService) Create(ctx context.Context, input domain.CreateWorkOrderInput) (domain.WorkOrder, error) {
	if err := domain.ValidateCreateInput(input); err != nil {
		return domain.WorkOrder{}, err
	}
	asset, found, err := s.repository.GetAsset(ctx, input.AssetID)
	if err != nil {
		return domain.WorkOrder{}, err
	}
	if !found || asset == nil {
		return domain.WorkOrder{}, domain.ErrAssetNotFound
	}
	if !asset.Active {
		return domain.WorkOrder{}, domain.ErrAssetUnavailable
	}

	id := fmt.Sprintf("WO-%04d", s.nextID.Add(1))
	order, err := domain.NewWorkOrder(id, input, s.now())
	if err != nil {
		return domain.WorkOrder{}, err
	}
	if err := s.repository.SaveWorkOrder(ctx, order); err != nil {
		return domain.WorkOrder{}, err
	}
	return order, nil
}

func (s *WorkOrderService) Assign(ctx context.Context, id, technician string) (domain.WorkOrder, error) {
	order, found, err := s.repository.GetWorkOrder(ctx, id)
	if err != nil {
		return domain.WorkOrder{}, err
	}
	if !found || order == nil {
		return domain.WorkOrder{}, domain.ErrWorkOrderNotFound
	}
	if err := domain.Assign(order, technician); err != nil {
		return domain.WorkOrder{}, err
	}
	if err := s.repository.UpdateWorkOrder(ctx, *order); err != nil {
		return domain.WorkOrder{}, err
	}
	return *order, nil
}

func (s *WorkOrderService) Complete(ctx context.Context, id string) (domain.WorkOrder, error) {
	order, found, err := s.repository.GetWorkOrder(ctx, id)
	if err != nil {
		return domain.WorkOrder{}, err
	}
	if !found || order == nil {
		return domain.WorkOrder{}, domain.ErrWorkOrderNotFound
	}
	if err := domain.Complete(order, s.now()); err != nil {
		return domain.WorkOrder{}, err
	}
	if err := s.repository.UpdateWorkOrder(ctx, *order); err != nil {
		return domain.WorkOrder{}, err
	}
	return *order, nil
}
