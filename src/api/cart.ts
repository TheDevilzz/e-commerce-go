import { apiFetch } from "./client";
import type { ProductApiItem } from "./product";

export type CartItem = {
  id: number;
  product_id: number;
  product: ProductApiItem;
  quantity: number;
};

type ApiResponse<T> = { message: string; data: T };

export const notifyCartUpdated = () => {
  window.dispatchEvent(new Event("shophub-cart-updated"));
};

export const getCart = () => apiFetch<ApiResponse<CartItem[]>>("/cart", { auth: true });

export const addCartItem = (productId: number, quantity: number) =>
  apiFetch<ApiResponse<CartItem>>("/cart", {
    method: "POST",
    auth: true,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ product_id: productId, quantity }),
  });

export const updateCartItem = (id: number, productId: number, quantity: number) =>
  apiFetch<ApiResponse<CartItem>>(`/cart/${id}`, {
    method: "PUT",
    auth: true,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ product_id: productId, quantity }),
  });

export const deleteCartItem = (id: number) =>
  apiFetch<ApiResponse<null>>(`/cart/${id}`, {
    method: "DELETE",
    auth: true,
  });
