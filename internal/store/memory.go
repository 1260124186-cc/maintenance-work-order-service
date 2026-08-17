package store

import (
	"context"
	"sync"

	"github.com/1260124186-cc/maintenance-work-order-service/internal/domain"
)

type MemoryRepository struct {
	mu             sync.Mutex
	assets         map[string]domain.Asset
	orders         map[string]domain.WorkOrder
	auditOpen      bool
	nextAuditClose error
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		assets: map[string]domain.Asset{
			"FAN-01":  {ID: "FAN-01", Name: "Air Handler", Location: "Plant Room A", Active: true, Tags: []string{"ventilation", "roof"}},
			"PUMP-02": {ID: "PUMP-02", Name: "Drainage Pump", Location: "Basement", Active: true, Tags: []string{"water"}},
			"GEN-03":  {ID: "GEN-03", Name: "Retired Generator", Location: "Annex", Active: false},
		},
		orders: make(map[string]domain.WorkOrder),
	}
}

func (m *MemoryRepository) ListAssets(ctx context.Context) ([]domain.Asset, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	assets := make([]domain.Asset, 0, len(m.assets))
	for _, asset := range m.assets {
		assets = append(assets, copyAsset(asset))
	}
	return assets, nil
}

func (m *MemoryRepository) GetAsset(ctx context.Context, id string) (*domain.Asset, bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	asset, found := m.assets[id]
	if !found {
		return nil, true, nil
	}
	copy := copyAsset(asset)
	return &copy, true, nil
}

func (m *MemoryRepository) SaveWorkOrder(ctx context.Context, order domain.WorkOrder) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.orders[order.ID] = copyOrder(order)
	return nil
}

func (m *MemoryRepository) GetWorkOrder(ctx context.Context, id string) (*domain.WorkOrder, bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	order, found := m.orders[id]
	if !found {
		return nil, false, nil
	}
	copy := copyOrder(order)
	return &copy, true, nil
}

func (m *MemoryRepository) UpdateWorkOrder(ctx context.Context, order domain.WorkOrder) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, found := m.orders[order.ID]; !found {
		return domain.ErrWorkOrderNotFound
	}
	m.orders[order.ID] = copyOrder(order)
	return nil
}

func (m *MemoryRepository) ListWorkOrders(ctx context.Context) ([]domain.WorkOrder, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.workOrdersLocked(), nil
}

func (m *MemoryRepository) OpenAudit(ctx context.Context) (domain.AuditCursor, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.auditOpen {
		return nil, domain.ErrAuditAlreadyOpen
	}
	m.auditOpen = true
	return &memoryAuditCursor{
		repository: m,
		orders:     m.workOrdersLocked(),
	}, nil
}

func (m *MemoryRepository) SetNextAuditCloseError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextAuditClose = err
}

func (m *MemoryRepository) workOrdersLocked() []domain.WorkOrder {
	orders := make([]domain.WorkOrder, 0, len(m.orders))
	for _, order := range m.orders {
		orders = append(orders, copyOrder(order))
	}
	return orders
}

func copyAsset(asset domain.Asset) domain.Asset {
	asset.Tags = append([]string(nil), asset.Tags...)
	return asset
}

func copyOrder(order domain.WorkOrder) domain.WorkOrder {
	order.Labels = append([]string(nil), order.Labels...)
	if order.ClosedAt != nil {
		closedAt := *order.ClosedAt
		order.ClosedAt = &closedAt
	}
	return order
}

type memoryAuditCursor struct {
	repository *MemoryRepository
	orders     []domain.WorkOrder
	closed     bool
}

func (c *memoryAuditCursor) WorkOrders() []domain.WorkOrder {
	return append([]domain.WorkOrder(nil), c.orders...)
}

func (c *memoryAuditCursor) Close() error {
	c.repository.mu.Lock()
	defer c.repository.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	c.repository.auditOpen = false
	err := c.repository.nextAuditClose
	c.repository.nextAuditClose = nil
	return err
}
