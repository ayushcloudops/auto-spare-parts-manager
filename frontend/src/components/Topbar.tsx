import { useLocation } from "react-router-dom";
import clsx from "clsx";
import { navItems } from "../lib/nav";
import { useHealth } from "../hooks/useHealth";
import ThemeToggle from "./ThemeToggle";

/**
 * Topbar shows the current page title, a live backend-health indicator (proves
 * the frontend↔Go↔DB round-trip), and global controls.
 */
export default function Topbar() {
  const { pathname } = useLocation();
  const { info, error } = useHealth();
  const current =
    navItems.find((n) => n.to === pathname) ??
    (pathname === "/" ? navItems[0] : undefined);

  const healthy = !error && info?.status === "ok";

  return (
    <header className="flex h-14 shrink-0 items-center justify-between border-b border-slate-200 bg-white px-6 dark:border-slate-800 dark:bg-slate-900">
      <h1 className="text-base font-semibold text-slate-800 dark:text-slate-100">
        {current?.label ?? "Auto Spare Parts Manager"}
      </h1>
      <div className="flex items-center gap-4">
        {/* Backend health + shop identity, fetched from Go on mount. */}
        <div className="flex items-center gap-2 text-sm text-slate-500 dark:text-slate-400">
          <span
            className={clsx(
              "h-2 w-2 rounded-full",
              healthy ? "bg-emerald-500" : "bg-slate-300"
            )}
            title={healthy ? "Backend connected" : "Connecting…"}
          />
          <span>{info?.shopName ?? "…"}</span>
        </div>
        <ThemeToggle />
      </div>
    </header>
  );
}
