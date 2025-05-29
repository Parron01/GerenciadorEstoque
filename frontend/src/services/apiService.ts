import type { Product } from "@/models/product";
import type { Lote, LotePayload } from "@/models/lote";
import type {
  BackendHistoryRecord,
  ParsedHistoryRecord,
} from "@/models/history";
import { useAuthStore } from "@/stores/authStore";

const getApiBaseUrl = () => {
  const authStore = useAuthStore(); // Access store instance here or pass as arg if preferred
  return authStore.API_BASE_URL;
};

const getAuthHeaders = () => {
  const authStore = useAuthStore();
  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };
  if (authStore.token) {
    headers["Authorization"] = `Bearer ${authStore.token}`;
  }
  return headers;
};

// Helper to handle API responses
async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const errorData = await response
      .json()
      .catch(() => ({ message: response.statusText }));
    throw new Error(errorData.message || `API error: ${response.status}`);
  }
  return response.json() as Promise<T>;
}

// Product Services
export async function fetchProducts(): Promise<Product[]> {
  const response = await fetch(`${getApiBaseUrl()}/api/products`, {
    headers: getAuthHeaders(),
  });
  const products = await handleResponse<Product[]>(response);
  // Backend sends product_id for lotes, map to productId
  return products.map((p) => ({
    ...p,
    lotes: p.lotes?.map((l) => ({
      ...l,
      productId: l.productId || (l as any).product_id,
    })),
  }));
}

export async function fetchProductById(productId: string): Promise<Product> {
  const response = await fetch(`${getApiBaseUrl()}/api/products/${productId}`, {
    headers: getAuthHeaders(),
  });
  const product = await handleResponse<Product>(response);
  return {
    ...product,
    lotes: product.lotes?.map((l) => ({
      ...l,
      productId: l.productId || (l as any).product_id,
    })),
  };
}

export async function createProductApi(
  productData: Omit<Product, "id" | "lotes" | "quantity"> & {
    quantity?: number;
  }
): Promise<Product> {
  const response = await fetch(`${getApiBaseUrl()}/api/products`, {
    method: "POST",
    headers: getAuthHeaders(),
    body: JSON.stringify(productData),
  });
  return handleResponse<Product>(response);
}

export async function updateProductApi(
  productId: string,
  productData: Partial<Pick<Product, "name" | "unit">>
): Promise<Product> {
  const response = await fetch(`${getApiBaseUrl()}/api/products/${productId}`, {
    method: "PUT",
    headers: getAuthHeaders(),
    body: JSON.stringify(productData),
  });
  return handleResponse<Product>(response);
}

export async function deleteProductApi(productId: string): Promise<void> {
  const response = await fetch(`${getApiBaseUrl()}/api/products/${productId}`, {
    method: "DELETE",
    headers: getAuthHeaders(),
  });
  if (!response.ok) {
    const errorData = await response
      .json()
      .catch(() => ({ message: response.statusText }));
    throw new Error(errorData.message || `API error: ${response.status}`);
  }
}

// Lote Services
export async function createLoteApi(
  productId: string,
  loteData: LotePayload
): Promise<Lote> {
  const response = await fetch(
    `${getApiBaseUrl()}/api/products/${productId}/lotes`,
    {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(loteData),
    }
  );
  const lote = await handleResponse<Lote>(response);
  return { ...lote, productId: lote.productId || (lote as any).product_id };
}

export async function updateLoteApi(
  loteId: string,
  loteData: LotePayload
): Promise<Lote> {
  const response = await fetch(`${getApiBaseUrl()}/api/lotes/${loteId}`, {
    method: "PUT",
    headers: getAuthHeaders(),
    body: JSON.stringify(loteData),
  });
  const lote = await handleResponse<Lote>(response);
  return { ...lote, productId: lote.productId || (lote as any).product_id };
}

export async function deleteLoteApi(loteId: string): Promise<void> {
  const response = await fetch(`${getApiBaseUrl()}/api/lotes/${loteId}`, {
    method: "DELETE",
    headers: getAuthHeaders(),
  });
  if (!response.ok) {
    const errorData = await response
      .json()
      .catch(() => ({ message: response.statusText }));
    throw new Error(errorData.message || `API error: ${response.status}`);
  }
}

// History Service
export async function fetchHistoryApi(): Promise<ParsedHistoryRecord[]> {
  const response = await fetch(`${getApiBaseUrl()}/api/history`, {
    headers: getAuthHeaders(),
  });
  const backendRecords = await handleResponse<BackendHistoryRecord[]>(response);

  // Parse ChangeDetails and map to ParsedHistoryRecord
  return backendRecords.map((record) => {
    let parsedDetails;
    try {
      parsedDetails = JSON.parse(record.ChangeDetails);
    } catch (e) {
      console.error(
        "Failed to parse ChangeDetails for history record:",
        record.ID,
        e
      );
      parsedDetails = { error: "Failed to parse details" };
    }
    return {
      id: record.ID,
      entityType: record.EntityType,
      entityId: record.EntityID,
      details: parsedDetails,
      createdAt: record.CreatedAt,
    };
  });
}
