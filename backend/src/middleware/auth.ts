import { Request, Response, NextFunction } from "express";
import jwt from "jsonwebtoken";

// Estender o tipo Request para incluir o usuário
declare global {
  namespace Express {
    interface Request {
      user?: any;
    }
  }
}

// Modificado o tipo de retorno para permitir Response ou void
export function authenticateToken(
  req: Request,
  res: Response,
  next: NextFunction
): Response | void {
  // Skip authentication in development mode
  if (process.env.NODE_ENV === "development") {
    next();
    return;
  }

  const authHeader = req.headers["authorization"];
  const token = authHeader && authHeader.split(" ")[1];

  if (!token) {
    return res
      .status(401)
      .json({ message: "Token de autenticação não fornecido" });
  }

  try {
    const jwtSecret = process.env.JWT_SECRET || "fallback_secret";
    const user = jwt.verify(token, jwtSecret);
    req.user = user;
    next();
  } catch (err) {
    return res.status(403).json({ message: "Token inválido ou expirado" });
  }
}
