package domain

import (
	"testing"
	"time"
)

func TestNewWorkOrderCopiesLabels(t *testing.T) {
	labels := []string{"safety", "roof"}
	order, err := NewWorkOrder("WO-0001", CreateWorkOrderInput{
		AssetID: "FAN-01", Title: "Inspect belt", Priority: "urgent", Labels: labels,
	}, time.Date(2026, 8, 17, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("NewWorkOrder() error = %v", err)
	}
	labels[0] = "changed"
	if order.Labels[0] != "safety" {
		t.Fatalf("order labels were changed through caller input: %q", order.Labels[0])
	}
}

func TestAssignThenComplete(t *testing.T) {
	order, err := NewWorkOrder("WO-0002", CreateWorkOrderInput{
		AssetID: "PUMP-02", Title: "Check pressure", Priority: "normal",
	}, time.Now())
	if err != nil {
		t.Fatalf("NewWorkOrder() error = %v", err)
	}
	if err := Assign(&order, "Kai"); err != nil {
		t.Fatalf("Assign() error = %v", err)
	}
	if err := Complete(&order, time.Now()); err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if order.Status != StatusCompleted || order.ClosedAt == nil {
		t.Fatalf("completed order = %#v", order)
	}
}
