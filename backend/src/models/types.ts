// Tipos compartilhados entre cliente e servidor
export interface Product {
  id: string
  name: string
  unit: 'L' | 'kg'
  quantity: number
}

export interface ProductChange {
  productId: string
  productName: string
  action: 'add' | 'remove'
  quantityChanged: number
  quantityBefore: number
  quantityAfter: number
  isNewProduct?: boolean
  isProductRemoval?: boolean
}

export interface ProductHistory {
  id: string
  date: string
  changes: ProductChange[]
}
