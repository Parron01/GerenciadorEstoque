export interface Product {
  id: string
  name: string
  unit: 'L' | 'kg'
  quantity: number
}

export interface ProductChange {
  productId: string
  productName: string // Store the product name directly in the history record
  action: 'add' | 'remove'
  quantityChanged: number
  quantityBefore: number
  quantityAfter: number
  isNewProduct?: boolean // Flag to indicate a new product was added
  isProductRemoval?: boolean // Flag to indicate a product was completely removed
}

export interface ProductHistory {
  id: string
  date: string
  changes: ProductChange[]
}
