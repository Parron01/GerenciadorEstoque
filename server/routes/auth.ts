import express, { RequestHandler } from 'express'
import jwt from 'jsonwebtoken' // Changed from namespace import

const router = express.Router()

const loginHandler: RequestHandler = (req, res) => {
  const { username, password } = req.body

  if (username !== process.env.ADMIN_USERNAME) {
    res.status(401).json({ message: 'Credenciais inválidas' })
    return
  }
  if (password !== process.env.ADMIN_PASSWORD) {
    res.status(401).json({ message: 'Credenciais inválidas' })
    return
  }

  const currentJwtSecret = process.env.JWT_SECRET
  if (!currentJwtSecret) {
    // This check ensures currentJwtSecret is a string in the code path below
    res.status(500).json({ message: 'JWT_SECRET não definido no .env' })
    return
  }
  // At this point, TypeScript knows currentJwtSecret is of type 'string'

  const expirationTime = process.env.JWT_EXPIRATION || '7d' // This will be a string

  // Explicitly define the type for the options object
  const tokenOptions: jwt.SignOptions = {
    expiresIn: expirationTime as import('ms').StringValue,
  }

  try {
    // Pass the correctly typed secret and options
    const token = jwt.sign({ username }, currentJwtSecret, tokenOptions)
    res.json({ token, username })
  } catch (error) {
    // This catch is for runtime errors, your current issue is compile-time
    console.error('Erro ao assinar o token JWT:', error)
    res.status(500).json({ message: 'Erro interno ao gerar o token de autenticação' })
  }
}

const verifyTokenHandler: RequestHandler = (req, res) => {
  const authHeader = req.headers.authorization
  const token = authHeader?.split(' ')[1]

  if (!token) {
    res.status(401).json({ valid: false, message: 'Token não fornecido' })
    return
  }

  const currentJwtSecret = process.env.JWT_SECRET
  if (!currentJwtSecret) {
    res.status(500).json({ message: 'JWT_SECRET não definido no .env' })
    return
  }

  jwt.verify(token, currentJwtSecret, (err, decoded) => {
    if (err) {
      res.status(401).json({ valid: false, message: 'Token inválido ou expirado' })
      return
    }
    res.json({ valid: true, user: decoded })
  })
}

router.post('/login', loginHandler)
router.get('/verify', verifyTokenHandler)

export default router
