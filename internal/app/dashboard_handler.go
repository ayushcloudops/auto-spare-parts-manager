package app

import (
	"context"

	"autoshop/internal/domain"
	"autoshop/internal/service"
)

// DashboardHandler is the Wails binding for the home dashboard.
type DashboardHandler struct {
	svc *service.DashboardService
}

// NewDashboardHandler constructs the handler.
func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// Stats returns the aggregated dashboard figures.
func (h *DashboardHandler) Stats() (domain.DashboardStats, error) {
	s, err := h.svc.Stats(context.Background())
	return s, bindError(err)
}
