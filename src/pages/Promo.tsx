import { useState, useEffect } from "react";
import { ProductCard } from "../components/ProductCard";
import { fetchProducts, type ProductApiItem } from "../api/product";
import { toProductCardItem } from "../api/mappers";
import { applyPromotion, getPromotions } from "../api/promotion";
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { Badge } from "../components/ui/badge";
import { Clock } from "lucide-react";
import { toast } from "sonner";

export function PromoPage() {
  const [promoProducts, setPromoProducts] = useState<ProductApiItem[]>([]);
  const [promoCode, setPromoCode] = useState("");
  const [timeLeft, setTimeLeft] = useState({
    hours: 23,
    minutes: 45,
    seconds: 30,
  });

  useEffect(() => {
    const loadPromos = async () => {
      try {
        await getPromotions(true);
        const response = await fetchProducts();
        setPromoProducts(response.data.slice(0, 8));
      } catch (error) {
        toast.error(error instanceof Error ? error.message : "Failed to load promotions");
      }
    };
    void loadPromos();

    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev.seconds > 0) {
          return { ...prev, seconds: prev.seconds - 1 };
        } else if (prev.minutes > 0) {
          return { ...prev, minutes: prev.minutes - 1, seconds: 59 };
        } else if (prev.hours > 0) {
          return { hours: prev.hours - 1, minutes: 59, seconds: 59 };
        }
        return prev;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, []);

  const handleApplyPromo = async () => {
    try {
      await applyPromotion(promoCode, 100);
      toast.success(`Promo code "${promoCode}" applied!`);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Invalid promo code");
    }
  };

  return (
    <div>
      <section className="bg-gradient-to-r from-destructive to-destructive/80 text-destructive-foreground py-16">
        <div className="container mx-auto px-4 text-center">
          <Badge variant="outline" className="mb-4 bg-background/20 border-destructive-foreground/30">
            Limited Time Offer
          </Badge>
          <h1 className="text-4xl md:text-5xl mb-4">Flash Sale - Up to 40% Off!</h1>
          <p className="text-lg mb-8 opacity-90">
            Don't miss out on our biggest deals of the season
          </p>

          <div className="flex items-center justify-center gap-4 mb-8">
            <Clock className="size-6" />
            <div className="flex gap-2">
              <div className="bg-background/20 backdrop-blur-sm rounded-lg px-4 py-2">
                <div className="text-2xl font-bold">{String(timeLeft.hours).padStart(2, "0")}</div>
                <div className="text-xs opacity-75">Hours</div>
              </div>
              <div className="text-2xl font-bold self-center">:</div>
              <div className="bg-background/20 backdrop-blur-sm rounded-lg px-4 py-2">
                <div className="text-2xl font-bold">{String(timeLeft.minutes).padStart(2, "0")}</div>
                <div className="text-xs opacity-75">Minutes</div>
              </div>
              <div className="text-2xl font-bold self-center">:</div>
              <div className="bg-background/20 backdrop-blur-sm rounded-lg px-4 py-2">
                <div className="text-2xl font-bold">{String(timeLeft.seconds).padStart(2, "0")}</div>
                <div className="text-xs opacity-75">Seconds</div>
              </div>
            </div>
          </div>

          <Button size="lg" variant="secondary">
            Shop Now
          </Button>
        </div>
      </section>

      <section className="container mx-auto px-4 py-12">
        <div className="max-w-md mx-auto mb-12">
          <h2 className="text-center mb-4">Have a Promo Code?</h2>
          <div className="flex gap-2">
            <Input
              placeholder="Enter promo code"
              value={promoCode}
              onChange={(e) => setPromoCode(e.target.value)}
            />
            <Button onClick={handleApplyPromo}>Apply</Button>
          </div>
        </div>

        <h2 className="mb-6">Featured Deals</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {promoProducts.map((product) => (
            <ProductCard key={product.id} product={toProductCardItem(product)} />
          ))}
        </div>
      </section>
    </div>
  );
}
