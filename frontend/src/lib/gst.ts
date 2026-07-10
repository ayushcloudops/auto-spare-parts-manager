// Client-side GST preview, mirroring the Go engine (internal/service/gst.go) so
// the cashier sees live totals as they build a bill. The AUTHORITATIVE totals
// are always recomputed by the backend on save — this is display only.
//
// All amounts are in paise (integers), matching money.Money.

export interface PreviewLine {
  unitPrice: number;
  quantity: number;
  discount: number;
  gstRate: number;
}

export interface PreviewTotals {
  subTotal: number;
  discountTotal: number;
  cgst: number;
  sgst: number;
  igst: number;
  taxTotal: number;
  roundOff: number;
  grandTotal: number;
}

export function previewTotals(lines: PreviewLine[], interState = false): PreviewTotals {
  let subTotal = 0;
  let discountTotal = 0;
  let cgst = 0;
  let sgst = 0;
  let igst = 0;

  for (const l of lines) {
    const taxable = Math.max(0, l.unitPrice * l.quantity - l.discount);
    subTotal += taxable;
    discountTotal += l.discount;
    if (interState) {
      igst += Math.round((taxable * l.gstRate) / 100);
    } else {
      cgst += Math.round((taxable * (l.gstRate / 2)) / 100);
      sgst += Math.round((taxable * (l.gstRate / 2)) / 100);
    }
  }

  const taxTotal = cgst + sgst + igst;
  const preRound = subTotal + taxTotal;
  const grandTotal = Math.round(preRound / 100) * 100;
  const roundOff = grandTotal - preRound;

  return { subTotal, discountTotal, cgst, sgst, igst, taxTotal, roundOff, grandTotal };
}
