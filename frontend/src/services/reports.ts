import {
  Sales,
  TopProducts,
  Profit,
  LowStock,
  ExportSalesCSV,
} from "../../wailsjs/go/app/ReportsHandler";
import { domain } from "../../wailsjs/go/models";

export type SalesSummary = domain.SalesSummary;
export type TopProduct = domain.TopProduct;
export type ProfitReport = domain.ProfitReport;
export type Product = domain.Product;

// Dates are "YYYY-MM-DD" strings.
export const reportsApi = {
  sales: (from: string, to: string) => Sales(from, to),
  topProducts: (from: string, to: string, limit = 10) => TopProducts(from, to, limit),
  profit: (from: string, to: string) => Profit(from, to),
  lowStock: () => LowStock(),
  exportCSV: (from: string, to: string) => ExportSalesCSV(from, to),
};
