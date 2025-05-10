import express from 'express'
import { authenticateToken } from '../middleware/auth.ts'
import {
  getAllProducts,
  getProductById,
  createProduct,
  updateProduct,
  deleteProduct,
} from '../controllers/productController.ts'

const router = express.Router()

// Rotas públicas
router.get('/', getAllProducts)
router.get('/:id', getProductById)

// Rotas protegidas
router.post('/', authenticateToken, createProduct)
router.put('/:id', authenticateToken, updateProduct)
router.delete('/:id', authenticateToken, deleteProduct)

export default router
