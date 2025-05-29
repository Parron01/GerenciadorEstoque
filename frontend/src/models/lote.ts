export interface Lote {
  id: string;
  productId: string; // Provided by backend as product_id
  quantity: number;
  dataValidade: string; // YYYY-MM-DD
  createdAt?: string;
  updatedAt?: string;
}

export interface LotePayload {
  quantity: number;
  dataValidade: string; // YYYY-MM-DD, ensure this format is sent
}

// Mirrored from backend's LoteChangeDetail struct
export interface LoteChangeDetails {
  loteId?: string; // This is often the EntityID from the history record
  productId?: string;
  action: "created" | "updated" | "deleted" | string; // Allow other actions if any
  quantityBefore?: number;
  quantityAfter?: number;
  quantityChanged?: number;
  dataValidade?: string; // Current/New for create/delete actions
  dataValidadeOld?: string;
  dataValidadeNew?: string;
}
