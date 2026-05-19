import { useState } from "react";
import { toast } from "sonner";
import { createCategory } from "../api/product";
import { Button } from "./ui/button";
import { Input } from "./ui/input";

type CategoryFormProps = {
  onCreated?: () => void;
  initialName?: string;
  submitLabel?: string;
  onSubmitCategory?: (name: string) => Promise<void>;
};

function CategoryForm({
  onCreated,
  initialName = "",
  submitLabel = "Add Category",
  onSubmitCategory,
}: CategoryFormProps) {
  const [name, setName] = useState<string>(initialName);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (
    e: React.FormEvent
  ) => {
    e.preventDefault();

    const trimmedName = name.trim();
    if (!trimmedName) {
      toast.error("Please enter category name");
      return;
    }

    setIsSubmitting(true);

    try {
      if (onSubmitCategory) {
        await onSubmitCategory(trimmedName);
      } else {
        await createCategory(trimmedName);
      }

        toast.success(onSubmitCategory ? "Category updated successfully" : "Category added successfully");
        setName("");
        onCreated?.();

    } catch (error) {
        toast.error(error instanceof Error ? error.message : "Failed to add category");
      } finally {
        setIsSubmitting(false);
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="flex w-full flex-col gap-4"
    >
      <Input
        type="text"
        placeholder="Category name"
        value={name}
        onChange={(e) => setName(e.target.value)}
        required
      />
      <Button
        type="submit"
        disabled={isSubmitting}
      >
        {isSubmitting ? "Saving..." : submitLabel}
      </Button>
    </form>
  );
}

export default CategoryForm;
