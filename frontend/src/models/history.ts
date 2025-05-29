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
  entityType: "product" | "lote";
  entityId: string;
  details: ProductChange | LoteChangeDetails | Record<string, any>; // Parsed JSON from ChangeDetails
  createdAt: string; // ISO string
  batchId: string;
  productNameContext?: string; // Optional: For lotes, the name of the product they belong to
}

// For sending a batch of history entries (less used now, but kept for completeness)
export interface HistoryBatchInput {
  entityType: "product" | "lote";
  entityId: string;
  changes: Record<string, any>; // The raw changes object
  // BatchID will be assigned by the backend for this type of input
}

// Represents a group of history records for a single batch operation (from /api/history/grouped)
export interface HistoryBatchGroup {
  batchId: string;
  createdAt: string; // Timestamp of the first entry in the batch, for ordering
  records: ParsedHistoryRecord[];
  recordCount: number;
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
