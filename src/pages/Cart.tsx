import { type FormEvent, useEffect, useMemo, useState } from "react";
import { Link } from "react-router";
import { Loader2, MapPin, Minus, Plus, Tag, Trash2 } from "lucide-react";
import { toast } from "sonner";

import { applyPromotion } from "../api/promotion";
import { deleteCartItem, getCart, notifyCartUpdated, updateCartItem, type CartItem } from "../api/cart";
import { createAddress, getAddresses, type Address } from "../api/customer";
import { createOrder } from "../api/order";
import {
  findThaiAddressRow,
  getThaiAddressData,
  getThaiDistricts,
  getThaiProvinces,
  getThaiSubdistricts,
  type ThaiAddressData,
} from "../api/thaiAddress";
import { toProductCardItem } from "../api/mappers";
import { formatPrice } from "../lib/utils";
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card";
import { Separator } from "../components/ui/separator";
import { Badge } from "../components/ui/badge";
import { Checkbox } from "../components/ui/checkbox";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "../components/ui/dialog";
import { Label } from "../components/ui/label";
import { RadioGroup, RadioGroupItem } from "../components/ui/radio-group";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../components/ui/select";

const emptyAddressForm = {
  type: "home",
  name: "",
  street: "",
  provinceCode: "",
  districtCode: "",
  subdistrictCode: "",
  is_default: false,
};

