import { useCallback, useEffect, useState } from "react";
import { Download, TrendingUp, AlertTriangle } from "lucide-react";
import { Button } from "../components/ui/Button";
import {
  reportsApi,
  SalesSummary,
  ProfitReport,
  TopProduct,
  Product,
} from "../services/reports";
import { formatPaise } from "../lib/format";

function todayStr(offsetDays = 0): string {
  const d = new Date();
  d.setDate(d.getDate() + offsetDays);
  return d.toISOString().slice(0, 10);
}

export default function Reports() {
  const [from, setFrom] = useState(todayStr(-30));
  const [to, setTo] = useState(todayStr(0));
  const [sales, setSales] = useState<SalesSummary | null>(null);
  const [profit, setProfit] = useState<ProfitReport | null>(null);
  const [top, setTop] = useState<TopProduct[]>([]);
  const [lowStock, setLowStock] = useState<Product[]>([]);
  const [exporting, setExporting] = useState(false);

  const load = useCallback(() => {
    reportsApi.sales(from, to).then(setSales).catch(() => {});
    reportsApi.profit(from, to).then(setProfit).catch(() => {});
    reportsApi.topProducts(from, to, 10).then(setTop).catch(() => {});
    reportsApi.lowStock().then(setLowStock).catch(() => {});
  }, [from, to]);

  useEffect(() => {
    load();
  }, [load]);

  const exportCSV = async () => {
    setExporting(true);
    try {
      const csv = await reportsApi.exportCSV(from, to);
      const blob = new Blob([csv], { type: "text/csv;charset=utf-8;" });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `sales_${from}_to_${to}.csv`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    } finally {
      setExporting(false);
    }
  };

  const setPreset = (days: number) => {
    setFrom(todayStr(-days));
    setTo(todayStr(0));
  };

  return (
    <div className="space-y-6">
      {/* controls */}
      <div className="flex flex-wrap items-end gap-3">
        <label className="block">
          <span className="mb-1 block text-xs font-medium text-slate-500">From</span>
          <input type="date" value={from} onChange={(e) => setFrom(e.target.value)} className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-800" />
        </label>
        <label className="block">
          <span className="mb-1 block text-xs font-medium text-slate-500">To</span>
          <input type="date" value={to} onChange={(e) => setTo(e.target.value)} className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-800" />
        </label>
        <div className="flex gap-1">
          <Button variant="secondary" onClick={() => setPreset(0)}>Today</Button>
          <Button variant="secondary" onClick={() => setPreset(7)}>7d</Button>
          <Button variant="secondary" onClick={() => setPreset(30)}>30d</Button>
        </div>
        <Button className="ml-auto" onClick={exportCSV} disabled={exporting}>
          <Download size={16} /> {exporting ? "Exporting…" : "Export CSV"}
        </Button>
      </div>

      {/* summary cards */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <Stat label="Total Sales" value={formatPaise(sales?.totalSales ?? 0)} tone="text-emerald-600" />
        <Stat label="Invoices" value={String(sales?.invoiceCount ?? 0)} tone="text-brand-600" />
        <Stat label="GST Collected" value={formatPaise(sales?.totalTax ?? 0)} tone="text-sky-600" />
        <Stat label="Profit" value={formatPaise(profit?.profit ?? 0)} tone="text-amber-600" />
      </div>

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        {/* top products */}
        <div className="card overflow-hidden">
          <div className="flex items-center gap-2 border-b border-slate-200 px-5 py-3 dark:border-slate-800">
            <TrendingUp size={16} className="text-brand-600" />
            <h2 className="font-semibold text-slate-800 dark:text-slate-100">Top Selling Products</h2>
          </div>
          <table className="w-full text-sm">
            <thead className="bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500 dark:bg-slate-800/50 dark:text-slate-400">
              <tr>
                <th className="px-5 py-2.5">Product</th>
                <th className="px-5 py-2.5 text-right">Qty Sold</th>
                <th className="px-5 py-2.5 text-right">Revenue</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
              {top.length === 0 ? (
                <tr><td colSpan={3} className="px-5 py-8 text-center text-slate-400">No sales in this period.</td></tr>
              ) : (
                top.map((t) => (
                  <tr key={t.productId}>
                    <td className="px-5 py-2.5 text-slate-700 dark:text-slate-200">{t.productName}</td>
                    <td className="px-5 py-2.5 text-right text-slate-500">{t.qtySold}</td>
                    <td className="px-5 py-2.5 text-right font-medium text-slate-700 dark:text-slate-200">{formatPaise(t.revenue)}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* low stock */}
        <div className="card overflow-hidden">
          <div className="flex items-center gap-2 border-b border-slate-200 px-5 py-3 dark:border-slate-800">
            <AlertTriangle size={16} className="text-amber-600" />
            <h2 className="font-semibold text-slate-800 dark:text-slate-100">Low Stock Report</h2>
          </div>
          <table className="w-full text-sm">
            <thead className="bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500 dark:bg-slate-800/50 dark:text-slate-400">
              <tr>
                <th className="px-5 py-2.5">Product</th>
                <th className="px-5 py-2.5 text-right">Stock</th>
                <th className="px-5 py-2.5 text-right">Min</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
              {lowStock.length === 0 ? (
                <tr><td colSpan={3} className="px-5 py-8 text-center text-slate-400">All products above minimum. 🎉</td></tr>
              ) : (
                lowStock.map((p) => (
                  <tr key={p.id}>
                    <td className="px-5 py-2.5 text-slate-700 dark:text-slate-200">{p.name}</td>
                    <td className="px-5 py-2.5 text-right font-medium text-amber-600">{p.currentStock}</td>
                    <td className="px-5 py-2.5 text-right text-slate-500">{p.minimumStock}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function Stat({ label, value, tone }: { label: string; value: string; tone: string }) {
  return (
    <div className="card p-5">
      <div className="text-sm text-slate-500 dark:text-slate-400">{label}</div>
      <div className={"mt-2 text-2xl font-semibold " + tone}>{value}</div>
    </div>
  );
}
