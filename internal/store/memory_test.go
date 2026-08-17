package store

import (
	"context"
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
