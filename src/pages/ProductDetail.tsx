import { useEffect, useState } from "react";
import { Link, useParams } from "react-router";
import { Heart, Minus, Plus, RotateCcw, Shield, Star, Truck } from "lucide-react";
import { toast } from "sonner";

import { addCartItem, notifyCartUpdated } from "../api/cart";
import { addWishlistItem } from "../api/customer";
import { ASSET_BASE_URL } from "../api/client";
import { fetchProductById, fetchProducts, type ProductApiItem } from "../api/product";
import { toProductCardItem } from "../api/mappers";
import { formatPrice } from "../lib/utils";
import { Button } from "../components/ui/button";
import { Badge } from "../components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../components/ui/tabs";
import { ProductCard } from "../components/ProductCard";

export function ProductDetailPage() {
  const { productId } = useParams();
  const [product, setProduct] = useState<ProductApiItem | null>(null);
  const [relatedProducts, setRelatedProducts] = useState<ProductApiItem[]>([]);
  const [quantity, setQuantity] = useState(1);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const loadProduct = async () => {
      if (!productId) {
        return;
      }
      try {
        const [productResponse, productsResponse] = await Promise.all([
          fetchProductById(productId),
          fetchProducts(),
        ]);
        setProduct(productResponse.data);
        setRelatedProducts(
          productsResponse.data
            .filter((item) => String(item.id) !== productId && item.category_id === productResponse.data.category_id)
            .slice(0, 4),
        );
      } catch (error) {
        toast.error(error instanceof Error ? error.message : "Failed to load product");
      } finally {
        setIsLoading(false);
      }
    };

    void loadProduct();
  }, [productId]);

  if (!product && !isLoading) {
    return (
      <div className="container mx-auto px-4 py-12 text-center">
        <h2 className="mb-4">Product not found</h2>
        <Link to="/category">
          <Button>Back to Shopping</Button>
        </Link>
      </div>
    );
  }

  if (!product) {
    return <div className="container mx-auto px-4 py-12 text-center text-muted-foreground">Loading product...</div>;
  }

  const productImage = product.image ? `${ASSET_BASE_URL}${product.image}` : "https://placehold.co/800x800?text=Product";
  const inStock = Number(product.stock ?? 0) > 0;

  const handleAddToCart = async () => {
    try {
      await addCartItem(Number(product.id), quantity);
      notifyCartUpdated();
      toast.success("Added to cart!");
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Please login before adding to cart");
    }
  };

  const handleAddWishlist = async () => {
    try {
      await addWishlistItem(Number(product.id));
      toast.success("Added to wishlist!");
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Please login before adding to wishlist");
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-16 grid grid-cols-1 gap-12 lg:grid-cols-2">
        <div className="relative aspect-square overflow-hidden rounded-lg bg-muted">
          <img src={productImage} alt={product.name} className="size-full object-cover" />
          <Button
            variant="ghost"
            size="icon"
            className="absolute right-4 top-4 bg-background/80 hover:bg-background"
            onClick={handleAddWishlist}
          >
            <Heart className="size-5" />
          </Button>
        </div>

        <div>
          <div className="mb-2">
            <Badge variant="outline">{product.category?.name ?? `Category ${product.category_id}`}</Badge>
          </div>
          <h1 className="mb-4">{product.name}</h1>

          <div className="mb-6 flex items-center gap-4">
            <div className="flex items-center gap-1">
              {[...Array(5)].map((_, i) => (
                <Star key={i} className={`size-4 ${i < 0 ? "fill-warning text-warning" : "text-muted"}`} />
              ))}
            </div>
            <span>0</span>
            <span className="text-muted-foreground">(0 reviews)</span>
          </div>

          <div className="mb-6 flex items-baseline gap-3">
            <span className="text-3xl font-bold">{formatPrice(Number(product.price ?? 0))}</span>
          </div>

          <p className="mb-6 text-muted-foreground">{product.description}</p>

          <div className="mb-6">
            <label className="mb-3 block">Quantity</label>
            <div className="flex w-fit items-center gap-2 rounded-lg border">
              <Button variant="ghost" size="icon" onClick={() => setQuantity(Math.max(1, quantity - 1))}>
                <Minus className="size-4" />
              </Button>
              <span className="w-12 text-center">{quantity}</span>
              <Button variant="ghost" size="icon" onClick={() => setQuantity(Math.min(Number(product.stock ?? 1), quantity + 1))}>
                <Plus className="size-4" />
              </Button>
            </div>
          </div>

          <Button className="mb-8 w-full" size="lg" onClick={handleAddToCart} disabled={!inStock}>
            {inStock ? "Add to Cart" : "Out of Stock"}
          </Button>

          <div className="space-y-3 border-t pt-6">
            <div className="flex items-center gap-3 text-sm">
              <Truck className="size-5 text-primary" />
              <span>Free shipping on orders over $50</span>
            </div>
            <div className="flex items-center gap-3 text-sm">
              <RotateCcw className="size-5 text-primary" />
              <span>30-day return policy</span>
            </div>
            <div className="flex items-center gap-3 text-sm">
              <Shield className="size-5 text-primary" />
              <span>Secure checkout</span>
            </div>
          </div>
        </div>
      </div>

      <Tabs defaultValue="description" className="mb-16">
        <TabsList>
          <TabsTrigger value="description">Description</TabsTrigger>
          <TabsTrigger value="reviews">Reviews (0)</TabsTrigger>
        </TabsList>
        <TabsContent value="description" className="mt-6">
          <p className="text-muted-foreground">{product.description}</p>
        </TabsContent>
        <TabsContent value="reviews" className="mt-6">
          <p className="text-muted-foreground">Customer reviews will be displayed here.</p>
        </TabsContent>
      </Tabs>

      {relatedProducts.length > 0 && (
        <div>
          <h2 className="mb-6">Related Products</h2>
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
            {relatedProducts.map((item) => (
              <ProductCard key={item.id} product={toProductCardItem(item)} />
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
