import { useCallback, useEffect, useMemo, useState } from "react";
import { Search, Trash2, PackagePlus, CheckCircle2 } from "lucide-react";
import { Button } from "../components/ui/Button";
import { Select } from "../components/ui/form";
import { Product, productsApi } from "../services/products";
import { Supplier, suppliersApi } from "../services/suppliers";
import { purchasesApi, Purchase, NewPurchase } from "../services/purchases";
import { formatPaise, formatDate, rupeesToPaise } from "../lib/format";

interface Line {
  product: Product;
  quantity: number;
  costRupees: string;
}

export default function Purchases() {
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [supplierId, setSupplierId] = useState<number | "">("");
  const [invNo, setInvNo] = useState("");
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<Product[]>([]);
  const [lines, setLines] = useState<Line[]>([]);
  const [notes, setNotes] = useState("");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [saved, setSaved] = useState(false);
  const [recent, setRecent] = useState<Purchase[]>([]);

  const loadRecent = useCallback(() => {
    purchasesApi.list({ limit: 10 }).then(setRecent).catch(() => {});
  }, []);

  useEffect(() => {
    suppliersApi.list("").then(setSuppliers).catch(() => {});
    loadRecent();
  }, [loadRecent]);

  useEffect(() => {
    const t = setTimeout(() => {
      productsApi.list({ search: query, limit: 8 }).then(setResults).catch(() => setResults([]));
    }, 150);
    return () => clearTimeout(t);
  }, [query]);

  const addLine = (p: Product) => {
    setLines((prev) =>
      prev.find((l) => l.product.id === p.id)
        ? prev
        : [...prev, { product: p, quantity: 1, costRupees: String(p.purchasePrice / 100 || "") }]
    );
  };
  const updateLine = (id: number, patch: Partial<Line>) =>
    setLines((prev) => prev.map((l) => (l.product.id === id ? { ...l, ...patch } : l)));
  const removeLine = (id: number) => setLines((prev) => prev.filter((l) => l.product.id !== id));

  const grandTotal = useMemo(
    () =>
      lines.reduce((sum, l) => {
        const cost = rupeesToPaise(parseFloat(l.costRupees) || 0);
        const net = cost * l.quantity;
        return sum + net + Math.round((net * l.product.gstRate) / 100);
      }, 0),
    [lines]
  );

  const canSave = supplierId !== "" && lines.length > 0 && !saving;

  const save = async () => {
    if (!canSave) return;
    setSaving(true);
    setError("");
    const payload: NewPurchase = {
      supplierId: Number(supplierId),
      supplierInvNo: invNo,
      lines: lines.map((l) => ({
        productId: l.product.id,
        quantity: l.quantity,
        costPrice: rupeesToPaise(parseFloat(l.costRupees) || 0),
      })),
      notes,
    };
    try {
      await purchasesApi.create(payload);
      setSaved(true);
      setLines([]);
      setInvNo("");
      setNotes("");
      setQuery("");
      loadRecent();
      setTimeout(() => setSaved(false), 2500);
    } catch (e: any) {
      setError(String(e?.message ?? e));
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 gap-4 lg:grid-cols-5">
        {/* product search */}
        <div className="lg:col-span-2 space-y-3">
          <div className="relative">
            <Search size={16} className="pointer-events-none absolute left-3 top-2.5 text-slate-400" />
            <input value={query} onChange={(e) => setQuery(e.target.value)} placeholder="Search product to receive…"
              className="w-full rounded-lg border border-slate-300 bg-white py-2 pl-9 pr-3 text-sm outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-500/20 dark:border-slate-700 dark:bg-slate-800" />
          </div>
          <div className="card divide-y divide-slate-100 dark:divide-slate-800">
            {results.length === 0 ? (
              <div className="p-6 text-center text-sm text-slate-400">No products</div>
            ) : (
              results.map((p) => (
                <button key={p.id} onClick={() => addLine(p)} className="flex w-full items-center justify-between px-4 py-3 text-left hover:bg-slate-50 dark:hover:bg-slate-800/40">
                  <div>
                    <div className="text-sm font-medium text-slate-800 dark:text-slate-100">{p.name}</div>
                    <div className="text-xs text-slate-400">Stock {p.currentStock}</div>
                  </div>
                  <div className="text-sm text-slate-500">{formatPaise(p.purchasePrice)}</div>
                </button>
              ))
            )}
          </div>
        </div>

        {/* purchase form */}
        <div className="lg:col-span-3 space-y-4">
          <div className="grid grid-cols-2 gap-3">
            <label className="block">
              <span className="mb-1 block text-xs font-medium text-slate-500">Supplier *</span>
              <Select value={supplierId} onChange={(e) => setSupplierId(e.target.value === "" ? "" : Number(e.target.value))}>
                <option value="">Select supplier…</option>
                {suppliers.map((s) => <option key={s.id} value={s.id}>{s.name}</option>)}
              </Select>
            </label>
            <label className="block">
              <span className="mb-1 block text-xs font-medium text-slate-500">Supplier Invoice No</span>
              <input value={invNo} onChange={(e) => setInvNo(e.target.value)} className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-800" />
            </label>
          </div>

          <div className="card overflow-hidden">
            <table className="w-full text-sm">
              <thead className="bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500 dark:bg-slate-800/50 dark:text-slate-400">
                <tr>
                  <th className="px-3 py-2">Item</th>
                  <th className="px-3 py-2 w-20 text-center">Qty</th>
                  <th className="px-3 py-2 w-28 text-right">Cost ₹</th>
                  <th className="px-3 py-2"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
                {lines.length === 0 ? (
                  <tr><td colSpan={4} className="px-3 py-8 text-center text-slate-400">Add products received from the supplier.</td></tr>
                ) : (
                  lines.map((l) => (
                    <tr key={l.product.id}>
                      <td className="px-3 py-2 font-medium text-slate-800 dark:text-slate-100">{l.product.name}</td>
                      <td className="px-3 py-2">
                        <input type="number" min={1} value={l.quantity}
                          onChange={(e) => updateLine(l.product.id, { quantity: Math.max(1, parseInt(e.target.value, 10) || 1) })}
                          className="w-16 rounded border border-slate-300 bg-white px-2 py-1 text-center text-sm dark:border-slate-700 dark:bg-slate-800" />
                      </td>
                      <td className="px-3 py-2">
                        <input type="number" min={0} value={l.costRupees}
                          onChange={(e) => updateLine(l.product.id, { costRupees: e.target.value })}
                          className="w-24 rounded border border-slate-300 bg-white px-2 py-1 text-right text-sm dark:border-slate-700 dark:bg-slate-800" />
                      </td>
                      <td className="px-3 py-2 text-right">
                        <button onClick={() => removeLine(l.product.id)} className="rounded p-1 text-slate-400 hover:text-rose-600"><Trash2 size={15} /></button>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>

          {error && <div className="rounded-lg bg-rose-50 px-3 py-2 text-sm text-rose-700 dark:bg-rose-900/30 dark:text-rose-300">{error}</div>}

          <div className="flex items-center justify-between">
            <div className="text-sm text-slate-500">
              Est. total <span className="font-semibold text-slate-800 dark:text-slate-100">{formatPaise(grandTotal)}</span>
            </div>
            <div className="flex items-center gap-3">
              {saved && <span className="flex items-center gap-1 text-sm text-emerald-600"><CheckCircle2 size={16} /> Saved · stock updated</span>}
              <Button onClick={save} disabled={!canSave}><PackagePlus size={16} /> {saving ? "Saving…" : "Save Purchase"}</Button>
            </div>
          </div>
        </div>
      </div>

      {/* recent purchases */}
      <div className="card overflow-hidden">
        <div className="border-b border-slate-200 px-5 py-3 dark:border-slate-800">
          <h2 className="font-semibold text-slate-800 dark:text-slate-100">Recent Purchases</h2>
        </div>
        <table className="w-full text-sm">
          <thead className="bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500 dark:bg-slate-800/50 dark:text-slate-400">
            <tr>
              <th className="px-5 py-2.5">Date</th>
              <th className="px-5 py-2.5">Supplier</th>
              <th className="px-5 py-2.5">Invoice No</th>
              <th className="px-5 py-2.5 text-right">Total</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
            {recent.length === 0 ? (
              <tr><td colSpan={4} className="px-5 py-8 text-center text-slate-400">No purchases yet.</td></tr>
            ) : (
              recent.map((p) => (
                <tr key={p.id}>
                  <td className="px-5 py-2.5 text-slate-500">{formatDate(p.date)}</td>
                  <td className="px-5 py-2.5 text-slate-700 dark:text-slate-200">{p.supplier?.name ?? "—"}</td>
                  <td className="px-5 py-2.5 text-slate-500">{p.supplierInvNo || "—"}</td>
                  <td className="px-5 py-2.5 text-right font-medium text-slate-700 dark:text-slate-200">{formatPaise(p.grandTotal)}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
