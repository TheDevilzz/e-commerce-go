import { apiFetch } from "./client";

export type ApiUser = {
  id: number;
  username: string;
  email: string;
  name: string;
  phone: string;
  date_of_birth: string;
  role: "admin" | "user" | string;
  created_at?: string;
};

type ApiResponse<T> = {
  message: string;
  data: T;
};

export const getMe = async () => {
  return apiFetch<ApiResponse<ApiUser>>("/me", { auth: true });
};

export const updateMe = async (payload: {
  name: string;
  phone: string;
  date_of_birth: string;
}) => {
  return apiFetch<ApiResponse<ApiUser>>("/me", {
    method: "PUT",
    auth: true,
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });
};

export const getUsers = async () => {
  return apiFetch<ApiResponse<ApiUser[]>>("/users", { auth: true });
};
