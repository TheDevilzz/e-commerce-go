import { useEffect, useState } from "react";
import { login } from "../api/auth";
import { Link, useNavigate } from "react-router";
import { Eye, EyeOff } from "lucide-react";
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { Label } from "../components/ui/label";
import { Card, CardHeader, CardTitle, CardContent } from "../components/ui/card";
import { Checkbox } from "../components/ui/checkbox";
import { toast } from "sonner";

const getHomePathByRole = (role?: string) => {
  if (role === "admin") {
    return "/admin";
  }

  return "/dashboard";
};

export function LoginPage() {
  const navigate = useNavigate();
  const [showPassword, setShowPassword] = useState(false);
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      const data = await login(identifier, identifier, password);
      if (!data.user) {
        throw new Error("Login response did not include user data");
      }
      toast.success("Login successful!");
      localStorage.setItem("shophub_user", JSON.stringify(data.user));
      const accessToken = data.access_token ?? data.token;
      if (accessToken) {
        localStorage.setItem("shophub_token", accessToken);
      }
      if (data.refresh_token) {
        localStorage.setItem("shophub_refresh_token", data.refresh_token);
      }
      navigate(getHomePathByRole(data.user.role));
    } catch {
      toast.error("Login failed. Please check your credentials.");
    } finally {
      setIsSubmitting(false);
    }
  };
  useEffect(() => {
    const rawUser = localStorage.getItem("shophub_user");

    if (!rawUser) {
      return;
    }

    try {
      const user = JSON.parse(rawUser);
      navigate(getHomePathByRole(user.role), { replace: true });
    } catch {
      localStorage.removeItem("shophub_user");
      localStorage.removeItem("shophub_token");
      localStorage.removeItem("shophub_refresh_token");
    }
  }, [navigate]);

  return (
    <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center px-4 py-12">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Login to ShopHub</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="username">Username OR Email</Label>
              <Input
                id="username"
                type="text"
                placeholder="Enter your username or email"
                value={identifier}
                onChange={(e) => setIdentifier(e.target.value)}
                required
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  placeholder="Enter your password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-1 top-1/2 size-8 -translate-y-1/2 text-muted-foreground"
                >
                  {showPassword ? (
                    <EyeOff className="size-4" />
                  ) : (
                    <Eye className="size-4" />
                  )}
                </Button>
              </div>
            </div>

            <div className="flex items-center justify-between">
              <label className="flex items-center gap-2 text-sm text-muted-foreground">
                <Checkbox />
                Remember me
              </label>
            </div>

            <Button type="submit" className="w-full" disabled={isSubmitting}>
              {isSubmitting ? "Logging in..." : "Login"}
            </Button>

            <p className="text-center text-sm text-muted-foreground">
              Don't have an account?{" "}
              <Link to="/register" className="text-primary hover:underline">
                Register
              </Link>
            </p>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
