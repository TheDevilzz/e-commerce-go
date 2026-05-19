import { API_BASE_URL, apiFetch, getAuthHeaders, getErrorMessage } from "./client";

type ApiResponse<T> = {
    message: string;
    data: T;
};

export type ProductApiItem = {
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

export type CategoryApiItem = {
    id: number | string;
    name: string;
};

const uploadProductImage = async (image: File, authHeaders: Record<string, string>) => {
    const formData = new FormData();
    formData.append("image", image);

    const response = await fetch(`${API_BASE_URL}/upload`, {
        method: "POST",
        headers: authHeaders,
        body: formData,
    });

    if (!response.ok) {
        throw new Error(await getErrorMessage(response, "Failed to upload product image"));
    }

    const data = await response.json();

    if (!data.url) {
        throw new Error("Upload response did not include an image URL");
    }

    return data.url as string;
};

export const fetchProducts = async () => {
    try {
        return await apiFetch<ApiResponse<ProductApiItem[]>>("/products");
    } catch (error) {
        console.error("Error fetching products:", error);
        throw error;
    }
};

export const fetchProductById = async (id: string) => {
    try {
        return await apiFetch<ApiResponse<ProductApiItem>>(`/products/${id}`);
    }
    catch (error) {
        console.error(`Error fetching product with id ${id}:`, error);
        throw error;
    }
};

export const fetchCategories = async () => {
    try {
        return await apiFetch<ApiResponse<CategoryApiItem[]>>("/categories");
    } catch (error) {
        console.error("Error fetching categories:", error);
        throw error;
    }
};

export const fetchProductsByCategory = async (categoryId: string) => {
    try {
        return await apiFetch<ApiResponse<ProductApiItem[]>>(`/categories/${categoryId}/products`);
    } catch (error) {
        console.error(`Error fetching products for category ${categoryId}:`, error);
        throw error;
    }
};

export const searchProducts = async (query: string) => {
    try {
        return await apiFetch<ApiResponse<ProductApiItem[]>>(`/products/search?q=${encodeURIComponent(query)}`);
    } catch (error) {
        console.error(`Error searching products for query "${query}":`, error);
        throw error;
    }
};

export const fetchFeaturedaddProducts = async (
  name: string,
  description: string,
  price: number,
  stock: number,
  categoryId: number,
  image: File
) => {
  try {
    const authHeaders = getAuthHeaders();
    const imageUrl = await uploadProductImage(image, authHeaders);

    const response = await fetch(`${API_BASE_URL}/products`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...authHeaders,
      },
      body: JSON.stringify({
        name,
        description,
        price: Number(price),
        stock: Number(stock),
        category_id: Number(categoryId),
        image: imageUrl,
      }),
    });

    if (!response.ok) {
      throw new Error(
        await getErrorMessage(
          response,
          "Failed to add product"
        )
      );
    }

    return await response.json();
  } catch (error) {
    console.error(
      "Error adding product:",
      error
    );

    throw error;
  }
};

export const getCategory = async () => {
    try {
        return await apiFetch<ApiResponse<CategoryApiItem[]>>("/categories");
    } catch (error) {
        console.error("Error fetching categories:", error);
        throw error;
    }
};

export const createCategory = async (name: string) => {
    try {
        const response = await fetch(`${API_BASE_URL}/categories`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                ...getAuthHeaders(),
            },
            body: JSON.stringify({ name })
        });
        if (!response.ok) {
            throw new Error(await getErrorMessage(response, "Failed to create category"));
        }
        const data = await response.json();
        return data;
    } catch (error) {
        console.error("Error creating category:", error);
        throw error;
    }
};

export const createProduct = fetchFeaturedaddProducts;

export const updateProduct = async (
    id: string,
    payload: {
        name: string;
        description: string;
        price: number;
        stock: number;
        category_id: number;
        image: string;
    },
) => {
    try {
        return await apiFetch<ApiResponse<ProductApiItem>>(`/products/${id}`, {
            method: "PUT",
            auth: true,
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
        });
    } catch (error) {
        console.error("Error updating product:", error);
        throw error;
    }
};

export const deleteProduct = async (id: string) => {
    try {
        return await apiFetch(`/products/${id}`, {
            method: "DELETE",
            auth: true,
        });
    } catch (error) {
        console.error("Error deleting product:", error);
        throw error;
    }
};

export const updateCategory = async (id: string, name: string) => {
    try {
        return await apiFetch(`/categories/${id}`, {
            method: "PUT",
            auth: true,
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ name }),
        });
    } catch (error) {
        console.error("Error updating category:", error);
        throw error;
    }
};

export const deleteCategory = async (id: string) => {
    try {
        return await apiFetch(`/categories/${id}`, {
            method: "DELETE",
            auth: true,
        });
    } catch (error) {
        console.error("Error deleting category:", error);
        throw error;
    }
};
