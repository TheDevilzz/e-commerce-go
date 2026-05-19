import { Link } from "react-router";
import { Star, Heart } from "lucide-react";
import { formatPrice } from "../lib/utils";
import { Card } from "./ui/card";
import { Button } from "./ui/button";
import { Badge } from "./ui/badge";

export type ProductCardItem = {
  id: string;
  name: string;
  price: number;
  image: string;
  category: string;
  rating?: number;
  reviews?: number;
  inStock?: boolean;
  originalPrice?: number;
};

interface ProductCardProps {
  product: ProductCardItem;
}

export function ProductCard({ product }: ProductCardProps) {
  const discount = product.originalPrice
    ? Math.round(((product.originalPrice - product.price) / product.originalPrice) * 100)
    : 0;

  return (
    <Card className="overflow-hidden group hover:shadow-md transition-shadow">
      <Link to={`/product/${product.id}`} className="block">
        <div className="relative aspect-square overflow-hidden bg-muted">
          <img
            src={product.image}
            alt={product.name}
            className="size-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
          {discount > 0 && (
            <Badge variant="destructive" className="absolute top-2 left-2">
              -{discount}%
            </Badge>
          )}
          {product.inStock === false && (
            <Badge variant="outline" className="absolute top-2 right-2 bg-background">
              Out of Stock
            </Badge>
          )}
          <Button
            variant="ghost"
            size="icon"
            className="absolute top-2 right-2 size-8 opacity-0 group-hover:opacity-100 transition-opacity bg-background/80 hover:bg-background"
            onClick={(e) => {
              e.preventDefault();
            }}
          >
            <Heart className="size-4" />
          </Button>
        </div>

        <div className="p-4">
          <div className="text-xs text-muted-foreground mb-1">
            {product.category}
          </div>
          <h3 className="font-medium line-clamp-1 mb-2">{product.name}</h3>

          <div className="flex items-center gap-1 mb-2">
            <Star className="size-3 fill-warning text-warning" />
            <span className="text-sm">{product.rating ?? 0}</span>
            <span className="text-xs text-muted-foreground">
              ({product.reviews ?? 0})
            </span>
          </div>

          <div className="flex items-center gap-2">
            <span className="font-semibold">
              {formatPrice(product.price)}
            </span>
            {product.originalPrice && (
              <span className="text-sm text-muted-foreground line-through">
                {formatPrice(product.originalPrice)}
              </span>
            )}
          </div>
        </div>
      </Link>
    </Card>
  );
}
