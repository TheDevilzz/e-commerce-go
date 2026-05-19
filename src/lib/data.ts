export interface Product {
  id: string;
  name: string;
  price: number;
  originalPrice?: number;
  image: string;
  category: string;
  rating: number;
  reviews: number;
  description: string;
  sizes?: string[];
  colors?: string[];
  inStock: boolean;
}

export interface CartItem {
  product: Product;
  quantity: number;
  size?: string;
  color?: string;
}

export interface Order {
  id: string;
  date: string;
  total: number;
  status: "pending" | "processing" | "shipped" | "delivered" | "cancelled";
  items: CartItem[];
}

export interface User {
  id: string;
  name: string;
  email: string;
  role: "customer" | "admin";
}

export const mockProducts: Product[] = [
  {
    id: "1",
    name: "Classic White Sneakers",
    price: 89.99,
    originalPrice: 129.99,
    image: "https://images.unsplash.com/photo-1549298916-b41d501d3772?w=500&h=500&fit=crop",
    category: "Footwear",
    rating: 4.5,
    reviews: 128,
    description: "Comfortable and stylish white sneakers perfect for everyday wear.",
    sizes: ["US 7", "US 8", "US 9", "US 10", "US 11"],
    colors: ["White", "Black", "Navy"],
    inStock: true,
  },
  {
    id: "2",
    name: "Leather Laptop Bag",
    price: 149.99,
    image: "https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=500&h=500&fit=crop",
    category: "Accessories",
    rating: 4.8,
    reviews: 89,
    description: "Premium leather laptop bag with multiple compartments.",
    colors: ["Brown", "Black"],
    inStock: true,
  },
  {
    id: "3",
    name: "Wireless Headphones",
    price: 199.99,
    originalPrice: 249.99,
    image: "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=500&h=500&fit=crop",
    category: "Electronics",
    rating: 4.7,
    reviews: 256,
    description: "High-quality wireless headphones with noise cancellation.",
    colors: ["Black", "Silver", "Rose Gold"],
    inStock: true,
  },
  {
    id: "4",
    name: "Minimalist Watch",
    price: 249.99,
    image: "https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=500&h=500&fit=crop",
    category: "Accessories",
    rating: 4.6,
    reviews: 92,
    description: "Elegant minimalist watch with leather strap.",
    colors: ["Black", "Brown"],
    inStock: true,
  },
  {
    id: "5",
    name: "Canvas Backpack",
    price: 69.99,
    originalPrice: 89.99,
    image: "https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=500&h=500&fit=crop",
    category: "Bags",
    rating: 4.4,
    reviews: 156,
    description: "Durable canvas backpack perfect for daily commute.",
    colors: ["Olive", "Navy", "Black"],
    inStock: true,
  },
  {
    id: "6",
    name: "Running Shoes",
    price: 119.99,
    image: "https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=500&h=500&fit=crop",
    category: "Footwear",
    rating: 4.9,
    reviews: 412,
    description: "Lightweight running shoes with excellent cushioning.",
    sizes: ["US 7", "US 8", "US 9", "US 10", "US 11", "US 12"],
    colors: ["Black/Red", "White/Blue", "Gray"],
    inStock: true,
  },
  {
    id: "7",
    name: "Sunglasses",
    price: 159.99,
    image: "https://images.unsplash.com/photo-1572635196237-14b3f281503f?w=500&h=500&fit=crop",
    category: "Accessories",
    rating: 4.3,
    reviews: 67,
    description: "Polarized sunglasses with UV protection.",
    inStock: false,
  },
  {
    id: "8",
    name: "Casual T-Shirt",
    price: 29.99,
    image: "https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=500&h=500&fit=crop",
    category: "Clothing",
    rating: 4.2,
    reviews: 234,
    description: "Comfortable cotton t-shirt in various colors.",
    sizes: ["S", "M", "L", "XL", "XXL"],
    colors: ["White", "Black", "Navy", "Gray", "Red"],
    inStock: true,
  },
];

export const mockOrders: Order[] = [
  {
    id: "ORD-001",
    date: "2026-05-05",
    total: 339.97,
    status: "delivered",
    items: [
      { product: mockProducts[0], quantity: 2, size: "US 9", color: "White" },
      { product: mockProducts[3], quantity: 1, color: "Black" },
    ],
  },
  {
    id: "ORD-002",
    date: "2026-05-08",
    total: 199.99,
    status: "shipped",
    items: [{ product: mockProducts[2], quantity: 1, color: "Black" }],
  },
  {
    id: "ORD-003",
    date: "2026-05-10",
    total: 119.99,
    status: "processing",
    items: [{ product: mockProducts[5], quantity: 1, size: "US 10" }],
  },
];
