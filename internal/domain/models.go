package domain

import (
	"errors"
	"time"
)

var (
	ErrAssetNotFound       = errors.New("asset not found")
	ErrAssetUnavailable    = errors.New("asset is not available for maintenance")
	ErrWorkOrderNotFound   = errors.New("work order not found")
	ErrInvalidWorkOrder    = errors.New("invalid work order")
	ErrInvalidStatusChange = errors.New("invalid work order status change")
	ErrTechnicianRequired  = errors.New("technician is required")
	ErrAuditAlreadyOpen    = errors.New("daily audit is already open")
	ErrUnsupportedPriority = errors.New("unsupported priority")
)

type Asset struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Location string   `json:"location"`
	Active   bool     `json:"active"`
	Tags     []string `json:"tags,omitempty"`
}

type Status string

const (
	StatusOpen      Status = "open"
	StatusAssigned  Status = "assigned"
	StatusCompleted Status = "completed"
)

type WorkOrder struct {
	ID         string     `json:"id"`
	AssetID    string     `json:"asset_id"`
	Title      string     `json:"title"`
	Priority   string     `json:"priority"`
	Labels     []string   `json:"labels,omitempty"`
	Status     Status     `json:"status"`
	Technician string     `json:"technician,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	ClosedAt   *time.Time `json:"closed_at,omitempty"`
}

type CreateWorkOrderInput struct {
	AssetID  string   `json:"asset_id"`
	Title    string   `json:"title"`
	Priority string   `json:"priority"`
	Labels   []string `json:"labels"`
}

type DailySummary struct {
	Date      string `json:"date"`
	Open      int    `json:"open"`
	Assigned  int    `json:"assigned"`
	Completed int    `json:"completed"`
}

type AuditCursor interface {
	WorkOrders() []WorkOrder
	Close() error
}
