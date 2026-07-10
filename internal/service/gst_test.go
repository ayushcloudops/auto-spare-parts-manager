package service_test

import (
	"testing"

	"autoshop/internal/pkg/money"
	"autoshop/internal/service"
)

func TestGSTIntraState(t *testing.T) {
	// One line: ₹1000 × 1, 18% GST, no discount. Intra-state.
	res := service.CalculateGST([]service.GSTLineInput{
		{UnitPrice: money.FromRupees(1000), Quantity: 1, GSTRate: 18},
	}, false)

	if res.SubTotal != money.FromRupees(1000) {
		t.Fatalf("subtotal: got %s", res.SubTotal)
	}
	// 18% = ₹180 → CGST ₹90 + SGST ₹90.
	if res.CGST != money.FromRupees(90) || res.SGST != money.FromRupees(90) {
		t.Fatalf("cgst/sgst: got %s / %s", res.CGST, res.SGST)
	}
	if res.IGST != money.Zero {
		t.Fatalf("igst should be zero intra-state, got %s", res.IGST)
	}
	if res.GrandTotal != money.FromRupees(1180) {
		t.Fatalf("grand total: got %s want 1180.00", res.GrandTotal)
	}
}

func TestGSTInterState(t *testing.T) {
	res := service.CalculateGST([]service.GSTLineInput{
		{UnitPrice: money.FromRupees(1000), Quantity: 1, GSTRate: 18},
	}, true)

	if res.IGST != money.FromRupees(180) {
		t.Fatalf("igst: got %s want 180.00", res.IGST)
	}
	if res.CGST != money.Zero || res.SGST != money.Zero {
		t.Fatalf("cgst/sgst should be zero inter-state")
	}
	if res.GrandTotal != money.FromRupees(1180) {
		t.Fatalf("grand total: got %s", res.GrandTotal)
	}
}

func TestGSTWithDiscountAndQty(t *testing.T) {
	// 3 × ₹250 = ₹750, less ₹50 discount = ₹700 taxable, 28% = ₹196.
	res := service.CalculateGST([]service.GSTLineInput{
		{UnitPrice: money.FromRupees(250), Quantity: 3, Discount: money.FromRupees(50), GSTRate: 28},
	}, false)

	if res.SubTotal != money.FromRupees(700) {
		t.Fatalf("taxable: got %s want 700.00", res.SubTotal)
	}
	// 14% each of 700 = ₹98.
	if res.CGST != money.FromRupees(98) || res.SGST != money.FromRupees(98) {
		t.Fatalf("cgst/sgst: got %s / %s", res.CGST, res.SGST)
	}
	if res.GrandTotal != money.FromRupees(896) {
		t.Fatalf("grand total: got %s want 896.00", res.GrandTotal)
	}
}

func TestGSTRoundOff(t *testing.T) {
	// ₹99.50 × 1, 0% GST → preRound ₹99.50 → rounds to ₹100, roundOff +₹0.50.
	res := service.CalculateGST([]service.GSTLineInput{
		{UnitPrice: money.FromPaise(9950), Quantity: 1, GSTRate: 0},
	}, false)
	if res.GrandTotal != money.FromRupees(100) {
		t.Fatalf("grand total: got %s want 100.00", res.GrandTotal)
	}
	if res.RoundOff != money.FromPaise(50) {
		t.Fatalf("round off: got %s want 0.50", res.RoundOff)
	}
}

func TestGSTMultiLineTotals(t *testing.T) {
	res := service.CalculateGST([]service.GSTLineInput{
		{UnitPrice: money.FromRupees(500), Quantity: 2, GSTRate: 18}, // 1000 taxable
		{UnitPrice: money.FromRupees(300), Quantity: 1, GSTRate: 28}, // 300 taxable
	}, false)

	if res.SubTotal != money.FromRupees(1300) {
		t.Fatalf("subtotal: got %s", res.SubTotal)
	}
	// Tax: 180 + 84 = 264. Grand = 1564.
	if res.TaxTotal != money.FromRupees(264) {
		t.Fatalf("tax total: got %s want 264.00", res.TaxTotal)
	}
	if res.GrandTotal != money.FromRupees(1564) {
		t.Fatalf("grand total: got %s want 1564.00", res.GrandTotal)
	}
}
