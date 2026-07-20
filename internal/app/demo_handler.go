package app

import (
	"context"

	"autoshop/internal/service"
)

// DemoHandler is the Wails binding for loading sample data. Kept separate from
// SettingsHandler so demo tooling stays clearly isolated from real settings.
type DemoHandler struct {
	svc *service.DemoService
}

// NewDemoHandler constructs the handler.
func NewDemoHandler(svc *service.DemoService) *DemoHandler {
	return &DemoHandler{svc: svc}
}

// LoadDemoData populates the database with realistic sample products,
// customers, suppliers and invoices, returning a summary of what was created.
func (h *DemoHandler) LoadDemoData() (service.DemoSummary, error) {
	s, err := h.svc.Load(context.Background())
	return s, bindError(err)
}
