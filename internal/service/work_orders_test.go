package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/1260124186-cc/maintenance-work-order-service/internal/domain"
	"github.com/1260124186-cc/maintenance-work-order-service/internal/service"
	"github.com/1260124186-cc/maintenance-work-order-service/internal/store"
)

func TestCreateAssignAndCompleteWorkOrder(t *testing.T) {
	operations := service.NewWorkOrderService(store.NewMemoryRepository())
	created, err := operations.Create(context.Background(), domain.CreateWorkOrderInput{
		AssetID: "FAN-01", Title: "Inspect belt tension", Priority: "urgent", Labels: []string{"safety"},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	assigned, err := operations.Assign(context.Background(), created.ID, "Jordan")
	if err != nil {
		t.Fatalf("Assign() error = %v", err)
	}
	completed, err := operations.Complete(context.Background(), assigned.ID)
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if completed.Status != domain.StatusCompleted || completed.Technician != "Jordan" {
		t.Fatalf("completed work order = %#v", completed)
	}
}

func TestCreateRejectsInactiveAsset(t *testing.T) {
	operations := service.NewWorkOrderService(store.NewMemoryRepository())
	_, err := operations.Create(context.Background(), domain.CreateWorkOrderInput{
		AssetID: "GEN-03", Title: "Inspect fuel", Priority: "normal",
	})
	if !errors.Is(err, domain.ErrAssetUnavailable) {
		t.Fatalf("Create() error = %v, want inactive asset error", err)
	}
}

func TestDailySummaryReturnsCloseError(t *testing.T) {
	repository := store.NewMemoryRepository()
	repository.SetNextAuditCloseError(errors.New("audit close failed"))
	operations := service.NewWorkOrderService(repository)
	_, err := operations.DailySummary(context.Background(), "2026-08-17")
	if err == nil || err.Error() != "audit close failed" {
		t.Fatalf("DailySummary() error = %v, want audit close failed", err)
	}
}

func TestCreateHonorsCanceledContext(t *testing.T) {
	repository := store.NewMemoryRepository()
	operations := service.NewWorkOrderService(repository)
	context, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := operations.Create(context, domain.CreateWorkOrderInput{
		AssetID: "FAN-01", Title: "Inspect belt", Priority: "normal",
	})
	if !errors.Is(err, context.Err()) {
		t.Fatalf("Create() error = %v, want canceled context", err)
	}
}

func TestDailySummaryCompletesWithoutWorkOrders(t *testing.T) {
	operations := service.NewWorkOrderService(store.NewMemoryRepository())
	summary, err := operations.DailySummary(context.Background(), "2026-08-17")
	if err != nil {
		t.Fatalf("DailySummary() error = %v", err)
	}
	if summary.Completed != 0 || summary.Open != 0 || summary.Assigned != 0 {
		t.Fatalf("empty DailySummary() = %#v", summary)
	}
	if _, err := operations.DailySummary(context.Background(), "2026-08-18"); err != nil {
		t.Fatalf("second DailySummary() error = %v", err)
	}
}

func TestDailySummaryCountsWorkOrderStatuses(t *testing.T) {
	operations := service.NewWorkOrderService(store.NewMemoryRepository())
	open, err := operations.Create(context.Background(), domain.CreateWorkOrderInput{
		AssetID: "FAN-01", Title: "Inspect belt", Priority: "normal",
	})
	if err != nil {
		t.Fatalf("Create() open order error = %v", err)
	}
	assigned, err := operations.Create(context.Background(), domain.CreateWorkOrderInput{
		AssetID: "PUMP-02", Title: "Inspect pump", Priority: "urgent",
	})
	if err != nil {
		t.Fatalf("Create() assigned order error = %v", err)
	}
	if _, err := operations.Assign(context.Background(), assigned.ID, "Jordan"); err != nil {
		t.Fatalf("Assign() error = %v", err)
	}
	completed, err := operations.Create(context.Background(), domain.CreateWorkOrderInput{
		AssetID: "FAN-01", Title: "Replace filter", Priority: "low",
	})
	if err != nil {
		t.Fatalf("Create() completed order error = %v", err)
	}
	if _, err := operations.Assign(context.Background(), completed.ID, "Casey"); err != nil {
		t.Fatalf("Assign() completed order error = %v", err)
	}
	if _, err := operations.Complete(context.Background(), completed.ID); err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	summary, err := operations.DailySummary(context.Background(), "2026-08-17")
	if err != nil {
		t.Fatalf("DailySummary() error = %v", err)
	}
	if summary.Open != 1 || summary.Assigned != 1 || summary.Completed != 1 {
		t.Fatalf("DailySummary() = %#v, want one order in each status", summary)
	}
	if open.ID == "" {
		t.Fatal("Create() returned empty open order ID")
	}
}

func TestCreateRejectsInvalidPriorityBeforeLookingUpAsset(t *testing.T) {
	repository := &lookupTrackingRepository{}
	operations := service.NewWorkOrderService(repository)
	_, err := operations.Create(context.Background(), domain.CreateWorkOrderInput{
		AssetID: "FAN-01", Title: "Inspect belt", Priority: "immediate",
	})
	if !errors.Is(err, domain.ErrUnsupportedPriority) {
		t.Fatalf("Create() error = %v, want priority validation error", err)
	}
	if repository.lookups != 0 {
		t.Fatalf("GetAsset() calls = %d, want 0", repository.lookups)
	}
}

type lookupTrackingRepository struct {
	lookups int
}

func (r *lookupTrackingRepository) ListAssets(context.Context) ([]domain.Asset, error) {
	return nil, nil
}
func (r *lookupTrackingRepository) GetAsset(context.Context, string) (*domain.Asset, bool, error) {
	r.lookups++
	return nil, false, nil
}
func (r *lookupTrackingRepository) SaveWorkOrder(context.Context, domain.WorkOrder) error { return nil }
func (r *lookupTrackingRepository) GetWorkOrder(context.Context, string) (*domain.WorkOrder, bool, error) {
	return nil, false, nil
}
func (r *lookupTrackingRepository) UpdateWorkOrder(context.Context, domain.WorkOrder) error {
	return nil
}
func (r *lookupTrackingRepository) ListWorkOrders(context.Context) ([]domain.WorkOrder, error) {
	return nil, nil
}
func (r *lookupTrackingRepository) OpenAudit(context.Context) (domain.AuditCursor, error) {
	return nil, errors.New("not implemented")
}
