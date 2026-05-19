import { apiFetch } from "./client";

type ApiResponse<T> = { message: string; data: T };

export type Promotion = {
  id: number;
  code: string;
  description: string;
  discount: number;
  type: "percentage" | "fixed";
  start_date: string;
  end_date: string;
  status: "active" | "scheduled" | "expired";
  usage_count: number;
};

export const getPromotions = (activeOnly = false) =>
  apiFetch<ApiResponse<Promotion[]>>(`/promotions${activeOnly ? "?active=true" : ""}`);

export const createPromotion = (payload: Omit<Promotion, "id" | "usage_count">) =>
  apiFetch<ApiResponse<Promotion>>("/promotions", {
    method: "POST",
    auth: true,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

export const updatePromotion = (id: number, payload: Omit<Promotion, "id" | "usage_count">) =>
  apiFetch<ApiResponse<Promotion>>(`/promotions/${id}`, {
    method: "PUT",
    auth: true,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

export const deletePromotion = (id: number) =>
  apiFetch<ApiResponse<null>>(`/promotions/${id}`, { method: "DELETE", auth: true });

export const applyPromotion = (code: string, subtotal: number) =>
  apiFetch<ApiResponse<{ discount: number; promotion: Promotion }>>("/promotions/apply", {
    method: "POST",
    auth: true,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ code, subtotal }),
  });
