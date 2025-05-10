import express from 'express'
import { authenticateToken } from '../middleware/auth.ts'
import { getAllHistory, createHistoryEntry } from '../controllers/historyController.ts'

const router = express.Router()

// Rotas protegidas
router.get('/', authenticateToken, getAllHistory)
router.post('/', authenticateToken, createHistoryEntry)

export default router
