import { createBrowserRouter } from "react-router";
import { RootLayout } from "./layouts/RootLayout";
import { HomePage } from "./pages/Home";
import { LoginPage } from "./pages/Login";
import { RegisterPage } from "./pages/Register";
import { CategoryPage } from "./pages/Category";
import { PromoPage } from "./pages/Promo";
import { CartPage } from "./pages/Cart";
import { ProductDetailPage } from "./pages/ProductDetail";
import { CustomerDashboardLayout } from "./layouts/CustomerDashboardLayout";
import { CustomerProfilePage } from "./pages/customer/Profile";
import { CustomerOrdersPage } from "./pages/customer/Orders";
import { CustomerWishlistPage } from "./pages/customer/Wishlist";
import { CustomerAddressesPage } from "./pages/customer/Addresses";
import { CustomerPaymentPage } from "./pages/customer/Payment";
import { AdminDashboardLayout } from "./layouts/AdminDashboardLayout";
import { AdminOverviewPage } from "./pages/admin/Overview";
import { AdminProductsPage } from "./pages/admin/Products";
import { AdminCategorysPage } from "./pages/admin/Categorys";
import { AdminOrdersPage } from "./pages/admin/Orders";
import { AdminCustomersPage } from "./pages/admin/Customers";
import { AdminPromotionsPage } from "./pages/admin/Promotions";
import { AdminInventoryPage } from "./pages/admin/Inventory";
import { NotFoundPage } from "./pages/NotFound";
import { AuthGuard } from "./components/AuthGuard";

export const router = createBrowserRouter([
  {
    path: "/",
    Component: RootLayout,
    children: [
      { index: true, Component: HomePage },
      { path: "login", Component: LoginPage },
      { path: "register", Component: RegisterPage },
      { path: "category/:categoryName?", Component: CategoryPage },
      { path: "promo", Component: PromoPage },
      { path: "cart", Component: CartPage },
      { path: "product/:productId", Component: ProductDetailPage },
      {
        element: <AuthGuard requiredRole="user" />,
        children: [
          {
            path: "dashboard",
            Component: CustomerDashboardLayout,
            children: [
              { index: true, Component: CustomerProfilePage },
              { path: "orders", Component: CustomerOrdersPage },
              { path: "wishlist", Component: CustomerWishlistPage },
              { path: "addresses", Component: CustomerAddressesPage },
              { path: "payment", Component: CustomerPaymentPage },
            ],
          },
        ],
      },
      {
        element: <AuthGuard requiredRole="admin" />,
        children: [
          {
            path: "admin",
            Component: AdminDashboardLayout,
            children: [
              { index: true, Component: AdminOverviewPage },
              { path: "products", Component: AdminProductsPage },
              { path: "categories", Component: AdminCategorysPage },
              { path: "orders", Component: AdminOrdersPage },
              { path: "customers", Component: AdminCustomersPage },
              { path: "promotions", Component: AdminPromotionsPage },
              { path: "inventory", Component: AdminInventoryPage },
            ],
          },
        ],
      },
      { path: "*", Component: NotFoundPage },
    ],
  },
]);
