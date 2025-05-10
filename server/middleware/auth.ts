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
  // Adicionado tipo de retorno :void
  const authHeader = req.headers['authorization']
  const token = authHeader && authHeader.split(' ')[1]

  if (!token) {
    res.status(401).json({ message: 'Token de autenticação não fornecido' })
    return // Adicionado para garantir que a função retorne void neste caminho
  }

  jwt.verify(
    token,
    process.env.JWT_SECRET || 'fallback_secret',
    (err: VerifyErrors | null, user: any) => {
      // err tipado corretamente
      if (err) {
        res.status(403).json({ message: 'Token inválido ou expirado' })
        return // Adicionado para clareza e consistência no callback
      }

      req.user = user
      next()
    },
  )
  // A função implicitamente retorna undefined (void) se o primeiro if não for atingido
  // e após jwt.verify ser chamado (já que a callback é assíncrona).
}
