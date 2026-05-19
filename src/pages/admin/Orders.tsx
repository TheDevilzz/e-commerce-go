import { useEffect, useState } from "react";
import { toast } from "sonner";
import { getAdminOrders, type ApiOrder } from "../../api/order";
import { formatPrice, formatDate } from "../../lib/utils";
import { Badge } from "../../components/ui/badge";
import { Button } from "../../components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "../../components/ui/dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../components/ui/table";

export function AdminOrdersPage() {
  const [orders, setOrders] = useState<ApiOrder[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedOrder, setSelectedOrder] = useState<ApiOrder | null>(null);

  useEffect(() => {
    let isMounted = true;

    const loadOrders = async () => {
      try {
        const response = await getAdminOrders();
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

    void Promise.resolve().then(loadOrders);

    return () => {
      isMounted = false;
    };
  }, []);

  const getStatusVariant = (
    status: string
  ): "default" | "success" | "warning" | "destructive" => {
    switch (status) {
      case "delivered":
        return "success";
      case "shipped":
        return "default";
      case "processing":
        return "warning";
      case "cancelled":
        return "destructive";
      default:
        return "default";
    }
  };

  return (
    <div className="space-y-6">
      <h2>Orders</h2>

      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Order ID</TableHead>
              <TableHead>Date</TableHead>
              <TableHead>Customer</TableHead>
              <TableHead>Items</TableHead>
              <TableHead>Total</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {orders.map((order) => (
              <TableRow key={order.id}>
                <TableCell className="font-medium">#{order.id}</TableCell>
                <TableCell>{formatDate(order.created_at ?? order.CreatedAt ?? "")}</TableCell>
                <TableCell>User #{order.user_id}</TableCell>
                <TableCell>{order.items.length} item{order.items.length !== 1 ? "s" : ""}</TableCell>
                <TableCell className="font-medium">
                  {formatPrice(order.total_price)}
                </TableCell>
                <TableCell>
                  <Badge variant={getStatusVariant(order.status)}>
                    {order.status.charAt(0).toUpperCase() +
                      order.status.slice(1)}
                  </Badge>
                </TableCell>
                <TableCell className="text-right">
                  <Button variant="ghost" size="sm" onClick={() => setSelectedOrder(order)}>
                    View Details
                  </Button>
                </TableCell>
              </TableRow>
            ))}
            {!isLoading && orders.length === 0 && (
              <TableRow>
                <TableCell colSpan={7} className="py-8 text-center text-muted-foreground">
                  No orders found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <Dialog open={Boolean(selectedOrder)} onOpenChange={(open) => !open && setSelectedOrder(null)}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Order #{selectedOrder?.id}</DialogTitle>
            <DialogDescription>
              Customer User #{selectedOrder?.user_id} · {selectedOrder?.status}
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-3">
            {selectedOrder?.items.map((item) => (
              <div key={item.id} className="flex items-center justify-between rounded-md border p-3 text-sm">
                <div>
                  <div className="font-medium">{item.product_name || `Product #${item.product_id}`}</div>
                  <div className="text-muted-foreground">Quantity: {item.quantity}</div>
                </div>
                <div className="font-medium">{formatPrice(item.price * item.quantity)}</div>
              </div>
            ))}
            <div className="flex justify-between border-t pt-3 font-semibold">
              <span>Total</span>
              <span>{formatPrice(selectedOrder?.total_price ?? 0)}</span>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
