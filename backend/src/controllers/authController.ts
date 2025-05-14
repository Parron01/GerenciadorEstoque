import { Request, Response } from "express";
import jwt, { Secret, SignOptions } from "jsonwebtoken"; // Adicionado SignOptions
import bcrypt from "bcrypt";
import db from "../config/database.js";

interface User {
  id: string;
  username: string;
  password: string;
}

export const loginUser = (req: Request, res: Response) => {
  const { username, password } = req.body;

  try {
    // Verificar se o usuário existe no banco
    const user = db
      .prepare("SELECT * FROM users WHERE username = ?")
      .get(username) as User;

    if (!user) {
      return res.status(401).json({ message: "Credenciais inválidas" });
    }

    // Verificar senha com bcrypt
    const passwordMatch = bcrypt.compareSync(password, user.password);

    if (!passwordMatch) {
      return res.status(401).json({ message: "Credenciais inválidas" });
    }

    // Garantir que o JWT_SECRET existe e é do tipo correto
    const jwtSecret = process.env.JWT_SECRET || "fallback_secret";

    // Define um valor de expiração seguro
    // 7 dias em segundos = 604800
    const options: SignOptions = {
      expiresIn: (process.env.JWT_EXPIRATION ||
        "7d") as jwt.SignOptions["expiresIn"],
    };

    const token = jwt.sign(
      { username: user.username, id: user.id },
      jwtSecret as Secret,
      options
    );

    res.json({ token, username: user.username });
  } catch (error) {
    console.error("Erro ao autenticar usuário:", error);
    res
      .status(500)
      .json({ message: "Erro interno ao processar a autenticação" });
  }
};

export const verifyToken = (req: Request, res: Response) => {
  const authHeader = req.headers.authorization;
  const token = authHeader?.split(" ")[1];

  if (!token) {
    return res
      .status(401)
      .json({ valid: false, message: "Token não fornecido" });
  }

  const jwtSecret = process.env.JWT_SECRET || "fallback_secret";

  jwt.verify(token, jwtSecret as Secret, (err, decoded) => {
    if (err) {
      return res
        .status(401)
        .json({ valid: false, message: "Token inválido ou expirado" });
    }
    res.json({ valid: true, user: decoded });
  });
};
