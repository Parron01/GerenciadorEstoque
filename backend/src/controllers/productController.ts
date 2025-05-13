import { Request, Response } from 'express'
import db from '../config/database.ts'
import { Product } from '../models/types.ts'

export const getAllProducts = (req: Request, res: Response) => {
  try {
    const products = db.prepare('SELECT * FROM products').all()
    res.json(products)
  } catch (error) {
    console.error('Erro ao buscar produtos:', error)
    res.status(500).json({ message: 'Erro ao buscar produtos' })
  }
}

export const getProductById = (req: Request, res: Response): void => {
  try {
    const { id } = req.params
    const product = db.prepare('SELECT * FROM products WHERE id = ?').get(id)

    if (!product) {
      res.status(404).json({ message: 'Produto não encontrado' })
      return // Ensure the function returns void
    }

    res.json(product)
    return // Explicitly return void
  } catch (error) {
    console.error('Erro ao buscar produto:', error)
    res.status(500).json({ message: 'Erro ao buscar produto' })
    return // Explicitly return void
  }
}

export const createProduct = (req: Request, res: Response) => {
  const product: Product = req.body

  try {
    const stmt = db.prepare('INSERT INTO products (id, name, unit, quantity) VALUES (?, ?, ?, ?)')

    stmt.run(product.id, product.name, product.unit, product.quantity)

    res.status(201).json({
      message: 'Produto adicionado com sucesso',
      id: product.id,
    })
  } catch (error) {
    console.error('Erro ao adicionar produto:', error)
    res.status(500).json({ message: 'Erro ao adicionar produto' })
  }
}

export const updateProduct = (req: Request, res: Response): void => {
  const { id } = req.params
  const { quantity } = req.body

  try {
    const stmt = db.prepare('UPDATE products SET quantity = ? WHERE id = ?')
    const info = stmt.run(quantity, id)

    if (info.changes === 0) {
      res.status(404).json({ message: 'Produto não encontrado' })
      return // Ensure the function returns void
    }

    res.json({ message: 'Produto atualizado com sucesso' })
    return // Explicitly return void
  } catch (error) {
    console.error('Erro ao atualizar produto:', error)
    res.status(500).json({ message: 'Erro ao atualizar produto' })
    return // Explicitly return void
  }
}

export const deleteProduct = (req: Request, res: Response): void => {
  const { id } = req.params

  try {
    const stmt = db.prepare('DELETE FROM products WHERE id = ?')
    const info = stmt.run(id)

    if (info.changes === 0) {
      res.status(404).json({ message: 'Produto não encontrado' })
      return // Ensure the function returns void
    }

    res.json({ message: 'Produto removido com sucesso' })
    return // Explicitly return void
  } catch (error) {
    console.error('Erro ao remover produto:', error)
    res.status(500).json({ message: 'Erro ao remover produto' })
    return // Explicitly return void
  }
}
