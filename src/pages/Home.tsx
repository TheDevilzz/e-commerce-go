import { useEffect, useState } from "react";
import { Link } from "react-router";
import { ArrowRight } from "lucide-react";
import { toast } from "sonner";
import { Button } from "../components/ui/button";
import { ProductCard, type ProductCardItem } from "../components/ProductCard";
import { fetchProducts, type ProductApiItem } from "../api/product";
import { toProductCardItem } from "../api/mappers";

export function HomePage() {
  const [featuredProducts, setFeaturedProducts] = useState<ProductCardItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const loadProducts = async () => {
      try {
        const response = await fetchProducts();
        setFeaturedProducts((response.data as ProductApiItem[]).slice(0, 4).map(toProductCardItem));
      } catch (error) {
        toast.error(error instanceof Error ? error.message : "Failed to load products");
      } finally {
        setIsLoading(false);
      }
    };
    void loadProducts();
  }, []);

  return (
    <div>
      <section className="bg-gradient-to-br from-primary/10 to-primary/5 py-20">
        <div className="container mx-auto px-4">
          <div className="max-w-2xl">
            <h1 className="text-4xl md:text-5xl mb-4">
              Discover Amazing Products at Great Prices
            </h1>
            <p className="text-lg text-muted-foreground mb-8">
              Shop the latest trends and exclusive deals on quality products.
            </p>
            <div className="flex gap-4">
              <Link to="/category">
                <Button size="lg">
                  Shop Now <ArrowRight className="ml-2 size-4" />
                </Button>
              </Link>
              <Link to="/promo">
                <Button size="lg" variant="outline">
                  View Deals
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </section>

      <section className="py-16">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between mb-8">
            <h2>Featured Products</h2>
            <Link to="/category">
              <Button variant="ghost">
                View All <ArrowRight className="ml-2 size-4" />
              </Button>
            </Link>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {featuredProducts.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
          {!isLoading && featuredProducts.length === 0 && (
            <p className="py-8 text-center text-muted-foreground">No products found.</p>
          )}
        </div>
      </section>

      <section className="bg-muted py-16">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 text-center">
            <div>
              <div className="size-12 rounded-full bg-primary/10 flex items-center justify-center mx-auto mb-4">
                <svg
                  className="size-6 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 13l4 4L19 7"
                  />
                </svg>
              </div>
              <h3 className="mb-2">Free Shipping</h3>
              <p className="text-sm text-muted-foreground">
                On orders over $50
              </p>
            </div>

            <div>
              <div className="size-12 rounded-full bg-primary/10 flex items-center justify-center mx-auto mb-4">
                <svg
                  className="size-6 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
              </div>
              <h3 className="mb-2">30-Day Returns</h3>
              <p className="text-sm text-muted-foreground">
                Easy return policy
              </p>
            </div>

            <div>
              <div className="size-12 rounded-full bg-primary/10 flex items-center justify-center mx-auto mb-4">
                <svg
                  className="size-6 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                  />
                </svg>
              </div>
              <h3 className="mb-2">Secure Payment</h3>
              <p className="text-sm text-muted-foreground">
                100% secure transactions
              </p>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}
