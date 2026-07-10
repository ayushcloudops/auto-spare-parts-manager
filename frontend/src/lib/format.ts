// Formatting helpers shared across the UI.
//
// The backend sends money as an integer number of paise (see Go money.Money).
// We format to Indian-style ₹ grouping here, at the display edge only.

const inr = new Intl.NumberFormat("en-IN", {
  style: "currency",
  currency: "INR",
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
});

/** formatPaise renders an integer paise amount as "₹1,23,456.78". */
export function formatPaise(paise: number): string {
  return inr.format((paise ?? 0) / 100);
}

/** rupeesToPaise converts a user-entered rupee value to integer paise. */
export function rupeesToPaise(rupees: number): number {
  return Math.round(rupees * 100);
}

const dateFmt = new Intl.DateTimeFormat("en-IN", {
  day: "2-digit",
  month: "short",
  year: "numeric",
});

/** formatDate renders an ISO/date string as "29 Jun 2026". */
export function formatDate(value: string | Date): string {
  const d = typeof value === "string" ? new Date(value) : value;
  if (isNaN(d.getTime())) return "—";
  return dateFmt.format(d);
}
