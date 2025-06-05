import type { ProductChange } from "./product";
import type { LoteChangeDetails } from "./lote";

// Represents the raw history record from the backend
export interface BackendHistoryRecord {
  ID: string;
  Date: string; // ISO string
  EntityType: "product" | "lote";
  EntityID: string;
  ChangeDetails: string; // JSON string
  BatchID: string;
  CreatedAt: string; // ISO string, often same as Date
}

// Represents a parsed history record for frontend use
export interface ParsedHistoryRecord {
  id: string;
  entityType: "product" | "lote" | "product_batch_context"; // Added new type
  entityId: string;
  details:
    | ProductChange
    | LoteChangeDetails
    | ProductBatchContextChangeDetails // Added new type
    | Record<string, any>;
  createdAt: string; // ISO string
  batchId: string;
  productNameContext?: string;
  productCurrentTotalQuantity?: number; // This field will now be populated from product_batch_context snapshot
}

// For sending a batch of history entries (less used now, but kept for completeness)
export interface HistoryBatchInput {
  entityType: "product" | "lote";
  entityId: string;
  changes: Record<string, any>; // The raw changes object
  // BatchID will be assigned by the backend for this type of input
}

// New: Represents summary for a product within a batch
export interface ProductSummaryForBatch {
  productId: string;
  productName: string;
  totalQuantityBeforeBatch: number;
  totalQuantityAfterBatch: number;
  netQuantityChangeInBatch: number;
}

// Represents a group of history records for a single batch operation (from /api/history/grouped)
export interface HistoryBatchGroup {
  batchId: string;
  createdAt: string; // Timestamp of the first entry in the batch, for ordering
  records: ParsedHistoryRecord[];
  recordCount: number;
  productSummaries?: Record<string, ProductSummaryForBatch>; // Key: ProductID
}

// Represents the paginated response for grouped history (from /api/history/grouped)
export interface PaginatedHistoryBatchGroups {
  groups: HistoryBatchGroup[];
  totalBatches: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

// Local storage history structure (if different, adapt as needed)
// This is what `productStore.ts` might have been saving locally.
// The new `HistoryList.vue` will primarily consume `HistoryBatchGroup`.
export interface ProductHistory {
  id: string; // Can be a batchId
  date: string;
  changes: ProductChange[]; // Array of changes within this "local batch"
  batchId?: string; // Explicit batchId
}

// New type for sending product context history
export interface ProductBatchContextPayload {
  productId: string;
  productNameSnapshot: string;
  quantityBeforeBatch: number;
  quantityAfterBatch: number;
  // batchId will be handled by the X-Operation-Batch-ID header
}

// New type for the 'details' of a product_batch_context history record
export interface ProductBatchContextChangeDetails {
  productId: string; // Can be redundant if EntityID is used, but good for explicitness
  productNameSnapshot: string;
  quantityBeforeBatch: number;
  quantityAfterBatch: number;
}
