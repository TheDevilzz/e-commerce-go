import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import { fetchFeaturedaddProducts, getCategory } from "../api/product";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "./ui/select";
import { Textarea } from "./ui/textarea";

type CategoryOption = {
  id: string;
  name: string;
};

type ApiCategory = {
  id: number | string;
  name: string;
};

type ProductFormProps = {
  onCreated?: () => void;
};

function ProductForm({ onCreated }: ProductFormProps) {
  const [name, setName] = useState<string>("");
  const [description, setDescription] =
    useState<string>("");

  const [price, setPrice] =
    useState<number>(0);

  const [stock, setStock] =
    useState<number>(0);
  const [categoryId, setCategoryId] =
    useState<string>("");
  const [categories, setCategories] =
    useState<CategoryOption[]>([]);
  const [isLoadingCategories, setIsLoadingCategories] =
    useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const [image, setImage] =
    useState<File | null>(null);

  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    let isMounted = true;

    const loadCategories = async () => {
      setIsLoadingCategories(true);

      try {
        const response = await getCategory();
        const apiCategories = Array.isArray(response) ? response : response.data;

        if (!isMounted || !Array.isArray(apiCategories)) {
          return;
        }

        setCategories(
          (apiCategories as ApiCategory[]).map((category) => ({
            id: String(category.id),
            name: category.name,
          })),
        );
      } catch {
        toast.error("Failed to load categories");
      } finally {
        if (isMounted) {
          setIsLoadingCategories(false);
        }
      }
    };

    void Promise.resolve().then(loadCategories);

    return () => {
      isMounted = false;
    };
  }, []);

  const handleSubmit = async (
    e: React.FormEvent
  ) => {
    e.preventDefault();

    if (!image) {
      toast.error("Please select image");
      return;
    }

    const selectedCategoryId = Number(categoryId);
    if (!selectedCategoryId) {
      toast.error("Please select category");
      return;
    }

    setIsSubmitting(true);

    try {
        await fetchFeaturedaddProducts(
            name,
            description,
            price,
            stock,
            selectedCategoryId,
            image
        );
        toast.success("Product added successfully");
        setName("");
        setDescription("");
        setPrice(0);
        setStock(0);
        setCategoryId("");
        setImage(null);
        onCreated?.();
    } catch {
        toast.error("Failed to add product");
      } finally {
        setIsSubmitting(false);
      }
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="flex w-full flex-col gap-4"
    >
      <div className="space-y-2">
      <Label htmlFor="product-name">Product name</Label>
      <Input
        id="product-name"
        type="text"
        placeholder="Product name"
        value={name}
        onChange={(e) => setName(e.target.value)}
        required
      />
      </div>

      <div className="space-y-2">
      <Label htmlFor="product-description">Description</Label>
      <Textarea
        id="product-description"
        placeholder="Description"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        required
      />
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
      <div className="space-y-2">
      <Label htmlFor="product-price">Price</Label>
      <Input
        id="product-price"
        type="number"
        placeholder="Price"
        value={price}
        onChange={(e) => setPrice(Number(e.target.value))}
        required
        min={0.01}
        step="0.01"
      />
      </div>

      <div className="space-y-2">
      <Label htmlFor="product-stock">Stock</Label>
      <Input
        id="product-stock"
        type="number"
        placeholder="Stock"
        value={stock}
        onChange={(e) => setStock(Number(e.target.value))}
        required
        min={0}
      />
      </div>
      </div>

      <div className="space-y-2">
      <Label>Category</Label>
      <Select
        value={categoryId}
        onValueChange={setCategoryId}
        required
        disabled={isLoadingCategories}
      >
        <SelectTrigger>
          <SelectValue placeholder={isLoadingCategories ? "Loading categories..." : "Select category"} />
        </SelectTrigger>
        <SelectContent>
          {categories.map((category) => (
            <SelectItem key={category.id} value={category.id}>
              {category.name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      </div>

      <Input
        ref={inputRef}
        type="file"
        accept="image/*"
        className="hidden"
        onChange={(e) => {
          const selectedImage = e.target.files?.[0] ?? null;
          setImage(selectedImage);
        }}
      />
      
      {image ? (
        <div className="flex justify-center">
          <Button
            type="button"
            variant="outline"
            className="h-auto overflow-hidden p-0"
            onClick={() => inputRef.current?.click()}
          >
            <img
              src={URL.createObjectURL(image)}
              alt="Preview"
              className="h-48 w-72 object-cover"
            />
          </Button>
        </div>
      ) : (
        <div className="flex justify-center">
          <Button
            type="button"
            variant="outline"
            className="h-48 w-72 border-dashed text-muted-foreground"
            onClick={() => inputRef.current?.click()}
          >
            Click to select image
          </Button>
        </div>
      )}

      <Button
        type="submit"
        disabled={isSubmitting}
      >
        {isSubmitting ? "Adding..." : "Add Product"}
      </Button>
    </form>
  );
}

export default ProductForm;
