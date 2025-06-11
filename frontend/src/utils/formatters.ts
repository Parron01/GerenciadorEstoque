/**
 * Formats technical action names to user-friendly text
 */
export function formatActionName(action: string): string {
  const actionMap: Record<string, string> = {
    created: "Criado",
    deleted: "Removido",
    updated: "Atualizado",
    product_details_updated: "Detalhes atualizados",
    quantity_changed: "Quantidade alterada",
    lote_created: "Lote criado",
    lote_updated: "Lote atualizado",
    lote_deleted: "Lote removido",
    add: "Adicionado",
    remove: "Removido",
    update: "Atualizado",
  };

  return actionMap[action] || action.replace(/_/g, " ");
}

/**
 * Truncates UUIDs to a more readable format
 */
export function formatId(id: string, maxLength: number = 6): string {
  if (!id || id.length <= maxLength) return id;
  return id.substring(0, maxLength);
}

/**
 * Formats a date string to a localized date and time
 */
export function formatDateTime(dateStr: string): string {
  try {
    return new Date(dateStr).toLocaleString("pt-BR");
  } catch (e) {
    return dateStr || "-";
  }
}

/**
 * Formats a date string specifically to DD/MM/YYYY format
 */
export function formatDateOnly(dateStr: string): string {
  try {
    if (!dateStr) return "-";

    // Handle ISO date strings (YYYY-MM-DD or YYYY-MM-DDTHH:mm:ssZ)
    const date = new Date(dateStr);

    // Check if valid date
    if (isNaN(date.getTime())) return dateStr;

    // Format to DD/MM/YYYY using Portuguese locale
    return date.toLocaleDateString("pt-BR", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
    });
  } catch (e) {
    return dateStr || "-";
  }
}

/**
 * Returns appropriate CSS classes for action types
 */
export function getActionColorClass(action: string): string {
  const positiveActions = ["created", "add", "lote_created"];
  const negativeActions = ["deleted", "remove", "lote_deleted"];
  const updateActions = [
    "updated",
    "update",
    "product_details_updated",
    "lote_updated",
  ];

  if (positiveActions.includes(action)) return "text-emerald-600";
  if (negativeActions.includes(action)) return "text-red-600";
  if (updateActions.includes(action)) return "text-amber-600";

  return "text-gray-600";
}

/**
 * Returns appropriate badge classes for action types
 */
export function getActionBadgeClass(action: string): string {
  const positiveActions = ["created", "add", "lote_created"];
  const negativeActions = ["deleted", "remove", "lote_deleted"];
  const updateActions = [
    "updated",
    "update",
    "product_details_updated",
    "lote_updated",
  ];

  if (positiveActions.includes(action))
    return "bg-emerald-100 text-emerald-800";
  if (negativeActions.includes(action)) return "bg-red-100 text-red-800";
  if (updateActions.includes(action)) return "bg-amber-100 text-amber-800";

  return "bg-gray-100 text-gray-800";
}

/**
 * Returns appropriate CSS classes for quantity changes
 */
export function getQuantityChangeClass(
  change: number | undefined | null
): string {
  if (change === undefined || change === null) return "text-gray-600";
  return change > 0
    ? "text-emerald-600"
    : change < 0
      ? "text-red-600"
      : "text-gray-600";
}

/**
 * Formats quantity changes with appropriate prefix
 */
export function formatQuantityChange(
  change: number | undefined | null
): string {
  if (change === undefined || change === null) return "";
  return change > 0 ? `+${change}` : `${change}`;
}
