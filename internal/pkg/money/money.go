// Package money provides a fixed-point monetary type for the application.
//
// Why a dedicated type? Currency must NEVER be represented as a float64.
// Floating-point math introduces tiny rounding errors (e.g. 0.1 + 0.2 != 0.3)
// that silently corrupt invoice totals, GST amounts and outstanding balances.
// Instead we store every amount as an int64 number of *paise* (1 rupee = 100
// paise). All arithmetic is exact integer arithmetic; rounding only ever
// happens explicitly (e.g. when applying a GST percentage), in one place.
package money

import (
	"database/sql/driver"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Money is an amount in paise (the smallest Indian currency unit).
// 100 paise = ₹1. It is stored in the database as a plain INTEGER column.
type Money int64

// PaisePerRupee is the number of paise in one rupee.
const PaisePerRupee = 100

// Zero is the additive identity, provided for readability.
const Zero = Money(0)

// FromPaise wraps a raw paise value as Money.
func FromPaise(p int64) Money { return Money(p) }

// FromRupees converts a whole-rupee integer to Money.
func FromRupees(r int64) Money { return Money(r * PaisePerRupee) }

// ParseRupees parses a human string such as "1234.50" or "₹1,234.50" into
// Money, rounding to the nearest paise. It is the safe way to ingest a value
// typed by a user without ever going through binary floating point for storage.
func ParseRupees(s string) (Money, error) {
	clean := strings.TrimSpace(s)
	clean = strings.NewReplacer("₹", "", ",", "", " ", "").Replace(clean)
	if clean == "" {
		return 0, nil
	}
	f, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0, fmt.Errorf("money: invalid amount %q: %w", s, err)
	}
	// Round to nearest paise before truncating to int64.
	return Money(math.Round(f * PaisePerRupee)), nil
}

// Paise returns the underlying integer value in paise.
func (m Money) Paise() int64 { return int64(m) }

// Rupees returns the amount as a float64 in rupees. Use ONLY for display or
// JSON output that the frontend formats — never feed it back into storage.
func (m Money) Rupees() float64 { return float64(m) / PaisePerRupee }

// Add returns m + n.
func (m Money) Add(n Money) Money { return m + n }

// Sub returns m - n.
func (m Money) Sub(n Money) Money { return m - n }

// MulQty multiplies the amount by an integer quantity (exact).
func (m Money) MulQty(qty int) Money { return m * Money(qty) }

// Percent returns the given percentage of the amount, rounded to the nearest
// paise. This is the single place GST/discount rounding occurs.
// Example: ₹100.00.Percent(18) == ₹18.00.
func (m Money) Percent(pct float64) Money {
	return Money(math.Round(float64(m) * pct / 100.0))
}

// IsZero reports whether the amount is exactly zero.
func (m Money) IsZero() bool { return m == 0 }

// IsNegative reports whether the amount is below zero.
func (m Money) IsNegative() bool { return m < 0 }

// String renders the amount as a plain "1234.50" string (no symbol/grouping).
// Locale-aware formatting (₹, Indian digit grouping) is done in the frontend.
func (m Money) String() string {
	neg := m < 0
	v := int64(m)
	if neg {
		v = -v
	}
	rupees := v / PaisePerRupee
	paise := v % PaisePerRupee
	sign := ""
	if neg {
		sign = "-"
	}
	return fmt.Sprintf("%s%d.%02d", sign, rupees, paise)
}

// --- GORM / database/sql integration -------------------------------------
//
// Implementing driver.Valuer and sql.Scanner makes Money persist transparently
// as an INTEGER. GormDataType pins the column type across dialects.

// Value implements driver.Valuer (write path).
func (m Money) Value() (driver.Value, error) { return int64(m), nil }

// Scan implements sql.Scanner (read path).
func (m *Money) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		*m = Money(v)
	case nil:
		*m = 0
	case []byte:
		i, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return fmt.Errorf("money: cannot scan %q: %w", v, err)
		}
		*m = Money(i)
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("money: cannot scan %q: %w", v, err)
		}
		*m = Money(i)
	default:
		return fmt.Errorf("money: unsupported scan type %T", src)
	}
	return nil
}

// GormDataType tells GORM to use an INTEGER column for this type.
func (Money) GormDataType() string { return "integer" }
