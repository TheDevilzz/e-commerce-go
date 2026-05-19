import { useEffect, useState } from "react";
import { Plus, Trash2 } from "lucide-react";
import { toast } from "sonner";

import { createPromotion, deletePromotion, getPromotions, type Promotion } from "../../api/promotion";
import { Button } from "../../components/ui/button";
import { Badge } from "../../components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "../../components/ui/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "../../components/ui/dialog";
import { Input } from "../../components/ui/input";
import { Label } from "../../components/ui/label";

const emptyPromotion = {
  code: "",
  description: "",
  discount: 10,
  type: "percentage" as const,
  start_date: "2026-05-01",
  end_date: "2026-12-31",
  status: "active" as const,
};

export function AdminPromotionsPage() {
  const [promotions, setPromotions] = useState<Promotion[]>([]);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [formData, setFormData] = useState(emptyPromotion);

  const loadPromotions = async () => {
    try {
      const response = await getPromotions();
      setPromotions(response.data);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to load promotions");
    }
  };

  useEffect(() => {
    void Promise.resolve().then(loadPromotions);
  }, []);

  const addPromotion = async (event: React.FormEvent) => {
    event.preventDefault();
    try {
      await createPromotion(formData);
      toast.success("Promotion created");
      setFormData(emptyPromotion);
      setIsDialogOpen(false);
      await loadPromotions();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to create promotion");
    }
  };

  const removePromotion = async (id: number) => {
    try {
      await deletePromotion(id);
      toast.success("Promotion deleted");
      await loadPromotions();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to delete promotion");
    }
  };

  const getStatusVariant = (status: string): "default" | "success" | "warning" => {
    if (status === "active") return "success";
    if (status === "scheduled") return "warning";
    return "default";
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2>Promotions</h2>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="mr-2 size-4" />
              Create Promotion
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create Promotion</DialogTitle>
            </DialogHeader>
            <form onSubmit={addPromotion} className="space-y-3">
              <div className="space-y-1">
                <Label htmlFor="code">Code</Label>
                <Input id="code" value={formData.code} onChange={(event) => setFormData({ ...formData, code: event.target.value })} required />
              </div>
              <div className="space-y-1">
                <Label htmlFor="description">Description</Label>
                <Input id="description" value={formData.description} onChange={(event) => setFormData({ ...formData, description: event.target.value })} required />
              </div>
              <div className="space-y-1">
                <Label htmlFor="discount">Discount</Label>
                <Input id="discount" type="number" min={1} value={formData.discount} onChange={(event) => setFormData({ ...formData, discount: Number(event.target.value) })} required />
              </div>
              <Button type="submit" className="w-full">Save Promotion</Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        {promotions.map((promo) => (
          <Card key={promo.id}>
            <CardHeader>
              <div className="flex items-start justify-between">
                <div>
                  <CardTitle className="text-lg">{promo.code}</CardTitle>
                  <p className="mt-1 text-sm text-muted-foreground">{promo.description}</p>
                </div>
                <Badge variant={getStatusVariant(promo.status)}>{promo.status}</Badge>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <div>
                  <div className="text-2xl font-bold text-primary">
                    {promo.type === "percentage" ? `${promo.discount}%` : `$${promo.discount}`}
                  </div>
                  <div className="text-sm text-muted-foreground">{promo.type === "percentage" ? "Discount" : "Off"}</div>
                </div>
                <div className="text-sm">
                  <div className="text-muted-foreground">Valid Period</div>
                  <div>{promo.start_date} - {promo.end_date}</div>
                </div>
                <div className="text-sm">
                  <div className="text-muted-foreground">Usage</div>
                  <div className="font-medium">{promo.usage_count} times</div>
                </div>
                <Button variant="outline" size="sm" className="w-full" onClick={() => removePromotion(promo.id)}>
                  <Trash2 className="mr-1 size-3 text-destructive" />
                  Delete
                </Button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
