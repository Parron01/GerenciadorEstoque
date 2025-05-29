import type { ProductChange } from "./product";
import type { LoteChangeDetails } from "./lote";

// This represents the raw record from GET /api/history
export interface BackendHistoryRecord {
  ID: string;
  EntityType: "product" | "lote";
  EntityID: string;
  ChangeDetails: string; // JSON string from backend
  CreatedAt: string;
}

// This is the parsed version for frontend use
export interface ParsedHistoryRecord {
  id: string; // Mapped from ID
  entityType: "product" | "lote";
  entityId: string; // Mapped from EntityID
  // Parsed ChangeDetails from JSON string
  details: ProductChange | LoteChangeDetails | Record<string, any>;
  createdAt: string; // Mapped from CreatedAt
  // Optional: Denormalized product name, useful for displaying lote history
  productNameContext?: string;
}
