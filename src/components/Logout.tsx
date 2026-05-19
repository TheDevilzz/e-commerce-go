import { useNavigate } from "react-router";
import { logout } from "../api/auth";

const clearAuthStorage = () => {
  localStorage.removeItem("shophub_user");
  localStorage.removeItem("shophub_token");
  localStorage.removeItem("shophub_refresh_token");
};

export function useLogout() {
  const navigate = useNavigate();

  return async () => {
    const refreshToken = localStorage.getItem("shophub_refresh_token");
    if (refreshToken) {
      try {
        await logout(refreshToken);
      } catch {
        // Local logout should still complete even if the token was already revoked.
      }
    }
    clearAuthStorage();
    navigate("/login", { replace: true });
  };
}

export function Logout() {
  const refreshToken = localStorage.getItem("shophub_refresh_token");
  if (refreshToken) {
    logout(refreshToken).finally(() => {
      clearAuthStorage();
      window.location.replace("/login");
    });
    return null;
  }
  clearAuthStorage();
  window.location.replace("/login");
  return null;
}
