import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  IndianRupee,
  Receipt,
  Package,
  AlertTriangle,
  CreditCard,
  LucideIcon,
} from "lucide-react";
import { dashboardApi, DashboardStats } from "../services/dashboard";
import { formatPaise, formatDate } from "../lib/format";

export default function Dashboard() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    dashboardApi.stats().then(setStats).catch(() => {});
  }, []);

  const cards: { label: string; value: string; icon: LucideIcon; tone: string }[] = [
    { label: "Today's Sales", value: formatPaise(stats?.todaySales ?? 0), icon: IndianRupee, tone: "text-emerald-600" },
    { label: "Today's Bills", value: String(stats?.todayBills ?? 0), icon: Receipt, tone: "text-brand-600" },
    { label: "Total Products", value: String(stats?.totalProducts ?? 0), icon: Package, tone: "text-sky-600" },
    { label: "Low Stock", value: String(stats?.lowStockCount ?? 0), icon: AlertTriangle, tone: "text-amber-600" },
    { label: "Pending Credit", value: formatPaise(stats?.pendingCredit ?? 0), icon: CreditCard, tone: "text-rose-600" },
  ];

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-5">
        {cards.map((s) => {
          const Icon = s.icon;
          return (
            <div key={s.label} className="card p-5">
              <div className="flex items-center justify-between">
                <span className="text-sm text-slate-500 dark:text-slate-400">{s.label}</span>
                <Icon size={18} className={s.tone} />
              </div>
              <div className="mt-3 text-2xl font-semibold text-slate-800 dark:text-slate-100">
                {s.value}
              </div>
            </div>
          );
        })}
      </div>

      <div className="card overflow-hidden">
        <div className="border-b border-slate-200 px-5 py-3 dark:border-slate-800">
          <h2 className="font-semibold text-slate-800 dark:text-slate-100">Recent Bills</h2>
        </div>
        <table className="w-full text-sm">
          <thead className="bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500 dark:bg-slate-800/50 dark:text-slate-400">
            <tr>
              <th className="px-5 py-2.5">Invoice</th>
              <th className="px-5 py-2.5">Date</th>
              <th className="px-5 py-2.5">Customer</th>
              <th className="px-5 py-2.5">Payment</th>
              <th className="px-5 py-2.5 text-right">Total</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
            {!stats?.recentInvoices?.length ? (
              <tr>
                <td colSpan={5} className="px-5 py-8 text-center text-slate-400">
                  No bills yet. Create one from the Billing screen.
                </td>
              </tr>
            ) : (
              stats.recentInvoices.map((inv) => (
                <tr
                  key={inv.id}
                  onClick={() => navigate("/invoices")}
                  className="cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800/40"
                >
                  <td className="px-5 py-2.5 font-medium text-slate-700 dark:text-slate-200">{inv.number}</td>
                  <td className="px-5 py-2.5 text-slate-500">{formatDate(inv.date)}</td>
                  <td className="px-5 py-2.5 text-slate-500">{inv.customer?.name ?? "Walk-in"}</td>
                  <td className="px-5 py-2.5 text-slate-500">{String(inv.paymentMode).toUpperCase()}</td>
                  <td className="px-5 py-2.5 text-right font-medium text-slate-700 dark:text-slate-200">
                    {formatPaise(inv.grandTotal)}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
