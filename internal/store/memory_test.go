package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/1260124186-cc/maintenance-work-order-service/internal/domain"
)

func TestMemoryRepositoryCopiesWorkOrderLabels(t *testing.T) {
	repository := NewMemoryRepository()
	order, err := domain.NewWorkOrder("WO-0001", domain.CreateWorkOrderInput{
		AssetID: "FAN-01", Title: "Inspect belt", Priority: "normal", Labels: []string{"safety"},
	}, time.Now())
	if err != nil {
		t.Fatalf("NewWorkOrder() error = %v", err)
	}
	if err := repository.SaveWorkOrder(context.Background(), order); err != nil {
		t.Fatalf("SaveWorkOrder() error = %v", err)
	}
	order.Labels[0] = "changed"
	stored, found, err := repository.GetWorkOrder(context.Background(), order.ID)
	if err != nil || !found {
		t.Fatalf("GetWorkOrder() = (%v, %v, %v)", stored, found, err)
	}
	if stored.Labels[0] != "safety" {
		t.Fatalf("stored labels were changed through caller input: %q", stored.Labels[0])
	}
}

func TestMemoryRepositoryDoesNotSaveWithCanceledContext(t *testing.T) {
	repository := NewMemoryRepository()
	order, err := domain.NewWorkOrder("WO-0001", domain.CreateWorkOrderInput{
		AssetID: "FAN-01", Title: "Inspect belt", Priority: "normal",
	}, time.Now())
	if err != nil {
		t.Fatalf("NewWorkOrder() error = %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := repository.SaveWorkOrder(ctx, order); !errors.Is(err, context.Canceled) {
		t.Fatalf("SaveWorkOrder() error = %v, want canceled context", err)
	}

	orders, err := repository.ListWorkOrders(context.Background())
	if err != nil {
		t.Fatalf("ListWorkOrders() error = %v", err)
	}
	if len(orders) != 0 {
		t.Fatalf("SaveWorkOrder() stored %d work orders, want 0", len(orders))
	}
}

func TestMemoryRepositoryDoesNotSaveWhenContextCanceledWhileWaitingForWrite(t *testing.T) {
	repository := NewMemoryRepository()
	order, err := domain.NewWorkOrder("WO-0001", domain.CreateWorkOrderInput{
		AssetID: "FAN-01", Title: "Inspect belt", Priority: "normal",
	}, time.Now())
	if err != nil {
		t.Fatalf("NewWorkOrder() error = %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	repository.mu.Lock()
	saveErr := make(chan error, 1)
	go func() {
		saveErr <- repository.SaveWorkOrder(ctx, order)
	}()
	cancel()
	repository.mu.Unlock()

	if err := <-saveErr; !errors.Is(err, context.Canceled) {
		t.Fatalf("SaveWorkOrder() error = %v, want canceled context", err)
	}
	orders, err := repository.ListWorkOrders(context.Background())
	if err != nil {
		t.Fatalf("ListWorkOrders() error = %v", err)
	}
	if len(orders) != 0 {
		t.Fatalf("SaveWorkOrder() stored %d work orders, want 0", len(orders))
	}
}

func TestMemoryRepositoryDoesNotUpdateWhenContextCanceledWhileWaitingForWrite(t *testing.T) {
	repository := NewMemoryRepository()
	order, err := domain.NewWorkOrder("WO-0001", domain.CreateWorkOrderInput{
		AssetID: "FAN-01", Title: "Inspect belt", Priority: "normal",
	}, time.Now())
	if err != nil {
		t.Fatalf("NewWorkOrder() error = %v", err)
	}
	if err := repository.SaveWorkOrder(context.Background(), order); err != nil {
		t.Fatalf("SaveWorkOrder() error = %v", err)
	}

	updated := order
	updated.Status = domain.StatusAssigned
	updated.Technician = "Jordan"
	ctx, cancel := context.WithCancel(context.Background())
	checkingCtx := &preLockCheckingContext{
		Context: ctx,
		checked: make(chan struct{}, 1),
		release: make(chan struct{}),
	}
	repository.mu.Lock()
	updateErr := make(chan error, 1)
	go func() {
		updateErr <- repository.UpdateWorkOrder(checkingCtx, updated)
	}()
	select {
	case <-checkingCtx.checked:
	case <-time.After(50 * time.Millisecond):
	}
	cancel()
	close(checkingCtx.release)
	repository.mu.Unlock()

	if err := <-updateErr; !errors.Is(err, context.Canceled) {
		t.Fatalf("UpdateWorkOrder() error = %v, want canceled context", err)
	}
	stored, found, err := repository.GetWorkOrder(context.Background(), order.ID)
	if err != nil || !found {
		t.Fatalf("GetWorkOrder() = (%v, %v, %v)", stored, found, err)
	}
	if stored.Status != domain.StatusOpen || stored.Technician != "" {
		t.Fatalf("UpdateWorkOrder() changed canceled order to %#v", stored)
	}
}

type preLockCheckingContext struct {
	context.Context
	checked chan struct{}
	release chan struct{}
}

func (c *preLockCheckingContext) Err() error {
	err := c.Context.Err()
	select {
	case c.checked <- struct{}{}:
	default:
	}
	<-c.release
	return err
}

func TestAuditCloseReleasesNextReport(t *testing.T) {
	repository := NewMemoryRepository()
	cursor, err := repository.OpenAudit(context.Background())
	if err != nil {
		t.Fatalf("OpenAudit() error = %v", err)
	}
	if err := cursor.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if _, err := repository.OpenAudit(context.Background()); err != nil {
		t.Fatalf("second OpenAudit() error = %v", err)
	}
}
