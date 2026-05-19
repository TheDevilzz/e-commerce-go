import { Link } from "react-router";
import { ShoppingCart, Search, User, Menu } from "lucide-react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { useEffect, useState } from "react";
import { Badge } from "./ui/badge";
import { ThemeToggle } from "./ThemeToggle";
import { getCart } from "../api/cart";

export function Header() {
  const [cartCount, setCartCount] = useState(0);
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  useEffect(() => {
    const loadCartCount = async () => {
      if (!localStorage.getItem("shophub_token")) {
        setCartCount(0);
        return;
      }

      try {
        const response = await getCart();
        setCartCount(
          response.data.reduce((total, item) => total + item.quantity, 0),
        );
      } catch {
        setCartCount(0);
      }
    };

    void loadCartCount();
    window.addEventListener("focus", loadCartCount);
    window.addEventListener("shophub-cart-updated", loadCartCount);

    return () => {
      window.removeEventListener("focus", loadCartCount);
      window.removeEventListener("shophub-cart-updated", loadCartCount);
    };
  }, []);

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center gap-6">
            <Link to="/" className="flex items-center space-x-2">
              <div className="size-8 rounded-lg bg-primary" />
              <span className="font-semibold text-lg">ShopHub</span>
            </Link>

            <nav className="hidden md:flex items-center gap-6">
              <Link
                to="/category"
                className="text-sm font-medium text-foreground/60 transition-colors hover:text-foreground"
              >
                Categories
              </Link>
              <Link
                to="/promo"
                className="text-sm font-medium text-foreground/60 transition-colors hover:text-foreground"
              >
                Deals
              </Link>
            </nav>
          </div>

          <div className="hidden md:flex flex-1 max-w-md mx-6">
            <div className="relative w-full">
              <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search products..."
                className="pl-9 w-full"
              />
            </div>
          </div>

          <div className="flex items-center gap-3">
            <ThemeToggle />

            <Link to="/cart" className="relative">
              <Button variant="ghost" size="icon">
                <ShoppingCart className="size-5" />
              </Button>
              {cartCount > 0 && (
                <Badge className="absolute -top-1 -right-1 size-5 flex items-center justify-center p-0 text-xs">
                  {cartCount}
                </Badge>
              )}
            </Link>

            <Link to="/login">
              <Button variant="ghost" size="icon">
                <User className="size-5" />
              </Button>
            </Link>

            <Button
              variant="ghost"
              size="icon"
              className="md:hidden"
              onClick={() => setIsMenuOpen(!isMenuOpen)}
            >
              <Menu className="size-5" />
            </Button>
          </div>
        </div>

        {isMenuOpen && (
          <div className="md:hidden py-4 border-t">
            <div className="flex flex-col gap-4">
              <Input
                placeholder="Search products..."
                className="w-full"
              />
              <nav className="flex flex-col gap-2">
                <Link
                  to="/category"
                  className="text-sm font-medium px-2 py-1 rounded hover:bg-accent"
                  onClick={() => setIsMenuOpen(false)}
                >
                  Categories
                </Link>
                <Link
                  to="/promo"
                  className="text-sm font-medium px-2 py-1 rounded hover:bg-accent"
                  onClick={() => setIsMenuOpen(false)}
                >
                  Deals
                </Link>
              </nav>
            </div>
          </div>
        )}
      </div>
    </header>
  );
}
