import { apiFetch } from "./client";

export type ApiOrderItem = {
  id: number;
  order_id: number;
  product_id: number;
  product_name: string;
  quantity: number;
  price: number;
};

export type ApiOrder = {
  id: number;
  user_id: number;
  total_price: number;
  status: string;
  items: ApiOrderItem[];
  CreatedAt?: string;
  created_at?: string;
};

type ApiResponse<T> = {
  message: string;
  data: T;
};

export const getMyOrders = async () => {
  return apiFetch<ApiResponse<ApiOrder[]>>("/orders/my", { auth: true });
};

export const getAdminOrders = async () => {
  return apiFetch<ApiResponse<ApiOrder[]>>("/orders", { auth: true });
};

export const createOrder = async (items: Array<{ product_id: number; quantity: number }>, addressId?: number) =>
  apiFetch<ApiResponse<ApiOrder>>("/orders", {
    method: "POST",
    auth: true,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ items, ...(addressId ? { address_id: addressId } : {}) }),
  });
