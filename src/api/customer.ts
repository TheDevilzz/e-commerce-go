import { apiFetch } from "./client";
import type { ProductApiItem } from "./product";

type ApiResponse<T> = { message: string; data: T };

export type WishlistItem = {
  id: number;
  product_id: number;
  product: ProductApiItem;
};

export type Address = {
  id: number;
  type: "home" | "work" | string;
  name: string;
  street: string;
  city: string;
  state: string;
  zip: string;
  is_default: boolean;
};

export type PaymentMethod = {
  id: number;
  type: "visa" | "mastercard" | "amex" | string;
  last4: string;
  expiry: string;
  is_default: boolean;
};

export const getWishlist = () => apiFetch<ApiResponse<WishlistItem[]>>("/wishlist", { auth: true });

export const addWishlistItem = (productId: number) =>
  apiFetch<ApiResponse<WishlistItem>>("/wishlist", {
    method: "POST",
    auth: true,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ product_id: productId }),
  });

export const deleteWishlistItem = (id: number) =>
  apiFetch<ApiResponse<null>>(`/wishlist/${id}`, { method: "DELETE", auth: true });

export const getAddresses = () => apiFetch<ApiResponse<Address[]>>("/addresses", { auth: true });

export const createAddress = (payload: Omit<Address, "id">) =>
  apiFetch<ApiResponse<Address>>("/addresses", {
    method: "POST",
    auth: true,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

export const updateAddress = (id: number, payload: Omit<Address, "id">) =>
  apiFetch<ApiResponse<Address>>(`/addresses/${id}`, {
    method: "PUT",
    auth: true,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

export const deleteAddress = (id: number) =>
  apiFetch<ApiResponse<null>>(`/addresses/${id}`, { method: "DELETE", auth: true });

export const getPaymentMethods = () =>
  apiFetch<ApiResponse<PaymentMethod[]>>("/payment-methods", { auth: true });

export const createPaymentMethod = (payload: Omit<PaymentMethod, "id">) =>
  apiFetch<ApiResponse<PaymentMethod>>("/payment-methods", {
    method: "POST",
    auth: true,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

export const updatePaymentMethod = (id: number, payload: Omit<PaymentMethod, "id">) =>
  apiFetch<ApiResponse<PaymentMethod>>(`/payment-methods/${id}`, {
    method: "PUT",
    auth: true,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

export const deletePaymentMethod = (id: number) =>
  apiFetch<ApiResponse<null>>(`/payment-methods/${id}`, { method: "DELETE", auth: true });