export function CartPage() {
  const [cartItems, setCartItems] = useState<CartItem[]>([]);
  const [addresses, setAddresses] = useState<Address[]>([]);
  const [selectedAddressId, setSelectedAddressId] = useState("");
  const [thaiAddressData, setThaiAddressData] = useState<ThaiAddressData | null>(null);
  const [addressForm, setAddressForm] = useState(emptyAddressForm);
  const [provinceSearch, setProvinceSearch] = useState("");
  const [districtSearch, setDistrictSearch] = useState("");
  const [subdistrictSearch, setSubdistrictSearch] = useState("");
  const [isProvinceSuggestOpen, setIsProvinceSuggestOpen] = useState(false);
  const [isDistrictSuggestOpen, setIsDistrictSuggestOpen] = useState(false);
  const [isSubdistrictSuggestOpen, setIsSubdistrictSuggestOpen] = useState(false);
  const [isAddressDialogOpen, setIsAddressDialogOpen] = useState(false);
  const [isAddressLoading, setIsAddressLoading] = useState(true);
  const [isThailandApiLoading, setIsThailandApiLoading] = useState(true);
  const [promoCode, setPromoCode] = useState("");
  const [discount, setDiscount] = useState(0);
  const [isLoading, setIsLoading] = useState(true);

  const loadCart = async () => {
    try {
      const response = await getCart();
      setCartItems(response.data);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Please login to view cart");
    } finally {
      setIsLoading(false);
    }
  };

  const loadAddresses = async () => {
    try {
      const response = await getAddresses();
      setAddresses(response.data);
      setSelectedAddressId((current) => {
        if (current && response.data.some((address) => String(address.id) === current)) {
          return current;
        }

        return String(response.data.find((address) => address.is_default)?.id ?? response.data[0]?.id ?? "");
      });
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to load addresses");
    } finally {
      setIsAddressLoading(false);
    }
  };

  useEffect(() => {
    void Promise.resolve().then(async () => {
      await Promise.all([loadCart(), loadAddresses()]);
    });
  }, []);

  useEffect(() => {
    void Promise.resolve().then(async () => {
      try {
        setThaiAddressData(await getThaiAddressData());
      } catch (error) {
        toast.error(error instanceof Error ? error.message : "Thailand address API is unavailable");
      } finally {
        setIsThailandApiLoading(false);
      }
    });
  }, []);

  const subtotal = cartItems.reduce((sum, item) => sum + Number(item.product.price ?? 0) * item.quantity, 0);
  const shipping = subtotal > 50 ? 0 : 9.99;
  const total = subtotal - discount + shipping;
  const provinces = useMemo(() => getThaiProvinces(thaiAddressData), [thaiAddressData]);
  const filteredProvinces = useMemo(() => {
    const search = provinceSearch.trim().toLocaleLowerCase("th");

    if (!search) {
      return provinces.slice(0, 8);
    }

    return provinces
      .filter((province) => province.name.toLocaleLowerCase("th").startsWith(search))
      .slice(0, 8);
  }, [provinceSearch, provinces]);
  const districts = useMemo(
    () => getThaiDistricts(thaiAddressData, Number(addressForm.provinceCode)),
    [addressForm.provinceCode, thaiAddressData],
  );
  const filteredDistricts = useMemo(() => {
    const search = districtSearch.trim().toLocaleLowerCase("th");

    if (!search) {
      return districts.slice(0, 8);
    }

    return districts
      .filter((district) => district.name.toLocaleLowerCase("th").startsWith(search))
      .slice(0, 8);
  }, [districtSearch, districts]);
  const subdistricts = useMemo(
    () => getThaiSubdistricts(thaiAddressData, Number(addressForm.districtCode)),
    [addressForm.districtCode, thaiAddressData],
  );
  const filteredSubdistricts = useMemo(() => {
    const search = subdistrictSearch.trim().toLocaleLowerCase("th");

    if (!search) {
      return subdistricts.slice(0, 8);
    }

    return subdistricts
      .filter((subdistrict) => subdistrict.name.toLocaleLowerCase("th").startsWith(search))
      .slice(0, 8);
  }, [subdistrictSearch, subdistricts]);
  const selectedThaiAddress = useMemo(
    () => findThaiAddressRow(thaiAddressData, Number(addressForm.subdistrictCode)),
    [addressForm.subdistrictCode, thaiAddressData],
  );

  const updateQuantity = async (item: CartItem, change: number) => {
    const nextQuantity = item.quantity + change;
    if (nextQuantity <= 0) {
      return;
    }
    try {
      await updateCartItem(item.id, item.product_id, nextQuantity);
      await loadCart();
      notifyCartUpdated();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to update cart");
    }
  };

  const removeItem = async (id: number) => {
    try {
      await deleteCartItem(id);
      toast.success("Removed from cart");
      await loadCart();
      notifyCartUpdated();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to remove item");
    }
  };

  const handleApplyPromo = async () => {
    try {
      const response = await applyPromotion(promoCode, subtotal);
      setDiscount(response.data.discount);
      toast.success("Promo code applied");
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Invalid promo code");
    }
  };

  const addAddress = async (event: FormEvent) => {
    event.preventDefault();

    if (!selectedThaiAddress) {
      toast.error("Please select province, district, and subdistrict");
      return;
    }

    try {
      const street = `${addressForm.street}, ${selectedThaiAddress.subdistrictNameTh}`;
      const response = await createAddress({
        type: addressForm.type,
        name: addressForm.name,
        street,
        city: selectedThaiAddress.districtNameTh,
        state: selectedThaiAddress.provinceNameTh,
        zip: String(selectedThaiAddress.postalCode),
        is_default: addressForm.is_default || addresses.length === 0,
      });

      toast.success("Address added");
      setAddressForm(emptyAddressForm);
      setProvinceSearch("");
      setDistrictSearch("");
      setSubdistrictSearch("");
      setIsProvinceSuggestOpen(false);
      setIsDistrictSuggestOpen(false);
      setIsSubdistrictSuggestOpen(false);
      setIsAddressDialogOpen(false);
      await loadAddresses();
      setSelectedAddressId(String(response.data.id));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to add address");
    }
  };

  const handleCheckout = async () => {
    const addressId = Number(selectedAddressId);

    if (!addressId) {
      toast.error("Please select a shipping address");
      return;
    }

    try {
      await createOrder(
        cartItems.map((item) => ({ product_id: item.product_id, quantity: item.quantity })),
        addressId,
      );
      toast.success("Order created");
      await loadCart();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to create order");
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="mb-8">Shopping Cart</h1>

      {!isLoading && cartItems.length === 0 ? (
        <div className="py-12 text-center">
          <p className="mb-4 text-muted-foreground">Your cart is empty</p>
          <Link to="/category">
            <Button>Continue Shopping</Button>
          </Link>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
          <div className="space-y-4 lg:col-span-2">
            {cartItems.map((item) => {
              const product = toProductCardItem(item.product);
              return (
                <Card key={item.id}>
                  <CardContent className="p-4">
                    <div className="flex gap-4">
                      <img src={product.image} alt={product.name} className="size-24 rounded-lg object-cover" />
                      <div className="flex-1">
                        <div className="mb-2 flex items-start justify-between">
                          <div>
                            <h3 className="font-medium">{product.name}</h3>
                            <p className="text-sm text-muted-foreground">{product.category}</p>
                          </div>
                          <Button variant="ghost" size="icon" onClick={() => removeItem(item.id)}>
                            <Trash2 className="size-4 text-destructive" />
                          </Button>
                        </div>

                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-2 rounded-lg border">
                            <Button variant="ghost" size="icon" className="size-8" onClick={() => updateQuantity(item, -1)}>
                              <Minus className="size-3" />
                            </Button>
                            <span className="w-8 text-center">{item.quantity}</span>
                            <Button variant="ghost" size="icon" className="size-8" onClick={() => updateQuantity(item, 1)}>
                              <Plus className="size-3" />
                            </Button>
                          </div>
                          <div className="font-semibold">{formatPrice(product.price * item.quantity)}</div>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              );
            })}

            <Card>
              <CardHeader className="flex-row items-center justify-between space-y-0">
                <CardTitle>Shipping Address</CardTitle>
                <Dialog open={isAddressDialogOpen} onOpenChange={setIsAddressDialogOpen}>
                  <DialogTrigger asChild>
                    <Button variant="outline" size="sm">
                      <Plus className="mr-2 size-4" />
                      Add Address
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="sm:max-w-[560px]">
                    <DialogHeader>
                      <DialogTitle>Add Thailand Address</DialogTitle>
                    </DialogHeader>
                    <form onSubmit={addAddress} className="space-y-4 text-left">
                      <div className="grid gap-3 sm:grid-cols-2">
                        <div className="space-y-1">
                          <Label htmlFor="address-type">Type</Label>
                          <Select
                            value={addressForm.type}
                            onValueChange={(value) => setAddressForm({ ...addressForm, type: value })}
                          >
                            <SelectTrigger id="address-type">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="home">Home</SelectItem>
                              <SelectItem value="work">Work</SelectItem>
                            </SelectContent>
                          </Select>
                        </div>
                        <div className="space-y-1">
                          <Label htmlFor="address-name">Recipient name</Label>
                          <Input
                            id="address-name"
                            value={addressForm.name}
                            onChange={(event) => setAddressForm({ ...addressForm, name: event.target.value })}
                            required
                          />
                        </div>
                      </div>

                      <div className="space-y-1">
                        <Label htmlFor="address-street">Street, building, house no.</Label>
                        <Input
                          id="address-street"
                          value={addressForm.street}
                          onChange={(event) => setAddressForm({ ...addressForm, street: event.target.value })}
                          required
                        />
                      </div>

                      <div className="grid gap-3 sm:grid-cols-3">
                        <div className="relative space-y-1">
                          <Label>Province</Label>
                          <Input
                            value={provinceSearch}
                            onChange={(event) => {
                              setProvinceSearch(event.target.value);
                              setDistrictSearch("");
                              setSubdistrictSearch("");
                              setIsProvinceSuggestOpen(true);
                              setIsDistrictSuggestOpen(false);
                              setIsSubdistrictSuggestOpen(false);
                              setAddressForm({
                                ...addressForm,
                                provinceCode: "",
                                districtCode: "",
                                subdistrictCode: "",
                              });
                            }}
                            onFocus={() => setIsProvinceSuggestOpen(true)}
                            placeholder="พิมพ์จังหวัด เช่น ก"
                            disabled={isThailandApiLoading || provinces.length === 0}
                            required
                          />
                          {isProvinceSuggestOpen && filteredProvinces.length > 0 && (
                            <div className="absolute left-0 right-0 top-full z-50 mt-1 max-h-60 overflow-auto rounded-md border bg-popover p-1 text-popover-foreground shadow-md">
                              {filteredProvinces.map((province) => (
                                <Button
                                  key={province.code}
                                  type="button"
                                  variant="ghost"
                                  className="h-auto w-full justify-start px-2 py-2"
                                  onMouseDown={(event) => {
                                    event.preventDefault();
                                    setAddressForm({
                                      ...addressForm,
                                      provinceCode: String(province.code),
                                      districtCode: "",
                                      subdistrictCode: "",
                                    });
                                    setProvinceSearch(province.name);
                                    setDistrictSearch("");
                                    setSubdistrictSearch("");
                                    setIsProvinceSuggestOpen(false);
                                  }}
                                >
                                  {province.name}
                                </Button>
                              ))}
                            </div>
                          )}
                        </div>
                        <div className="relative space-y-1">
                          <Label>District</Label>
                          <Input
                            value={districtSearch}
                            onChange={(event) => {
                              setDistrictSearch(event.target.value);
                              setSubdistrictSearch("");
                              setIsDistrictSuggestOpen(true);
                              setIsSubdistrictSuggestOpen(false);
                              setAddressForm({ ...addressForm, districtCode: "", subdistrictCode: "" });
                            }}
                            onFocus={() => setIsDistrictSuggestOpen(true)}
                            placeholder="Type district"
                            disabled={!addressForm.provinceCode}
                            required
                          />
                          {isDistrictSuggestOpen && addressForm.provinceCode && filteredDistricts.length > 0 && (
                            <div className="absolute left-0 right-0 top-full z-40 mt-1 max-h-60 overflow-auto rounded-md border bg-popover p-1 text-popover-foreground shadow-md">
                              {filteredDistricts.map((district) => (
                                <Button
                                  key={district.code}
                                  type="button"
                                  variant="ghost"
                                  className="h-auto w-full justify-start px-2 py-2"
                                  onMouseDown={(event) => {
                                    event.preventDefault();
                                    setAddressForm({
                                      ...addressForm,
                                      districtCode: String(district.code),
                                      subdistrictCode: "",
                                    });
                                    setDistrictSearch(district.name);
                                    setSubdistrictSearch("");
                                    setIsDistrictSuggestOpen(false);
                                  }}
                                >
                                  {district.name}
                                </Button>
                              ))}
                            </div>
                          )}
                        </div>
                        <div className="relative space-y-1">
                          <Label>Subdistrict</Label>
                          <Input
                            value={subdistrictSearch}
                            onChange={(event) => {
                              setSubdistrictSearch(event.target.value);
                              setIsSubdistrictSuggestOpen(true);
                              setAddressForm({ ...addressForm, subdistrictCode: "" });
                            }}
                            onFocus={() => setIsSubdistrictSuggestOpen(true)}
                            placeholder="Type subdistrict"
                            disabled={!addressForm.districtCode}
                            required
                          />
                          {isSubdistrictSuggestOpen && addressForm.districtCode && filteredSubdistricts.length > 0 && (
                            <div className="absolute left-0 right-0 top-full z-30 mt-1 max-h-60 overflow-auto rounded-md border bg-popover p-1 text-popover-foreground shadow-md">
                              {filteredSubdistricts.map((subdistrict) => (
                                <Button
                                  key={subdistrict.code}
                                  type="button"
                                  variant="ghost"
                                  className="h-auto w-full justify-start px-2 py-2"
                                  onMouseDown={(event) => {
                                    event.preventDefault();
                                    setAddressForm({ ...addressForm, subdistrictCode: String(subdistrict.code) });
                                    setSubdistrictSearch(subdistrict.name);
                                    setIsSubdistrictSuggestOpen(false);
                                  }}
                                >
                                  {subdistrict.name}
                                </Button>
                              ))}
                            </div>
                          )}
                        </div>
                      </div>

                      <div className="grid gap-3 sm:grid-cols-[1fr_auto]">
                        <div className="space-y-1">
                          <Label htmlFor="address-zip">Postal code</Label>
                          <Input id="address-zip" value={selectedThaiAddress?.postalCode ?? ""} readOnly />
                        </div>
                        <label className="flex items-center gap-2 self-end rounded-md border px-3 py-2 text-sm">
                          <Checkbox
                            checked={addressForm.is_default}
                            onCheckedChange={(checked) =>
                              setAddressForm({ ...addressForm, is_default: checked === true })
                            }
                          />
                          Default
                        </label>
                      </div>

                      <Button type="submit" className="w-full" disabled={isThailandApiLoading}>
                        {isThailandApiLoading && <Loader2 className="mr-2 size-4 animate-spin" />}
                        Save Address
                      </Button>
                    </form>
                  </DialogContent>
                </Dialog>
              </CardHeader>
              <CardContent>
                {isAddressLoading ? (
                  <div className="flex items-center justify-center py-8 text-sm text-muted-foreground">
                    <Loader2 className="mr-2 size-4 animate-spin" />
                    Loading addresses
                  </div>
                ) : addresses.length === 0 ? (
                  <div className="rounded-md border border-dashed p-6 text-center text-sm text-muted-foreground">
                    Add a shipping address before checkout.
                  </div>
                ) : (
                  <RadioGroup value={selectedAddressId} onValueChange={setSelectedAddressId}>
                    {addresses.map((address) => (
                      <label
                        key={address.id}
                        className="flex cursor-pointer items-start gap-3 rounded-md border p-4 text-left transition-colors hover:bg-accent"
                      >
                        <RadioGroupItem value={String(address.id)} className="mt-1" />
                        <MapPin className="mt-0.5 size-4 shrink-0 text-primary" />
                        <div className="min-w-0 flex-1 space-y-1">
                          <div className="flex flex-wrap items-center gap-2">
                            <span className="font-medium">{address.name}</span>
                            <span className="text-sm capitalize text-muted-foreground">{address.type}</span>
                            {address.is_default && <Badge variant="outline">Default</Badge>}
                          </div>
                          <p className="text-sm text-muted-foreground">{address.street}</p>
                          <p className="text-sm text-muted-foreground">
                            {address.city}, {address.state} {address.zip}
                          </p>
                        </div>
                      </label>
                    ))}
                  </RadioGroup>
                )}
              </CardContent>
            </Card>
          </div>

          <Card className="sticky top-20 h-fit">
            <CardHeader>
              <CardTitle>Order Summary</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex gap-2">
                <Input placeholder="Promo code" value={promoCode} onChange={(e) => setPromoCode(e.target.value)} />
                <Button variant="outline" onClick={handleApplyPromo}>
                  <Tag className="size-4" />
                </Button>
              </div>
              <Separator />
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Subtotal</span>
                  <span>{formatPrice(subtotal)}</span>
                </div>
                {discount > 0 && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Discount</span>
                    <span className="text-success">-{formatPrice(discount)}</span>
                  </div>
                )}
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Shipping</span>
                  <span>{shipping === 0 ? "FREE" : formatPrice(shipping)}</span>
                </div>
              </div>
              <Separator />
              <div className="flex justify-between font-semibold">
                <span>Total</span>
                <span>{formatPrice(total)}</span>
              </div>
              <Button className="w-full" size="lg" onClick={handleCheckout} disabled={cartItems.length === 0}>
                Proceed to Checkout
              </Button>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
}
