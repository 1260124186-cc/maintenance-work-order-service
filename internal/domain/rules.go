package domain

import (
	"fmt"
	"strings"
	"time"
)

func ValidateCreateInput(input CreateWorkOrderInput) error {
	if strings.TrimSpace(input.AssetID) == "" || strings.TrimSpace(input.Title) == "" {
		return ErrInvalidWorkOrder
	}
	if input.Priority != "low" && input.Priority != "normal" && input.Priority != "urgent" {
		return fmt.Errorf("%w: %s", ErrUnsupportedPriority, input.Priority)
	}
	return nil
}

func NewWorkOrder(id string, input CreateWorkOrderInput, createdAt time.Time) (WorkOrder, error) {
	if err := ValidateCreateInput(input); err != nil {
		return WorkOrder{}, err
	}
	return WorkOrder{
		ID:        id,
		AssetID:   input.AssetID,
		Title:     strings.TrimSpace(input.Title),
		Priority:  input.Priority,
		Labels:    append([]string(nil), input.Labels...),
		Status:    StatusOpen,
		CreatedAt: createdAt.UTC(),
	}, nil
}

func Assign(order *WorkOrder, technician string) error {
	if strings.TrimSpace(technician) == "" {
		return ErrTechnicianRequired
	}
	if order.Status != StatusOpen {
		return ErrInvalidStatusChange
	}
	order.Technician = strings.TrimSpace(technician)
	order.Status = StatusAssigned
	return nil
}

func Complete(order *WorkOrder, completedAt time.Time) error {
	if order.Status != StatusAssigned {
		return ErrInvalidStatusChange
	}
	order.Status = StatusCompleted
	timestamp := completedAt.UTC()
	order.ClosedAt = &timestamp
	return nil
}

func Summarize(date string, orders []WorkOrder) DailySummary {
	summary := DailySummary{Date: date}
	for _, order := range orders {
		switch order.Status {
		case StatusOpen:
			summary.Open++
		case StatusAssigned:
			summary.Assigned++
		case StatusCompleted:
			summary.Completed++
		}
	}
	return summary
}
