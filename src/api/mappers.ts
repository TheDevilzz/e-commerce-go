import { ASSET_BASE_URL } from "./client";
import type { ProductApiItem } from "./product";
import type { ProductCardItem } from "../components/ProductCard";

export const toProductCardItem = (product: ProductApiItem): ProductCardItem => ({
  id: String(product.id),
  name: product.name,
  price: Number(product.price ?? 0),
  image: product.image ? `${ASSET_BASE_URL}${product.image}` : "https://placehold.co/600x600?text=Product",
  category: product.category?.name ?? `Category ${product.category_id ?? "-"}`,
  rating: 0,
  reviews: 0,
  inStock: Number(product.stock ?? 0) > 0,
});
