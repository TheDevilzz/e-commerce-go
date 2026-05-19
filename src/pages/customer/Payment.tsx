import { useEffect, useState } from "react";
import { CreditCard, Plus, Trash2 } from "lucide-react";
import { toast } from "sonner";

import { createPaymentMethod, deletePaymentMethod, getPaymentMethods, type PaymentMethod } from "../../api/customer";
import { Button } from "../../components/ui/button";
import { Card, CardContent } from "../../components/ui/card";
import { Badge } from "../../components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "../../components/ui/dialog";
import { Input } from "../../components/ui/input";
import { Label } from "../../components/ui/label";

const emptyMethod = { type: "visa", last4: "", expiry: "", is_default: false };

export function CustomerPaymentPage() {
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [formData, setFormData] = useState(emptyMethod);

  const loadPaymentMethods = async () => {
    try {
      const response = await getPaymentMethods();
      setPaymentMethods(response.data);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to load payment methods");
    }
  };

  useEffect(() => {
    void Promise.resolve().then(loadPaymentMethods);
  }, []);

  const addMethod = async (event: React.FormEvent) => {
    event.preventDefault();
    try {
      await createPaymentMethod({ ...formData, is_default: formData.is_default || paymentMethods.length === 0 });
      toast.success("Payment method added");
      setFormData(emptyMethod);
      setIsDialogOpen(false);
      await loadPaymentMethods();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to add payment method");
    }
  };

  const removeMethod = async (id: number) => {
    try {
      await deletePaymentMethod(id);
      toast.success("Payment method deleted");
      await loadPaymentMethods();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to delete payment method");
    }
  };

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h2>Payment Methods</h2>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="mr-2 size-4" />
              Add Card
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Add Payment Method</DialogTitle>
            </DialogHeader>
            <form onSubmit={addMethod} className="space-y-3">
              <div className="space-y-1">
                <Label htmlFor="type">Type</Label>
                <Input id="type" value={formData.type} onChange={(event) => setFormData({ ...formData, type: event.target.value })} required />
              </div>
              <div className="space-y-1">
                <Label htmlFor="last4">Last 4 Digits</Label>
                <Input id="last4" maxLength={4} value={formData.last4} onChange={(event) => setFormData({ ...formData, last4: event.target.value })} required />
              </div>
              <div className="space-y-1">
                <Label htmlFor="expiry">Expiry</Label>
                <Input id="expiry" placeholder="MM/YY" value={formData.expiry} onChange={(event) => setFormData({ ...formData, expiry: event.target.value })} required />
              </div>
              <Button type="submit" className="w-full">Save Card</Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        {paymentMethods.map((method) => (
          <Card key={method.id}>
            <CardContent className="p-6">
              <div className="mb-3 flex items-start justify-between">
                <div className="flex items-center gap-3">
                  <CreditCard className="size-6 text-primary" />
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="font-medium capitalize">{method.type}</span>
                      {method.is_default && <Badge variant="outline" className="text-xs">Default</Badge>}
                    </div>
                    <p className="text-sm text-muted-foreground">•••• {method.last4}</p>
                  </div>
                </div>
                <Button variant="ghost" size="icon" className="size-8" onClick={() => removeMethod(method.id)}>
                  <Trash2 className="size-3 text-destructive" />
                </Button>
              </div>
              <p className="text-sm text-muted-foreground">Expires {method.expiry}</p>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
