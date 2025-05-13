import express from 'express'
import { loginUser, verifyToken } from '../controllers/authController.ts'

const router = express.Router()

// Add health check endpoint
router.get('/health', (req, res) => {
  res.json({ status: 'ok', message: 'Server is running' })
})

router.post('/login', loginUser)
router.get('/verify', verifyToken)

export default router
