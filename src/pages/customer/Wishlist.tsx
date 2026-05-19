import { useEffect, useState } from "react";
import { toast } from "sonner";

import { deleteWishlistItem, getWishlist, type WishlistItem } from "../../api/customer";
import { toProductCardItem } from "../../api/mappers";
import { ProductCard } from "../../components/ProductCard";
import { Button } from "../../components/ui/button";

export function CustomerWishlistPage() {
  const [wishlistItems, setWishlistItems] = useState<WishlistItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const loadWishlist = async () => {
    try {
      const response = await getWishlist();
      setWishlistItems(response.data);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to load wishlist");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    void Promise.resolve().then(loadWishlist);
  }, []);

  const removeWishlistItem = async (id: number) => {
    try {
      await deleteWishlistItem(id);
      toast.success("Removed from wishlist");
      await loadWishlist();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to remove wishlist item");
    }
  };

  return (
    <div>
      <h2 className="mb-6">My Wishlist</h2>

      {!isLoading && wishlistItems.length === 0 ? (
        <div className="py-12 text-center">
          <p className="text-muted-foreground">Your wishlist is empty</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {wishlistItems.map((item) => (
            <div key={item.id} className="space-y-2">
              <ProductCard product={toProductCardItem(item.product)} />
              <Button variant="outline" className="w-full" onClick={() => removeWishlistItem(item.id)}>
                Remove
              </Button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
