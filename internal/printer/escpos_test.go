package printer_test

import (
	"strings"
	"testing"
	"time"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/money"
	"autoshop/internal/printer"
)

func sampleInvoice() (*domain.Invoice, *domain.ShopProfile) {
	inv := &domain.Invoice{
		Number:      "INV-2526-0001",
		Date:        time.Date(2026, 6, 29, 14, 30, 0, 0, time.UTC),
		PaymentMode: domain.PaymentCash,
		SubTotal:    money.FromRupees(1000),
		CGST:        money.FromRupees(90),
		SGST:        money.FromRupees(90),
		GrandTotal:  money.FromRupees(1180),
		AmountPaid:  money.FromRupees(1180),
		Items: []domain.InvoiceItem{
			{ProductName: "Brake Pad Set", Quantity: 1, UnitPrice: money.FromRupees(1000), LineTotal: money.FromRupees(1180)},
		},
	}
	shop := &domain.ShopProfile{
		ShopName:      "Sharma Auto Parts",
		Phone:         "9876543210",
		GSTIN:         "27ABCDE1234F1Z5",
		ReceiptFooter: "Thank You Visit Again",
	}
	return inv, shop
}

func TestBuildReceiptContainsInitAndCut(t *testing.T) {
	inv, shop := sampleInvoice()
	data := printer.BuildReceipt(inv, shop)

	// Must start with ESC @ (initialise).
	if len(data) < 2 || data[0] != 0x1B || data[1] != 0x40 {
		t.Fatalf("receipt should start with ESC @, got % x", data[:2])
	}
	// Must end with the partial-cut command (GS V 1).
	tail := data[len(data)-3:]
	if tail[0] != 0x1D || tail[1] != 0x56 || tail[2] != 0x01 {
		t.Fatalf("receipt should end with GS V 1 cut, got % x", tail)
	}
	// Should contain the invoice number and shop name somewhere.
	s := string(data)
	if !strings.Contains(s, "INV-2526-0001") || !strings.Contains(s, "Sharma Auto Parts") {
		t.Fatal("receipt missing invoice number or shop name")
	}
}

func TestPlainTextLayout(t *testing.T) {
	inv, shop := sampleInvoice()
	txt := printer.PlainText(inv, shop)

	for _, want := range []string{
		"Sharma Auto Parts",
		"INV-2526-0001",
		"Brake Pad Set",
		"GRAND TOTAL",
		"Thank You Visit Again",
	} {
		if !strings.Contains(txt, want) {
			t.Errorf("plain receipt missing %q\n---\n%s", want, txt)
		}
	}
}
