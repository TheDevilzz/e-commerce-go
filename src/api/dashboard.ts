import { apiFetch } from "./client";

type ApiResponse<T> = { message: string; data: T };

export type DashboardStats = {
  total_revenue: number;
  orders: number;
  customers: number;
  products: number;
  avg_order: number;
  category_sales: Array<{ category: string; sales: number }>;
};

export const getDashboardStats = () =>
  apiFetch<ApiResponse<DashboardStats>>("/dashboard/stats", { auth: true });
