import { Outlet, Link, useLocation } from "react-router";
import { User, ShoppingBag, Heart, MapPin, CreditCard, LogOut } from "lucide-react";
import { cn } from "../lib/utils";
import { useLogout } from "../components/Logout";
import { Button } from "../components/ui/button";
import { Card, CardContent } from "../components/ui/card";

const navItems = [
  { path: "/dashboard", label: "Profile", icon: User },
  { path: "/dashboard/orders", label: "Orders", icon: ShoppingBag },
  { path: "/dashboard/wishlist", label: "Wishlist", icon: Heart },
  { path: "/dashboard/addresses", label: "Addresses", icon: MapPin },
  { path: "/dashboard/payment", label: "Payment Methods", icon: CreditCard },
];

export function CustomerDashboardLayout() {
  const location = useLocation();
  const logout = useLogout();

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="mb-6">My Account</h1>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <aside className="md:col-span-1">
          <Card className="sticky top-20">
            <CardContent className="p-2">
          <nav className="space-y-1">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = location.pathname === item.path;

              return (
                <Button
                  key={item.path}
                  asChild
                  variant={isActive ? "default" : "ghost"}
                  className={cn("w-full justify-start", !isActive && "text-muted-foreground")}
                >
                  <Link
                  to={item.path}
                  className={cn(
                    "flex items-center gap-3",
                  )}
                >
                  <Icon className="size-4" />
                  <span className="text-sm">{item.label}</span>
                </Link>
                </Button>
              );
            })}

            <Button
              variant="ghost"
              className="w-full justify-start text-muted-foreground"
              onClick={logout}
            >
              <LogOut className="size-4" />
              <span className="text-sm">Logout</span>
            </Button>
          </nav>
            </CardContent>
          </Card>
        </aside>

        <div className="md:col-span-3">
          <Outlet />
        </div>
      </div>
    </div>
  );
}
