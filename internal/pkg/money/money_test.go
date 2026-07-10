package money_test

import (
	"testing"

	"autoshop/internal/pkg/money"
)

func TestParseRupees(t *testing.T) {
	cases := []struct {
		in   string
		want money.Money
	}{
		{"0", 0},
		{"1234.50", 123450},
		{"₹1,234.50", 123450},
		{"  99.99 ", 9999},
		{"100", 10000},
		{"0.1", 10},     // 10 paise
		{"19.995", 2000}, // rounds to nearest paise (₹20.00)
		{"", 0},
	}
	for _, c := range cases {
		got, err := money.ParseRupees(c.in)
		if err != nil {
			t.Fatalf("ParseRupees(%q): unexpected error %v", c.in, err)
		}
		if got != c.want {
			t.Errorf("ParseRupees(%q) = %d paise, want %d", c.in, got.Paise(), c.want.Paise())
		}
	}
}

func TestParseRupeesInvalid(t *testing.T) {
	if _, err := money.ParseRupees("abc"); err == nil {
		t.Fatal("expected error for non-numeric input")
	}
}

func TestPercentRounding(t *testing.T) {
	// 18% GST on ₹100.00 = ₹18.00
	if got := money.FromRupees(100).Percent(18); got != money.FromPaise(1800) {
		t.Errorf("18%% of ₹100 = %s, want 18.00", got)
	}
	// 2.5% (CGST half of 5%) on ₹99.00 = ₹2.475 -> rounds to ₹2.48 (248 paise)
	if got := money.FromPaise(9900).Percent(2.5); got != money.FromPaise(248) {
		t.Errorf("2.5%% of ₹99 = %s (%d paise), want 2.48", got, got.Paise())
	}
}

func TestArithmeticAndString(t *testing.T) {
	price := money.FromPaise(19999) // ₹199.99
	line := price.MulQty(3)         // ₹599.97
	if line != money.FromPaise(59997) {
		t.Errorf("MulQty: got %s, want 599.97", line)
	}
	if line.String() != "599.97" {
		t.Errorf("String: got %q, want 599.97", line.String())
	}
	if money.FromPaise(-5).String() != "-0.05" {
		t.Errorf("negative String: got %q, want -0.05", money.FromPaise(-5).String())
	}
}
