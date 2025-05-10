import express from 'express'
import { authenticateToken } from '../middleware/auth.ts'
import {
  getAllProducts,
  getProductById,
  createProduct,
  updateProduct,
  deleteProduct,
  updateBatchProducts,
} from '../controllers/productController.ts'

const router = express.Router()

// Rotas públicas
router.get('/', getAllProducts)
router.get('/:id', getProductById)

// Rotas protegidas
router.post('/', authenticateToken, createProduct)
router.put('/:id', authenticateToken, updateProduct)
router.delete('/:id', authenticateToken, deleteProduct)

// Rota para atualização em lote
router.post('/batch', authenticateToken, updateBatchProducts)

export default router
