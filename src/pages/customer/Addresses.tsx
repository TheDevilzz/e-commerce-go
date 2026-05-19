import { useEffect, useState } from "react";
import { MapPin, Plus, Trash2 } from "lucide-react";
import { toast } from "sonner";

import { createAddress, deleteAddress, getAddresses, type Address } from "../../api/customer";
import { Button } from "../../components/ui/button";
import { Card, CardContent } from "../../components/ui/card";
import { Badge } from "../../components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "../../components/ui/dialog";
import { Input } from "../../components/ui/input";
import { Label } from "../../components/ui/label";

const emptyAddress = {
  type: "home",
  name: "",
  street: "",
  city: "",
  state: "",
  zip: "",
  is_default: false,
};

export function CustomerAddressesPage() {
  const [addresses, setAddresses] = useState<Address[]>([]);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [formData, setFormData] = useState(emptyAddress);

  const loadAddresses = async () => {
    try {
      const response = await getAddresses();
      setAddresses(response.data);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to load addresses");
    }
  };

  useEffect(() => {
    void Promise.resolve().then(loadAddresses);
  }, []);

  const addAddress = async (event: React.FormEvent) => {
    event.preventDefault();
    try {
      await createAddress({ ...formData, is_default: formData.is_default || addresses.length === 0 });
      toast.success("Address added");
      setFormData(emptyAddress);
      setIsDialogOpen(false);
      await loadAddresses();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to add address");
    }
  };

  const removeAddress = async (id: number) => {
    try {
      await deleteAddress(id);
      toast.success("Address deleted");
      await loadAddresses();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to delete address");
    }
  };

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h2>Saved Addresses</h2>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="mr-2 size-4" />
              Add Address
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Add Address</DialogTitle>
            </DialogHeader>
            <form onSubmit={addAddress} className="space-y-3">
              {(["name", "street", "city", "state", "zip"] as const).map((field) => (
                <div key={field} className="space-y-1">
                  <Label htmlFor={field} className="capitalize">{field}</Label>
                  <Input id={field} value={formData[field]} onChange={(event) => setFormData({ ...formData, [field]: event.target.value })} required />
                </div>
              ))}
              <Button type="submit" className="w-full">Save Address</Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        {addresses.map((address) => (
          <Card key={address.id}>
            <CardContent className="p-6">
              <div className="mb-3 flex items-start justify-between">
                <div className="flex items-center gap-2">
                  <MapPin className="size-4 text-primary" />
                  <span className="font-medium capitalize">{address.type}</span>
                  {address.is_default && <Badge variant="outline" className="text-xs">Default</Badge>}
                </div>
                <Button variant="ghost" size="icon" className="size-8" onClick={() => removeAddress(address.id)}>
                  <Trash2 className="size-3 text-destructive" />
                </Button>
              </div>
              <div className="space-y-1 text-sm">
                <p className="font-medium">{address.name}</p>
                <p className="text-muted-foreground">{address.street}</p>
                <p className="text-muted-foreground">{address.city}, {address.state} {address.zip}</p>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
