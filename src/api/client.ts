export const API_BASE_URL: string =
  import.meta.env.MODE === "production" ? "/api" : "http://localhost:3001/api";

export const ASSET_BASE_URL: string =
  import.meta.env.MODE === "production" ? "" : "http://localhost:3001";

export const getErrorMessage = async (response: Response, fallback: string) => {
  try {
    const errorData = await response.json();
    return errorData.message || errorData.error || fallback;
  } catch {
    return fallback;
  }
};

export const getAuthHeaders = () => {
  const token = localStorage.getItem("shophub_token");

  if (!token) {
    throw new Error("Please login again");
  }

  return { Authorization: `Bearer ${token}` };
};

type RefreshPayload = {
  token?: string;
  access_token?: string;
  refresh_token?: string;
  user?: unknown;
};

type ApiFetchOptions = RequestInit & {
  auth?: boolean;
  retryOnUnauthorized?: boolean;
};

const clearAuthStorage = () => {
  localStorage.removeItem("shophub_user");
  localStorage.removeItem("shophub_token");
  localStorage.removeItem("shophub_refresh_token");
};

const refreshAccessToken = async () => {
  const refreshToken = localStorage.getItem("shophub_refresh_token");
  if (!refreshToken) {
    throw new Error("Please login again");
  }

  const response = await fetch(`${API_BASE_URL}/refresh`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ refresh_token: refreshToken }),
  });

  if (!response.ok) {
    clearAuthStorage();
    throw new Error(await getErrorMessage(response, "Please login again"));
  }

  const json = await response.json();
  const data = (json.data ?? json) as RefreshPayload;
  const accessToken = data.access_token ?? data.token;

  if (!accessToken || !data.refresh_token) {
    clearAuthStorage();
    throw new Error("Invalid refresh token response");
  }

  localStorage.setItem("shophub_token", accessToken);
  localStorage.setItem("shophub_refresh_token", data.refresh_token);
  if (data.user) {
    localStorage.setItem("shophub_user", JSON.stringify(data.user));
  }
};

export const apiFetch = async <T>(path: string, options: ApiFetchOptions = {}) => {
  const { auth = false, retryOnUnauthorized = true, headers, ...requestOptions } = options;
  const makeRequest = () =>
    fetch(`${API_BASE_URL}${path}`, {
      ...requestOptions,
      headers: {
        ...(auth ? getAuthHeaders() : {}),
        ...headers,
      },
    });

  let response = await makeRequest();

  if (auth && response.status === 401 && retryOnUnauthorized) {
    await refreshAccessToken();
    response = await makeRequest();
  }

  if (!response.ok) {
    throw new Error(await getErrorMessage(response, "Request failed"));
  }

  return (await response.json()) as T;
};
