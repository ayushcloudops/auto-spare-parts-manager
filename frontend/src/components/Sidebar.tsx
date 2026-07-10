import { NavLink } from "react-router-dom";
import clsx from "clsx";
import { Wrench } from "lucide-react";
import { navItems } from "../lib/nav";

/**
 * Sidebar is the primary navigation. Large, icon-led targets suit a
 * touch/keyboard shop counter. The active route is highlighted by NavLink.
 */
export default function Sidebar() {
  return (
    <aside className="flex w-60 shrink-0 flex-col border-r border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
      {/* Brand */}
      <div className="flex items-center gap-2 px-5 py-5">
        <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-brand-600 text-white">
          <Wrench size={18} />
        </div>
        <div className="leading-tight">
          <div className="text-sm font-semibold text-slate-800 dark:text-slate-100">
            AutoParts
          </div>
          <div className="text-[11px] text-slate-400">Shop Manager</div>
        </div>
      </div>

      {/* Nav */}
      <nav className="flex-1 space-y-1 px-3 py-2">
        {navItems.map((item) => {
          const Icon = item.icon;
          return (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.to === "/"}
              className={({ isActive }) =>
                clsx(
                  "group flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition",
                  isActive
                    ? "bg-brand-600 text-white shadow-sm"
                    : "text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800"
                )
              }
            >
              <Icon size={18} />
              <span className="flex-1">{item.label}</span>
              {item.shortcut && (
                <kbd className="hidden rounded border border-current px-1 text-[10px] opacity-40 group-hover:inline">
                  Alt+{item.shortcut}
                </kbd>
              )}
            </NavLink>
          );
        })}
      </nav>

      <div className="px-5 py-4 text-[11px] text-slate-400">
        v0.1.0 · Offline
      </div>
    </aside>
  );
}
