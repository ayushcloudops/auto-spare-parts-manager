// Single source of truth for the sidebar navigation and routing. Adding a
// module = adding one entry here plus its <Route> in App.tsx.
import {
  LayoutDashboard,
  Package,
  Receipt,
  History,
  Users,
  Truck,
  ClipboardList,
  BarChart3,
  Settings,
  LucideIcon,
} from "lucide-react";

export interface NavItem {
  to: string;
  label: string;
  icon: LucideIcon;
  /** Keyboard shortcut hint (Alt+key) shown in the sidebar. */
  shortcut?: string;
}

export const navItems: NavItem[] = [
  { to: "/", label: "Dashboard", icon: LayoutDashboard, shortcut: "1" },
  { to: "/products", label: "Products", icon: Package, shortcut: "2" },
  { to: "/billing", label: "Billing", icon: Receipt, shortcut: "3" },
  { to: "/invoices", label: "Invoices", icon: History, shortcut: "4" },
  { to: "/customers", label: "Customers", icon: Users, shortcut: "5" },
  { to: "/suppliers", label: "Suppliers", icon: Truck, shortcut: "6" },
  { to: "/purchases", label: "Purchases", icon: ClipboardList, shortcut: "7" },
  { to: "/reports", label: "Reports", icon: BarChart3, shortcut: "8" },
  { to: "/settings", label: "Settings", icon: Settings, shortcut: "9" },
];
