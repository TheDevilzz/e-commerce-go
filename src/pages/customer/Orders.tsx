import { useEffect, useState } from "react";
import { Link } from "react-router";
import { toast } from "sonner";

import { getMyOrders, type ApiOrder } from "../../api/order";
import { uploadPaymentSlip } from "../../api/payment";
import { formatDate, formatPrice } from "../../lib/utils";
import { Badge } from "../../components/ui/badge";
import { Button } from "../../components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "../../components/ui/card";

export function CustomerOrdersPage() {
  const [orders, setOrders] = useState<ApiOrder[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    let isMounted = true;

    const loadOrders = async () => {
      try {
        const response = await getMyOrders();
        if (isMounted) {
          setOrders(response.data);
        }
      } catch (error) {
        toast.error(error instanceof Error ? error.message : "Failed to load orders");
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    };

    void loadOrders();

    return () => {
      isMounted = false;
    };
  }, []);

  const getStatusVariant = (
    status: string,
  ): "default" | "success" | "warning" | "destructive" => {
    switch (status) {
      case "paid":
      case "delivered":
        return "success";
      case "pending":
      case "processing":
        return "warning";
      case "cancelled":
        return "destructive";
      default:
        return "default";
    }
  };

  const handleSlipUpload = async (order: ApiOrder, file?: File) => {
    if (!file) {
      return;
    }
    try {
      await uploadPaymentSlip(order.id, order.total_price, file);
      toast.success("Payment slip uploaded");
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to upload payment slip");
    }
  };

  return (
    <div className="space-y-4">
      <h2>My Orders</h2>

      {orders.map((order) => (
        <Card key={order.id}>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Order #{order.id}</CardTitle>
                <p className="mt-1 text-sm text-muted-foreground">
                  Placed on {formatDate(order.created_at ?? order.CreatedAt ?? "")}
                </p>
              </div>
              <Badge variant={getStatusVariant(order.status)}>{order.status}</Badge>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {order.items.map((item) => (
                <div key={item.id} className="flex gap-4">
                  <div className="size-16 rounded-lg bg-muted" />
                  <div className="flex-1">
                    <Link
                      to={`/product/${item.product_id}`}
                      className="font-medium hover:underline"
                    >
                      {item.product_name || `Product #${item.product_id}`}
                    </Link>
                    <p className="text-sm text-muted-foreground">
                      Quantity: {item.quantity}
                    </p>
                  </div>
                  <div className="font-medium">
                    {formatPrice(item.price * item.quantity)}
                  </div>
                </div>
              ))}

              <div className="flex items-center justify-between border-t pt-3">
                <span className="font-medium">Total</span>
                <span className="text-lg font-bold">
                  {formatPrice(order.total_price)}
                </span>
              </div>

              <div className="flex gap-2">
                <Button variant="outline" size="sm">
                  View Details
                </Button>
                <Button variant="outline" size="sm" asChild>
                  <label className="cursor-pointer">
                    Upload Slip
                    <input
                      type="file"
                      accept="image/*"
                      className="hidden"
                      onChange={(event) => handleSlipUpload(order, event.target.files?.[0])}
                    />
                  </label>
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      ))}

      {!isLoading && orders.length === 0 && (
        <Card>
          <CardContent className="py-8 text-center text-muted-foreground">
            No orders found.
          </CardContent>
        </Card>
      )}
    </div>
  );
}
