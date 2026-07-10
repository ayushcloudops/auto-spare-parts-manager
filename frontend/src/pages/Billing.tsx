import { useEffect, useMemo, useState } from "react";
import { Search, Plus, Trash2, Receipt, CheckCircle2 } from "lucide-react";
import { Button } from "../components/ui/Button";
import { Select } from "../components/ui/form";
import { Product, productsApi } from "../services/products";
import { customersApi, Customer } from "../services/customers";
import { billingApi, Invoice, NewBill, PAYMENT_MODES } from "../services/billing";
import { previewTotals } from "../lib/gst";
import { formatPaise, rupeesToPaise } from "../lib/format";

interface CartLine {
  product: Product;
  quantity: number;
  discountRupees: string;
}

export default function Billing() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<Product[]>([]);
  const [cart, setCart] = useState<CartLine[]>([]);
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [customerId, setCustomerId] = useState<number | "">("");
  const [paymentMode, setPaymentMode] = useState("cash");
  const [amountPaidRupees, setAmountPaidRupees] = useState("");
  const [notes, setNotes] = useState("");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [saved, setSaved] = useState<Invoice | null>(null);

  // Product search (debounced).
  useEffect(() => {
    const t = setTimeout(() => {
      productsApi.list({ search: query, limit: 8 }).then(setResults).catch(() => setResults([]));
    }, 150);
    return () => clearTimeout(t);
  }, [query]);

  useEffect(() => {
    customersApi.list("").then(setCustomers).catch(() => {});
  }, []);

  const addToCart = (p: Product) => {
    setCart((prev) => {
      const existing = prev.find((l) => l.product.id === p.id);
      if (existing) {
        return prev.map((l) =>
          l.product.id === p.id
            ? { ...l, quantity: Math.min(l.quantity + 1, p.currentStock) }
            : l
        );
      }
      return [...prev, { product: p, quantity: 1, discountRupees: "" }];
    });
  };

  const updateLine = (id: number, patch: Partial<CartLine>) =>
    setCart((prev) => prev.map((l) => (l.product.id === id ? { ...l, ...patch } : l)));
  const removeLine = (id: number) =>
    setCart((prev) => prev.filter((l) => l.product.id !== id));

  // Live totals preview (authoritative totals come back from the backend).
  const totals = useMemo(
    () =>
      previewTotals(
        cart.map((l) => ({
          unitPrice: l.product.sellingPrice,
          quantity: l.quantity,
          discount: rupeesToPaise(parseFloat(l.discountRupees) || 0),
          gstRate: l.product.gstRate,
        }))
      ),
    [cart]
  );

  const canSubmit = cart.length > 0 && !saving;

  const generate = async () => {
    if (!canSubmit) return;
    setSaving(true);
    setError("");
    const bill: NewBill = {
      customerId: customerId === "" ? undefined : Number(customerId),
      lines: cart.map((l) => ({
        productId: l.product.id,
        quantity: l.quantity,
        discount: rupeesToPaise(parseFloat(l.discountRupees) || 0),
      })),
      paymentMode,
      amountPaid:
        paymentMode === "credit" ? rupeesToPaise(parseFloat(amountPaidRupees) || 0) : 0,
      notes,
    };
    try {
      const inv = await billingApi.create(bill);
      setSaved(inv);
    } catch (e: any) {
      setError(String(e?.message ?? e));
    } finally {
      setSaving(false);
    }
  };

  const newBill = () => {
    setCart([]);
    setCustomerId("");
    setPaymentMode("cash");
    setAmountPaidRupees("");
    setNotes("");
    setSaved(null);
    setError("");
    setQuery("");
  };

  // --- success view ---------------------------------------------------------
  if (saved) {
    return (
      <div className="mx-auto max-w-lg">
        <div className="card p-8 text-center">
          <CheckCircle2 className="mx-auto text-emerald-500" size={48} />
          <h2 className="mt-4 text-lg font-semibold text-slate-800 dark:text-slate-100">
            Invoice {saved.number} saved
          </h2>
          <div className="mt-1 text-sm text-slate-500">
            {saved.items?.length} item(s) · {saved.paymentMode.toUpperCase()}
          </div>
          <div className="mt-4 text-3xl font-bold text-slate-800 dark:text-slate-100">
            {formatPaise(saved.grandTotal)}
          </div>
          {saved.amountDue > 0 && (
            <div className="mt-1 text-sm text-rose-600">
              Credit due: {formatPaise(saved.amountDue)}
            </div>
          )}
          <div className="mt-6 flex justify-center gap-2">
            <Button onClick={newBill}>
              <Plus size={16} /> New Bill
            </Button>
          </div>
          <p className="mt-4 text-xs text-slate-400">
            Thermal receipt printing arrives in the Printing step.
          </p>
        </div>
      </div>
    );
  }

  // --- billing view ---------------------------------------------------------
  return (
    <div className="grid grid-cols-1 gap-4 lg:grid-cols-5">
      {/* Left: product search */}
      <div className="lg:col-span-2 space-y-3">
        <div className="relative">
          <Search size={16} className="pointer-events-none absolute left-3 top-2.5 text-slate-400" />
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search product or part number…"
            autoFocus
            className="w-full rounded-lg border border-slate-300 bg-white py-2 pl-9 pr-3 text-sm outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-500/20 dark:border-slate-700 dark:bg-slate-800"
          />
        </div>
        <div className="card divide-y divide-slate-100 dark:divide-slate-800">
          {results.length === 0 ? (
            <div className="p-6 text-center text-sm text-slate-400">No products found</div>
          ) : (
            results.map((p) => {
              const out = p.currentStock <= 0;
              return (
                <button
                  key={p.id}
                  disabled={out}
                  onClick={() => addToCart(p)}
                  className="flex w-full items-center justify-between px-4 py-3 text-left hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-50 dark:hover:bg-slate-800/40"
                >
                  <div>
                    <div className="text-sm font-medium text-slate-800 dark:text-slate-100">{p.name}</div>
                    <div className="text-xs text-slate-400">
                      {p.partNumber || "—"} · Stock {p.currentStock}
                    </div>
                  </div>
                  <div className="text-sm font-medium text-slate-700 dark:text-slate-200">
                    {formatPaise(p.sellingPrice)}
                  </div>
                </button>
              );
            })
          )}
        </div>
      </div>

      {/* Right: cart + checkout */}
      <div className="lg:col-span-3 space-y-4">
        <div className="card overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500 dark:bg-slate-800/50 dark:text-slate-400">
              <tr>
                <th className="px-3 py-2">Item</th>
                <th className="px-3 py-2 w-20 text-center">Qty</th>
                <th className="px-3 py-2 w-24 text-right">Disc ₹</th>
                <th className="px-3 py-2 text-right">Amount</th>
                <th className="px-3 py-2"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
              {cart.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-3 py-10 text-center text-slate-400">
                    Search and click a product to add it to the bill.
                  </td>
                </tr>
              ) : (
                cart.map((l) => {
                  const disc = rupeesToPaise(parseFloat(l.discountRupees) || 0);
                  const amount = l.product.sellingPrice * l.quantity - disc;
                  return (
                    <tr key={l.product.id}>
                      <td className="px-3 py-2">
                        <div className="font-medium text-slate-800 dark:text-slate-100">{l.product.name}</div>
                        <div className="text-xs text-slate-400">
                          {formatPaise(l.product.sellingPrice)} · {l.product.gstRate}% GST
                        </div>
                      </td>
                      <td className="px-3 py-2">
                        <input
                          type="number"
                          min={1}
                          max={l.product.currentStock}
                          value={l.quantity}
                          onChange={(e) =>
                            updateLine(l.product.id, {
                              quantity: Math.max(1, Math.min(parseInt(e.target.value, 10) || 1, l.product.currentStock)),
                            })
                          }
                          className="w-16 rounded border border-slate-300 bg-white px-2 py-1 text-center text-sm dark:border-slate-700 dark:bg-slate-800"
                        />
                      </td>
                      <td className="px-3 py-2">
                        <input
                          type="number"
                          min={0}
                          value={l.discountRupees}
                          onChange={(e) => updateLine(l.product.id, { discountRupees: e.target.value })}
                          placeholder="0"
                          className="w-20 rounded border border-slate-300 bg-white px-2 py-1 text-right text-sm dark:border-slate-700 dark:bg-slate-800"
                        />
                      </td>
                      <td className="px-3 py-2 text-right font-medium text-slate-700 dark:text-slate-200">
                        {formatPaise(Math.max(0, amount))}
                      </td>
                      <td className="px-3 py-2 text-right">
                        <button onClick={() => removeLine(l.product.id)} className="rounded p-1 text-slate-400 hover:text-rose-600">
                          <Trash2 size={15} />
                        </button>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>

        {/* Checkout controls */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="space-y-3">
            <label className="block">
              <span className="mb-1 block text-xs font-medium text-slate-500">Customer</span>
              <Select value={customerId} onChange={(e) => setCustomerId(e.target.value === "" ? "" : Number(e.target.value))}>
                <option value="">Walk-in customer</option>
                {customers.map((c) => (
                  <option key={c.id} value={c.id}>
                    {c.name} {c.phone ? `(${c.phone})` : ""}
                  </option>
                ))}
              </Select>
            </label>
            <label className="block">
              <span className="mb-1 block text-xs font-medium text-slate-500">Payment Mode</span>
              <div className="flex gap-2">
                {PAYMENT_MODES.map((m) => (
                  <button
                    key={m.value}
                    onClick={() => setPaymentMode(m.value)}
                    className={
                      "flex-1 rounded-lg border px-2 py-1.5 text-sm font-medium transition " +
                      (paymentMode === m.value
                        ? "border-brand-600 bg-brand-600 text-white"
                        : "border-slate-300 text-slate-600 hover:bg-slate-100 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800")
                    }
                  >
                    {m.label}
                  </button>
                ))}
              </div>
            </label>
            {paymentMode === "credit" && (
              <label className="block">
                <span className="mb-1 block text-xs font-medium text-slate-500">Amount Paid Now (₹)</span>
                <input
                  type="number"
                  min={0}
                  value={amountPaidRupees}
                  onChange={(e) => setAmountPaidRupees(e.target.value)}
                  placeholder="0"
                  className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-800"
                />
              </label>
            )}
          </div>

          {/* Totals */}
          <div className="card space-y-1.5 p-4 text-sm">
            <Row label="Subtotal" value={formatPaise(totals.subTotal)} />
            {totals.discountTotal > 0 && <Row label="Discount" value={"− " + formatPaise(totals.discountTotal)} />}
            {totals.cgst > 0 && <Row label="CGST" value={formatPaise(totals.cgst)} />}
            {totals.sgst > 0 && <Row label="SGST" value={formatPaise(totals.sgst)} />}
            {totals.igst > 0 && <Row label="IGST" value={formatPaise(totals.igst)} />}
            {totals.roundOff !== 0 && <Row label="Round Off" value={formatPaise(totals.roundOff)} />}
            <div className="mt-2 flex items-center justify-between border-t border-slate-200 pt-2 dark:border-slate-700">
              <span className="font-semibold text-slate-800 dark:text-slate-100">Grand Total</span>
              <span className="text-xl font-bold text-slate-800 dark:text-slate-100">
                {formatPaise(totals.grandTotal)}
              </span>
            </div>
          </div>
        </div>

        {error && (
          <div className="rounded-lg bg-rose-50 px-3 py-2 text-sm text-rose-700 dark:bg-rose-900/30 dark:text-rose-300">
            {error}
          </div>
        )}

        <div className="flex justify-end">
          <Button onClick={generate} disabled={!canSubmit} className="px-6 py-2.5 text-base">
            <Receipt size={18} /> {saving ? "Saving…" : "Generate Invoice"}
          </Button>
        </div>
      </div>
    </div>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between text-slate-600 dark:text-slate-300">
      <span>{label}</span>
      <span>{value}</span>
    </div>
  );
}
