import type { Lote } from "./lote";

export interface Product {
  id: string;
  name: string;
  unit: "L" | "kg";
  quantity: number;
  lotes?: Lote[]; // Added Lotes
}

export interface ChangedField {
  field: string;
  oldValue?: any;
  newValue?: any;
  loteId?: string; // Added loteId for lote operations
}

export interface ProductChange {
  productId: string;
  productName: string; // Store the product name directly in the history record
  action:
    | "add"
    | "remove"
    | "update"
    | "created"
    | "deleted"
    | "product_details_updated"
    | "lote_created"
    | "lote_updated"
    | "lote_deleted"
    | string; // Allow more actions
  quantityChanged?: number; // Make optional as not all actions change quantity
  quantityBefore?: number; // Make optional
  quantityAfter?: number; // Make optional
  isNewProduct?: boolean; // Flag to indicate a new product was added
  isProductRemoval?: boolean; // Flag to indicate a product was completely removed
  changedFields?: ChangedField[]; // Updated to use the new interface
  // For lote changes, these might be part of details within changedFields or separate
  loteId?: string;
  loteData?: any; // For new/updated lote data
  originalLoteData?: any; // For old lote data in an update
}

export interface ProductHistory {
  id: string;
  date: string;
  changes: ProductChange[];
  batchId?: string; // Added batchId property
}
