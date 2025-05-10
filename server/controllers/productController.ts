import { Request, Response } from 'express'
import db from '../config/database.ts'
import { Product } from '../models/types.ts'
import { v4 as uuidv4 } from 'uuid'

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
      return
    }

    res.json(product)
    return
  } catch (error) {
    console.error('Erro ao buscar produto:', error)
    res.status(500).json({ message: 'Erro ao buscar produto' })
    return
  }
}

export const createProduct = (req: Request, res: Response) => {
  try {
    let product: Product = req.body

    // Se não foi fornecido um ID, gerar um novo
    if (!product.id) {
      product.id = uuidv4()
    }

    const stmt = db.prepare('INSERT INTO products (id, name, unit, quantity) VALUES (?, ?, ?, ?)')
    stmt.run(product.id, product.name, product.unit, product.quantity)

    res.status(201).json({
      message: 'Produto adicionado com sucesso',
      product,
    })
  } catch (error) {
    console.error('Erro ao adicionar produto:', error)
    res.status(500).json({ message: 'Erro ao adicionar produto' })
  }
}

export const updateProduct = (req: Request, res: Response): void => {
  const { id } = req.params
  const { quantity, name, unit } = req.body

  try {
    // Se temos apenas quantidade, atualizamos só ela, caso contrário atualizamos tudo
    let stmt
    let info

    if (quantity !== undefined && name === undefined && unit === undefined) {
      stmt = db.prepare('UPDATE products SET quantity = ? WHERE id = ?')
      info = stmt.run(quantity, id)
    } else {
      stmt = db.prepare('UPDATE products SET name = ?, unit = ?, quantity = ? WHERE id = ?')
      info = stmt.run(name || '', unit || '', quantity || 0, id)
    }

    if (info.changes === 0) {
      res.status(404).json({ message: 'Produto não encontrado' })
      return
    }

    // Busque o produto atualizado para retornar
    const updatedProduct = db.prepare('SELECT * FROM products WHERE id = ?').get(id)

    res.json({
      message: 'Produto atualizado com sucesso',
      product: updatedProduct,
    })
    return
  } catch (error) {
    console.error('Erro ao atualizar produto:', error)
    res.status(500).json({ message: 'Erro ao atualizar produto' })
    return
  }
}

export const deleteProduct = (req: Request, res: Response): void => {
  const { id } = req.params

  try {
    // Primeiro busca o produto para retornar seus dados
    const product = db.prepare('SELECT * FROM products WHERE id = ?').get(id)

    if (!product) {
      res.status(404).json({ message: 'Produto não encontrado' })
      return
    }

    const stmt = db.prepare('DELETE FROM products WHERE id = ?')
    stmt.run(id)

    res.json({
      message: 'Produto removido com sucesso',
      product,
    })
    return
  } catch (error) {
    console.error('Erro ao remover produto:', error)
    res.status(500).json({ message: 'Erro ao remover produto' })
    return
  }
}

// Nova função para atualização em lote de produtos
export const updateBatchProducts = (req: Request, res: Response): void => {
  const updates = req.body

  try {
    const stmt = db.prepare('UPDATE products SET quantity = ? WHERE id = ?')

    const results = db.transaction(() => {
      const results = []
      for (const update of updates) {
        stmt.run(update.quantity, update.id)
        results.push(update)
      }
      return results
    })()

    res.json({
      message: 'Produtos atualizados com sucesso',
      updates: results,
    })
  } catch (error) {
    console.error('Erro ao atualizar produtos em lote:', error)
    res.status(500).json({ message: 'Erro ao atualizar produtos em lote' })
  }
}
