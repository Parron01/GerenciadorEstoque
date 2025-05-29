import type { Lote } from "./lote";

export interface Product {
  id: string;
  name: string;
  unit: "L" | "kg";
  quantity: number;
  lotes?: Lote[]; // Added Lotes
}

export interface ProductChange {
  productId: string;
  productName: string; // Store the product name directly in the history record
  action: "add" | "remove" | "update" | "created" | "deleted" | string; // Allow more actions
  quantityChanged?: number; // Make optional as not all actions change quantity
  quantityBefore?: number; // Make optional
  quantityAfter?: number; // Make optional
  isNewProduct?: boolean; // Flag to indicate a new product was added
  isProductRemoval?: boolean; // Flag to indicate a product was completely removed
  // Optional: to store what was changed if it's an 'update' action
  changedFields?: { field: string; oldValue?: any; newValue?: any }[];
}

export interface ProductHistory {
  id: string;
  date: string;
  changes: ProductChange[];
}
