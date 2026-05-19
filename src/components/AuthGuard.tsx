import { Navigate, Outlet } from "react-router";

type AuthGuardProps = {
  requiredRole?: "admin" | "user";
};

const getHomePathByRole = (role?: string) => {
  if (role === "admin") {
    return "/admin";
  }

  return "/dashboard";
};

export function AuthGuard({ requiredRole }: AuthGuardProps) {
  const rawUser = localStorage.getItem("shophub_user");
  let userRole: string | undefined;

  if (!rawUser) {
    return <Navigate to="/login" replace />;
  }

  try {
    const user = JSON.parse(rawUser);
    userRole = user.role;
  } catch {
    localStorage.removeItem("shophub_user");
    localStorage.removeItem("shophub_token");
    localStorage.removeItem("shophub_refresh_token");
    return <Navigate to="/login" replace />;
  }

  if (requiredRole && userRole !== requiredRole) {
    return <Navigate to={getHomePathByRole(userRole)} replace />;
  }

  return <Outlet />;
}
