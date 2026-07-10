import { Stats } from "../../wailsjs/go/app/DashboardHandler";
import { domain } from "../../wailsjs/go/models";

export type DashboardStats = domain.DashboardStats;

export const dashboardApi = {
  stats: () => Stats(),
};
