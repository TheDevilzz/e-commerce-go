import { API_BASE_URL, getAuthHeaders, getErrorMessage } from "./client";

export const uploadPaymentSlip = async (orderId: number, amount: number, image: File) => {
  const formData = new FormData();
  formData.append("order_id", String(orderId));
  formData.append("amount", String(amount));
  formData.append("image", image);

  const response = await fetch(`${API_BASE_URL}/payment`, {
    method: "POST",
    headers: getAuthHeaders(),
    body: formData,
  });

  if (!response.ok) {
    throw new Error(await getErrorMessage(response, "Failed to upload payment slip"));
  }

  return response.json();
};
