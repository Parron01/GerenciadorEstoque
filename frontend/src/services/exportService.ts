import * as XLSX from "xlsx";
import type { Product } from "@/models/product";

interface ExcelRow {
  Tipo: "Produto" | "Lote";
  "Nome Produto": string;
  "Unidade Produto"?: string;
  "Qtd Total Produto"?: number;
  "Qtd Lote"?: number;
  "Validade Lote"?: string;
}

function calculateProductTotalQuantity(product: Product): number {
  if (product.lotes && product.lotes.length > 0) {
    return product.lotes.reduce((sum, lote) => sum + lote.quantity, 0);
  }
  return product.quantity;
}

function formatDate(dateString?: string): string {
  if (!dateString) return "";
  try {
    if (dateString.match(/^\d{4}-\d{2}-\d{2}$/)) {
      const parts = dateString.split("-");
      return new Date(
        parseInt(parts[0]),
        parseInt(parts[1]) - 1,
        parseInt(parts[2])
      ).toLocaleDateString();
    }
    return new Date(dateString).toLocaleDateString();
  } catch (e) {
    return dateString;
  }
}

export function exportProductsToExcel(
  products: Product[],
  fileName: string = "produtos_estoque.xlsx"
): void {
  // Create HTML table in memory (this approach often preserves more formatting)
  const table = document.createElement("table");
  const thead = document.createElement("thead");
  const tbody = document.createElement("tbody");

  // Create header row with styling
  const headerRow = document.createElement("tr");
  [
    "Tipo",
    "Nome Produto",
    "Unidade Produto",
    "Qtd Total Produto",
    "Qtd Lote",
    "Validade Lote",
  ].forEach((text) => {
    const th = document.createElement("th");
    th.textContent = text;
    th.style.backgroundColor = "#3B82F6"; // Indigo-500
    th.style.color = "white";
    th.style.fontWeight = "bold";
    th.style.textAlign = "center";
    th.style.padding = "8px";
    headerRow.appendChild(th);
  });

  thead.appendChild(headerRow);
  table.appendChild(thead);

  // Create data rows
  products.forEach((product) => {
    const totalQuantity = calculateProductTotalQuantity(product);

    // Product row
    const productRow = document.createElement("tr");
    productRow.style.backgroundColor = "#DBEAFE"; // Indigo-100

    // Add cells for product row
    const cells = [
      { text: "Produto", align: "center" },
      { text: product.name, align: "center" },
      { text: product.unit, align: "center" },
      { text: totalQuantity.toString(), align: "center" },
      { text: "", align: "center" },
      { text: "", align: "center" },
    ];

    cells.forEach((cell) => {
      const td = document.createElement("td");
      td.textContent = cell.text;
      td.style.textAlign = cell.align;
      td.style.padding = "6px";
      productRow.appendChild(td);
    });

    tbody.appendChild(productRow);

    // Add lote rows
    if (product.lotes && product.lotes.length > 0) {
      product.lotes.forEach((lote) => {
        const loteRow = document.createElement("tr");
        loteRow.style.backgroundColor = "#D1FAE5"; // Emerald-100

        const loteCells = [
          { text: "Lote", align: "center" },
          { text: "", align: "center" },
          { text: "", align: "center" },
          { text: "", align: "center" },
          { text: lote.quantity.toString(), align: "center" },
          { text: formatDate(lote.dataValidade), align: "center" },
        ];

        loteCells.forEach((cell) => {
          const td = document.createElement("td");
          td.textContent = cell.text;
          td.style.textAlign = cell.align;
          td.style.padding = "6px";
          loteRow.appendChild(td);
        });

        tbody.appendChild(loteRow);
      });
    }
  });

  table.appendChild(tbody);

  // Add the table to the document temporarily to let SheetJS parse it
  table.style.display = "none";
  document.body.appendChild(table);

  // Convert HTML table to worksheet
  const worksheet = XLSX.utils.table_to_sheet(table);

  // Auto-adjust columns width
  const maxWidths = [6, 30, 15, 15, 15, 20]; // Approximate widths for columns
  worksheet["!cols"] = maxWidths.map((w) => ({ wch: w }));

  // Remove the temporary table from the DOM
  document.body.removeChild(table);

  // Create workbook and export
  const workbook = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(workbook, worksheet, "Produtos");

  XLSX.writeFile(workbook, fileName, { bookType: "xlsx" });
}
