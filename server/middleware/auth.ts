import { Request, Response, NextFunction } from 'express'
import jwt, { VerifyErrors } from 'jsonwebtoken'

// Estender o tipo Request para incluir o usuário
declare global {
  namespace Express {
    interface Request {
      user?: any
    }
  }
}

export function authenticateToken(req: Request, res: Response, next: NextFunction): void {
  // Skip authentication in development mode
  if (process.env.NODE_ENV === 'development') {
    next()
    return
  }

  const authHeader = req.headers['authorization']
  const token = authHeader && authHeader.split(' ')[1]

  if (!token) {
    res.status(401).json({ message: 'Token de autenticação não fornecido' })
    return
  }

  jwt.verify(
    token,
    process.env.JWT_SECRET || 'fallback_secret',
    (err: VerifyErrors | null, user: any) => {
      if (err) {
        res.status(403).json({ message: 'Token inválido ou expirado' })
        return
      }

      req.user = user
      next()
    },
  )
}
