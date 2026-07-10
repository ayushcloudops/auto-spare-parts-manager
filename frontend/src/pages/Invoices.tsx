import { useCallback, useEffect, useState } from "react";
import { Search, Eye, Printer } from "lucide-react";
import { Button } from "../components/ui/Button";
import { Modal } from "../components/ui/Modal";
import { billingApi, Invoice } from "../services/billing";
import { printApi } from "../services/print";
import { formatPaise, formatDate } from "../lib/format";

export default function Invoices() {
  const [items, setItems] = useState<Invoice[]>([]);
  const [search, setSearch] = useState("");
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");
  const [loading, setLoading] = useState(true);
  const [viewing, setViewing] = useState<Invoice | null>(null);
  const [receipt, setReceipt] = useState("");
  const [printMsg, setPrintMsg] = useState("");

  const load = useCallback(() => {
    setLoading(true);
    const filter: any = { search };
    if (from) filter.from = new Date(from).toISOString();
    if (to) filter.to = new Date(to + "T23:59:59").toISOString();
    billingApi.list(filter).then(setItems).catch(() => setItems([])).finally(() => setLoading(false));
  }, [search, from, to]);

  useEffect(() => {
    const t = setTimeout(load, 200);
    return () => clearTimeout(t);
  }, [load]);

  const openView = async (inv: Invoice) => {
    setViewing(inv);
    setPrintMsg("");
    setReceipt("");
    try {
      setReceipt(await printApi.preview(inv.id));
    } catch {
      setReceipt("(could not load receipt preview)");
    }
  };

  const doPrint = async () => {
    if (!viewing) return;
    setPrintMsg("Printing…");
    try {
      await printApi.print(viewing.id);
      setPrintMsg("Sent to printer (or saved to receipts folder if no printer).");
    } catch (e: any) {
      setPrintMsg(String(e?.message ?? e));
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center gap-3">
        <div className="relative min-w-[220px] flex-1">
          <Search size={16} className="pointer-events-none absolute left-3 top-2.5 text-slate-400" />
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search invoice number or customer…"
            className="w-full rounded-lg border border-slate-300 bg-white py-2 pl-9 pr-3 text-sm outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-500/20 dark:border-slate-700 dark:bg-slate-800"
          />
        </div>
        <input type="date" value={from} onChange={(e) => setFrom(e.target.value)} className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-800" />
        <span className="text-slate-400">to</span>
        <input type="date" value={to} onChange={(e) => setTo(e.target.value)} className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-800" />
      </div>

      <div className="card overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500 dark:bg-slate-800/50 dark:text-slate-400">
            <tr>
              <th className="px-4 py-3">Invoice</th>
              <th className="px-4 py-3">Date</th>
              <th className="px-4 py-3">Customer</th>
              <th className="px-4 py-3">Payment</th>
              <th className="px-4 py-3 text-right">Total</th>
              <th className="px-4 py-3 text-right">Due</th>
              <th className="px-4 py-3"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
            {loading ? (
              <tr><td colSpan={7} className="px-4 py-10 text-center text-slate-400">Loading…</td></tr>
            ) : items.length === 0 ? (
              <tr><td colSpan={7} className="px-4 py-10 text-center text-slate-400">No invoices found.</td></tr>
            ) : (
              items.map((inv) => (
                <tr key={inv.id} className="hover:bg-slate-50 dark:hover:bg-slate-800/40">
                  <td className="px-4 py-3 font-medium text-slate-800 dark:text-slate-100">{inv.number}</td>
                  <td className="px-4 py-3 text-slate-500">{formatDate(inv.date)}</td>
                  <td className="px-4 py-3 text-slate-500">{inv.customer?.name ?? "Walk-in"}</td>
                  <td className="px-4 py-3 text-slate-500">{String(inv.paymentMode).toUpperCase()}</td>
                  <td className="px-4 py-3 text-right font-medium text-slate-700 dark:text-slate-200">{formatPaise(inv.grandTotal)}</td>
                  <td className={"px-4 py-3 text-right " + (inv.amountDue > 0 ? "text-rose-600" : "text-slate-400")}>{formatPaise(inv.amountDue)}</td>
                  <td className="px-4 py-3 text-right">
                    <button onClick={() => openView(inv)} className="rounded-md p-1.5 text-slate-400 hover:text-brand-600" title="View"><Eye size={16} /></button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <Modal
        open={!!viewing}
        onClose={() => setViewing(null)}
        title={viewing ? `Invoice ${viewing.number}` : ""}
        footer={
          <>
            {printMsg && <span className="mr-auto text-xs text-slate-500">{printMsg}</span>}
            <Button variant="secondary" onClick={() => setViewing(null)}>Close</Button>
            <Button onClick={doPrint}><Printer size={16} /> Print / Reprint</Button>
          </>
        }
      >
        <pre className="max-h-[60vh] overflow-auto rounded-lg bg-slate-50 p-4 font-mono text-xs leading-relaxed text-slate-700 dark:bg-slate-800 dark:text-slate-200">
          {receipt || "Loading receipt…"}
        </pre>
      </Modal>
    </div>
  );
}
