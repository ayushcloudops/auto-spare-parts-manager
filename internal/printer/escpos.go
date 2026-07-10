// Package printer builds ESC/POS receipts and sends them to a thermal printer.
//
// The byte-building is pure and unit-tested; the OS-specific sending is isolated
// behind the Printer interface so it can be swapped/mocked and so Windows raw
// spooling can be added later without touching receipt formatting.
package printer

import (
	"fmt"
	"strings"

	"autoshop/internal/domain"
)

// width is the character columns for an 80mm printer at Font A (~48 cols).
const width = 48

// ESC/POS control byte sequences.
var (
	escInit      = []byte{0x1B, 0x40}       // initialise printer
	alignLeft    = []byte{0x1B, 0x61, 0x00} // ESC a 0
	alignCenter  = []byte{0x1B, 0x61, 0x01} // ESC a 1
	boldOn       = []byte{0x1B, 0x45, 0x01} // ESC E 1
	boldOff      = []byte{0x1B, 0x45, 0x00} // ESC E 0
	doubleSize   = []byte{0x1D, 0x21, 0x11} // GS ! (double width+height)
	normalSize   = []byte{0x1D, 0x21, 0x00} // GS ! 0
	feedAndCut   = []byte{0x0A, 0x0A, 0x0A, 0x1D, 0x56, 0x01} // feed + partial cut
	lineFeed     = []byte{0x0A}
)

// BuildReceipt renders an invoice into an ESC/POS byte stream for an 80mm
// thermal printer.
func BuildReceipt(inv *domain.Invoice, shop *domain.ShopProfile) []byte {
	body := receiptText(inv, shop)

	out := []byte{}
	out = append(out, escInit...)

	// Header (centered, shop name emphasised).
	out = append(out, alignCenter...)
	out = append(out, boldOn...)
	out = append(out, doubleSize...)
	out = append(out, []byte(shop.ShopName+"\n")...)
	out = append(out, normalSize...)
	out = append(out, boldOff...)
	for _, l := range headerLines(shop) {
		out = append(out, []byte(l+"\n")...)
	}

	// Body (left aligned, monospace layout).
	out = append(out, alignLeft...)
	out = append(out, []byte(body)...)

	// Footer (centered).
	out = append(out, alignCenter...)
	out = append(out, lineFeed...)
	if shop.ReceiptFooter != "" {
		out = append(out, boldOn...)
		out = append(out, []byte(shop.ReceiptFooter+"\n")...)
		out = append(out, boldOff...)
	}
	out = append(out, feedAndCut...)

	return out
}

// PlainText renders the same receipt as plain text (for on-screen preview and a
// universal browser-print fallback).
func PlainText(inv *domain.Invoice, shop *domain.ShopProfile) string {
	var b strings.Builder
	b.WriteString(center(shop.ShopName))
	b.WriteString("\n")
	for _, l := range headerLines(shop) {
		b.WriteString(center(l))
		b.WriteString("\n")
	}
	b.WriteString(receiptText(inv, shop))
	if shop.ReceiptFooter != "" {
		b.WriteString("\n")
		b.WriteString(center(shop.ReceiptFooter))
		b.WriteString("\n")
	}
	return b.String()
}

func headerLines(shop *domain.ShopProfile) []string {
	var lines []string
	addr := strings.TrimSpace(strings.Join([]string{shop.AddressLine1, shop.AddressLine2}, " "))
	if addr != "" {
		lines = append(lines, addr)
	}
	cityLine := strings.TrimSpace(strings.Join([]string{shop.City, shop.State, shop.Pincode}, " "))
	if cityLine != "" {
		lines = append(lines, cityLine)
	}
	if shop.Phone != "" {
		lines = append(lines, "Ph: "+shop.Phone)
	}
	if shop.GSTIN != "" {
		lines = append(lines, "GSTIN: "+shop.GSTIN)
	}
	return lines
}

// receiptText builds the invoice body shared by BuildReceipt and PlainText.
func receiptText(inv *domain.Invoice, shop *domain.ShopProfile) string {
	var b strings.Builder
	b.WriteString(rule())
	b.WriteString(lr("Invoice: "+inv.Number, inv.Date.Format("02-01-2006 15:04")))
	b.WriteString("\n")
	if inv.Customer != nil {
		b.WriteString("Customer: " + inv.Customer.Name + "\n")
	}
	b.WriteString(rule())

	// Column header.
	b.WriteString(fmt.Sprintf("%-24s%4s%9s%11s\n", "Item", "Qty", "Rate", "Amount"))
	b.WriteString(rule())
	for _, it := range inv.Items {
		name := truncate(it.ProductName, 24)
		b.WriteString(fmt.Sprintf("%-24s%4d%9s%11s\n",
			name, it.Quantity, it.UnitPrice.String(), it.LineTotal.String()))
	}
	b.WriteString(rule())

	// Totals.
	b.WriteString(lr("Subtotal", inv.SubTotal.String()))
	b.WriteString("\n")
	if inv.DiscountTotal.Paise() > 0 {
		b.WriteString(lr("Discount", "-"+inv.DiscountTotal.String()))
		b.WriteString("\n")
	}
	if inv.CGST.Paise() > 0 {
		b.WriteString(lr("CGST", inv.CGST.String()))
		b.WriteString("\n")
		b.WriteString(lr("SGST", inv.SGST.String()))
		b.WriteString("\n")
	}
	if inv.IGST.Paise() > 0 {
		b.WriteString(lr("IGST", inv.IGST.String()))
		b.WriteString("\n")
	}
	if inv.RoundOff.Paise() != 0 {
		b.WriteString(lr("Round Off", inv.RoundOff.String()))
		b.WriteString("\n")
	}
	b.WriteString(rule())
	b.WriteString(lr("GRAND TOTAL", "Rs "+inv.GrandTotal.String()))
	b.WriteString("\n")
	b.WriteString(lr("Payment", strings.ToUpper(string(inv.PaymentMode))))
	b.WriteString("\n")
	if inv.AmountDue.Paise() > 0 {
		b.WriteString(lr("Paid", inv.AmountPaid.String()))
		b.WriteString("\n")
		b.WriteString(lr("Due", inv.AmountDue.String()))
		b.WriteString("\n")
	}
	b.WriteString(rule())
	return b.String()
}

// --- layout helpers --------------------------------------------------------

func rule() string { return strings.Repeat("-", width) + "\n" }

// lr left-justifies l and right-justifies r on one line of the receipt width.
func lr(l, r string) string {
	space := width - len(l) - len(r)
	if space < 1 {
		space = 1
	}
	return l + strings.Repeat(" ", space) + r
}

func center(s string) string {
	if len(s) >= width {
		return s
	}
	pad := (width - len(s)) / 2
	return strings.Repeat(" ", pad) + s
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] // hard cut keeps monospace columns aligned
}
