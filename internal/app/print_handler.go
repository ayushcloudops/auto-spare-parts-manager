package app

import (
	"context"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"
	"autoshop/internal/printer"
)

// PrintHandler is the Wails binding for receipt preview and printing.
type PrintHandler struct {
	invoices domain.InvoiceRepository
	settings domain.SettingsRepository
	printer  printer.Printer
}

// NewPrintHandler constructs the handler.
func NewPrintHandler(invoices domain.InvoiceRepository, settings domain.SettingsRepository, p printer.Printer) *PrintHandler {
	return &PrintHandler{invoices: invoices, settings: settings, printer: p}
}

func (h *PrintHandler) ctx() context.Context { return context.Background() }

// load fetches the invoice and shop profile needed to render a receipt.
func (h *PrintHandler) load(invoiceID uint) (*domain.Invoice, *domain.ShopProfile, error) {
	inv, err := h.invoices.FindByID(h.ctx(), invoiceID)
	if err != nil {
		return nil, nil, err
	}
	profile, err := h.settings.GetShopProfile(h.ctx())
	if err != nil {
		return nil, nil, err
	}
	return inv, profile, nil
}

// PreviewReceipt returns a plain-text rendering of the receipt for on-screen
// display (and a universal browser-print fallback).
func (h *PrintHandler) PreviewReceipt(invoiceID uint) (string, error) {
	inv, profile, err := h.load(invoiceID)
	if err != nil {
		return "", bindError(err)
	}
	return printer.PlainText(inv, profile), nil
}

// PrintReceipt builds an ESC/POS receipt and sends it to the configured thermal
// printer.
func (h *PrintHandler) PrintReceipt(invoiceID uint) error {
	inv, profile, err := h.load(invoiceID)
	if err != nil {
		return bindError(err)
	}
	name, _ := h.settings.Get(h.ctx(), domain.SettingPrinterName)
	data := printer.BuildReceipt(inv, profile)
	if err := h.printer.Print(data, name); err != nil {
		return bindError(apperr.Internal(err, "could not print receipt"))
	}
	return nil
}
