package service

import "autoshop/internal/pkg/money"

// GST engine — pure functions, no database. This is the tax heart of billing and
// is exhaustively unit-tested. India GST rules implemented:
//
//   - Tax is charged on the taxable value = (unit price × qty) − line discount.
//   - Intra-state sale: tax splits equally into CGST + SGST (rate/2 each).
//   - Inter-state sale: tax charged wholly as IGST (full rate).
//   - Each line's tax is rounded to the nearest paise (in money.Percent).
//   - The final payable is rounded to the nearest rupee; the difference is
//     recorded as RoundOff (standard on Indian invoices).

// GSTLineInput is one billing line fed to the engine.
type GSTLineInput struct {
	UnitPrice money.Money
	Quantity  int
	Discount  money.Money // absolute, on the whole line
	GSTRate   float64     // percent, e.g. 18
}

// GSTLineResult is the computed tax breakdown for one line.
type GSTLineResult struct {
	TaxableValue money.Money
	CGST         money.Money
	SGST         money.Money
	IGST         money.Money
	LineTotal    money.Money
}

// GSTResult is the fully computed invoice tax summary.
type GSTResult struct {
	Lines         []GSTLineResult
	SubTotal      money.Money // sum of taxable values
	DiscountTotal money.Money // sum of line discounts
	CGST          money.Money
	SGST          money.Money
	IGST          money.Money
	TaxTotal      money.Money
	RoundOff      money.Money // adjustment to reach a whole-rupee grand total
	GrandTotal    money.Money // final payable (rounded to nearest rupee)
}

// CalculateGST computes the tax breakdown for a set of lines. interState selects
// IGST (true) vs CGST/SGST (false).
func CalculateGST(lines []GSTLineInput, interState bool) GSTResult {
	res := GSTResult{Lines: make([]GSTLineResult, 0, len(lines))}

	for _, in := range lines {
		gross := in.UnitPrice.MulQty(in.Quantity)
		taxable := gross.Sub(in.Discount)
		if taxable.IsNegative() {
			taxable = money.Zero
		}

		var lr GSTLineResult
		lr.TaxableValue = taxable
		if interState {
			lr.IGST = taxable.Percent(in.GSTRate)
		} else {
			half := in.GSTRate / 2
			lr.CGST = taxable.Percent(half)
			lr.SGST = taxable.Percent(half)
		}
		lr.LineTotal = taxable.Add(lr.CGST).Add(lr.SGST).Add(lr.IGST)

		res.Lines = append(res.Lines, lr)
		res.SubTotal = res.SubTotal.Add(taxable)
		res.DiscountTotal = res.DiscountTotal.Add(in.Discount)
		res.CGST = res.CGST.Add(lr.CGST)
		res.SGST = res.SGST.Add(lr.SGST)
		res.IGST = res.IGST.Add(lr.IGST)
	}

	res.TaxTotal = res.CGST.Add(res.SGST).Add(res.IGST)
	preRound := res.SubTotal.Add(res.TaxTotal)
	rounded := roundToRupee(preRound)
	res.RoundOff = rounded.Sub(preRound)
	res.GrandTotal = rounded
	return res
}

// roundToRupee rounds an amount to the nearest whole rupee (100 paise).
func roundToRupee(m money.Money) money.Money {
	p := m.Paise()
	neg := p < 0
	if neg {
		p = -p
	}
	r := ((p + 50) / 100) * 100
	if neg {
		r = -r
	}
	return money.FromPaise(r)
}
