import { useEffect, useState } from "react";
import { Plus, Search, Pencil, Trash2 } from "lucide-react";
import { toast } from "sonner";
import { ASSET_BASE_URL } from "../../api/client";
import { formatPrice } from "../../lib/utils";
import { Button } from "../../components/ui/button";
import { Input } from "../../components/ui/input";
import { Badge } from "../../components/ui/badge";
import ProductForm from "../../components/ProductForm";
import { deleteProduct, fetchProducts, updateProduct } from "../../api/product";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../../components/ui/dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../components/ui/table";
type AdminProductRow = {
  id: string;
  name: string;
  description: string;
  price: number;
  image: string;
  rawImage: string;
  stock: number;
  categoryId: number;
  category: string;
  rating: number;
  reviews: number;
  inStock: boolean;
  originalPrice?: number;
};

type ApiProduct = {
  id: number | string;
  name: string;
  description?: string;
  price?: number | string;
  stock?: number | string;
  category_id?: number | string;
  category?: {
    name?: string;
  };
  image?: string;
};

export function AdminProductsPage() {
  const [searchTerm, setSearchTerm] = useState("");
  const [products, setProducts] = useState<AdminProductRow[]>([]);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<AdminProductRow | null>(null);
  const [stockValue, setStockValue] = useState("");
  const [isLoading, setIsLoading] = useState(true);

  const loadProducts = async () => {
      try {
        const response = await fetchProducts();
        const apiProducts = Array.isArray(response) ? response : response.data;

        if (!Array.isArray(apiProducts)) {
          return;
        }

        setProducts(
          (apiProducts as ApiProduct[]).map((product) => ({
            id: String(product.id),
            name: product.name,
            description: product.description ?? "",
            price: Number(product.price ?? 0),
            image: product.image ? `${ASSET_BASE_URL}${product.image}` : "https://placehold.co/96x96?text=Product",
            rawImage: product.image ?? "",
            stock: Number(product.stock ?? 0),
            categoryId: Number(product.category_id ?? 0),
            category: product.category?.name ?? `Category ${product.category_id ?? "-"}`,
            rating: 0,
            reviews: 0,
            inStock: Number(product.stock ?? 0) > 0,
          })),
        );
      } catch (error) {
        setProducts([]);
        toast.error(error instanceof Error ? error.message : "Failed to load products");
      } finally {
        setIsLoading(false);
      }
    };

  useEffect(() => {
    let isMounted = true;

    void Promise.resolve().then(async () => {
      if (isMounted) {
        await loadProducts();
      }
    });

    return () => {
      isMounted = false;
    };
  }, []);

  const handleDeleteProduct = async (id: string) => {
    try {
      await deleteProduct(id);
      toast.success("Product deleted");
      await loadProducts();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to delete product");
    }
  };

  const handleUpdateStock = async (product: AdminProductRow) => {
    const nextStock = Number(stockValue);
    if (!Number.isFinite(nextStock) || nextStock < 0) {
      toast.error("Stock must be zero or more");
      return;
    }
    try {
      await updateProduct(product.id, {
        name: product.name,
        description: product.description,
        price: product.price,
        stock: nextStock,
        category_id: product.categoryId,
        image: product.rawImage,
      });
      toast.success("Product updated");
      setEditingProduct(null);
      await loadProducts();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to update product");
    }
  };

  const filteredProducts = products.filter((product) =>
    product.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="space-y-6 ">
      <div className="flex items-center justify-between">
        <h2>Products</h2>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button variant="outline" size="sm">
              <Plus className="size-4" />
              Add Product
            </Button>
          </DialogTrigger>
          <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-xl">
            <DialogHeader>
              <DialogTitle>Add Product</DialogTitle>
              <DialogDescription>
                Fill in product details and upload an image.
              </DialogDescription>
            </DialogHeader>
            <ProductForm
              onCreated={() => {
                setIsDialogOpen(false);
                loadProducts();
              }}
            />
          </DialogContent>
        </Dialog>
        <Dialog open={Boolean(editingProduct)} onOpenChange={(open) => !open && setEditingProduct(null)}>
          <DialogContent className="sm:max-w-sm">
            <DialogHeader>
              <DialogTitle>Update Stock</DialogTitle>
              <DialogDescription>
                Adjust available stock for {editingProduct?.name}.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <Input
                type="number"
                min={0}
                value={stockValue}
                onChange={(event) => setStockValue(event.target.value)}
              />
              <Button
                className="w-full"
                onClick={() => editingProduct && handleUpdateStock(editingProduct)}
              >
                Save Stock
              </Button>
            </div>
          </DialogContent>
        </Dialog>
      </div>

      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
        <Input
          placeholder="Search products..."
          className="pl-9"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
      </div>

      <div className="border rounded-lg ">
        <Table >
          <TableHeader>
            <TableRow>
              <TableHead>Product</TableHead>
              <TableHead>Category</TableHead>
              <TableHead>Price</TableHead>
              <TableHead>Stock</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredProducts.map((product) => (
              <TableRow key={product.id}>
                <TableCell>
                  <div className="flex items-center gap-3">
                    <img
                      src={product.image}
                      alt={product.name}
                      className="size-12 rounded-lg object-cover"
                    />
                    <div>
                      <div className="font-medium">{product.name}</div>
                      <div className="text-sm text-muted-foreground">
                        {product.reviews} reviews
                      </div>
                    </div>
                  </div>
                </TableCell>
                <TableCell>{product.category}</TableCell>
                <TableCell>
                  <div>
                    <div className="font-medium">
                      {formatPrice(product.price)}
                    </div>
                    {product.originalPrice && (
                      <div className="text-sm text-muted-foreground line-through">
                        {formatPrice(product.originalPrice)}
                      </div>
                    )}
                  </div>
                </TableCell>
                <TableCell>
                  <Badge variant={product.inStock ? "success" : "destructive"}>
                    {product.inStock ? "In Stock" : "Out of Stock"}
                  </Badge>
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex justify-end gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => {
                        setEditingProduct(product);
                        setStockValue(String(product.stock));
                      }}
                    >
                      <Pencil className="size-4" />
                    </Button>
                    <Button variant="ghost" size="icon" onClick={() => handleDeleteProduct(product.id)}>
                      <Trash2 className="size-4 text-destructive" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
            {filteredProducts.length === 0 && (
              <TableRow>
                <TableCell colSpan={5} className="py-8 text-center text-muted-foreground">
                  {isLoading ? "Loading products..." : "No products found."}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
