import { useEffect, useMemo, useState } from "react";
import { useParams } from "react-router";
import { SlidersHorizontal } from "lucide-react";
import { toast } from "sonner";

import { ProductCard, type ProductCardItem } from "../components/ProductCard";
import { fetchCategories, fetchProducts, type CategoryApiItem, type ProductApiItem } from "../api/product";
import { toProductCardItem } from "../api/mappers";
import { Button } from "../components/ui/button";
import { Label } from "../components/ui/label";
import { Slider } from "../components/ui/slider";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../components/ui/select";
import { Card } from "../components/ui/card";

export function CategoryPage() {
  const { categoryName } = useParams();
  const [selectedCategory, setSelectedCategory] = useState(categoryName || "All");
  const [categories, setCategories] = useState<CategoryApiItem[]>([]);
  const [products, setProducts] = useState<ProductCardItem[]>([]);
  const [priceRange, setPriceRange] = useState([0, 0]);
  const [maxPrice, setMaxPrice] = useState(0);
  const [sortBy, setSortBy] = useState("featured");
  const [showFilters, setShowFilters] = useState(true);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const loadData = async () => {
      try {
        const [productResponse, categoryResponse] = await Promise.all([
          fetchProducts(),
          fetchCategories(),
        ]);
        const mappedProducts = (productResponse.data as ProductApiItem[]).map(toProductCardItem);
        const nextMaxPrice = Math.max(...mappedProducts.map((product) => product.price), 0);
        setProducts(mappedProducts);
        setCategories(categoryResponse.data);
        setMaxPrice(nextMaxPrice);
        setPriceRange([0, nextMaxPrice]);
      } catch (error) {
        toast.error(error instanceof Error ? error.message : "Failed to load category data");
      } finally {
        setIsLoading(false);
      }
    };

    void loadData();
  }, []);

  const filteredProducts = useMemo(() => {
    const filtered = products.filter((product) => {
      const categoryMatch = selectedCategory === "All" || product.category === selectedCategory;
      const priceMatch = product.price >= priceRange[0] && product.price <= priceRange[1];
      return categoryMatch && priceMatch;
    });

    switch (sortBy) {
      case "price-low":
        filtered.sort((a, b) => a.price - b.price);
        break;
      case "price-high":
        filtered.sort((a, b) => b.price - a.price);
        break;
    }

    return filtered;
  }, [selectedCategory, priceRange, sortBy, products]);

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-6 flex items-center justify-between">
        <h1>Shop by Category</h1>
        <Button
          variant="outline"
          size="sm"
          onClick={() => setShowFilters(!showFilters)}
          className="lg:hidden"
        >
          <SlidersHorizontal className="mr-2 size-4" />
          Filters
        </Button>
      </div>

      <div className="flex gap-6">
        {showFilters && (
          <aside className="w-full shrink-0 lg:w-64">
            <Card className="sticky top-20 space-y-6 p-6">
              <div>
                <h3 className="mb-3">Categories</h3>
                <div className="space-y-2">
                  {["All", ...categories.map((category) => category.name)].map((category) => (
                    <Button
                      key={category}
                      type="button"
                      variant={selectedCategory === category ? "default" : "ghost"}
                      onClick={() => setSelectedCategory(category)}
                      className="w-full justify-start"
                    >
                      {category}
                    </Button>
                  ))}
                </div>
              </div>

              <div>
                <h3 className="mb-3">Price Range</h3>
                <Slider
                  value={priceRange}
                  onValueChange={setPriceRange}
                  max={maxPrice}
                  step={100}
                  className="mb-2"
                  disabled={maxPrice === 0}
                />
                <div className="flex items-center justify-between text-sm text-muted-foreground">
                  <span>${priceRange[0]}</span>
                  <span>${priceRange[1]}</span>
                </div>
              </div>

              <Button
                variant="outline"
                className="w-full"
                onClick={() => {
                  setSelectedCategory("All");
                  setPriceRange([0, maxPrice]);
                }}
              >
                Clear Filters
              </Button>
            </Card>
          </aside>
        )}

        <div className="flex-1">
          <div className="mb-6 flex items-center justify-between">
            <p className="text-muted-foreground">{filteredProducts.length} products found</p>
            <div className="flex items-center gap-2">
              <Label>Sort by:</Label>
              <Select value={sortBy} onValueChange={setSortBy}>
                <SelectTrigger className="w-40">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="featured">Featured</SelectItem>
                  <SelectItem value="price-low">Price: Low to High</SelectItem>
                  <SelectItem value="price-high">Price: High to Low</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {filteredProducts.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>

          {!isLoading && filteredProducts.length === 0 && (
            <div className="py-12 text-center">
              <p className="text-muted-foreground">No products found matching your filters.</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
