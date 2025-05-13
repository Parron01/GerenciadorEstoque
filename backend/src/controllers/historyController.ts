import { Request, Response } from 'express'
import db from '../config/database.ts'
import { ProductHistory } from '../models/types.ts'

export const getAllHistory = (req: Request, res: Response) => {
  try {
    const historyRecords = db.prepare('SELECT * FROM history ORDER BY date DESC').all()

    // Converter as alterações de volta para objetos
    const formattedHistory = historyRecords.map((entry: any) => ({
      ...entry,
      changes: JSON.parse(entry.changes),
    }))

    res.json(formattedHistory)
  } catch (error) {
    console.error('Erro ao buscar histórico:', error)
    res.status(500).json({ message: 'Erro ao buscar histórico' })
  }
}

export const createHistoryEntry = (req: Request, res: Response) => {
  const historyEntry: ProductHistory = req.body

  try {
    // Converter changes para string JSON para armazenamento
    const stmt = db.prepare('INSERT INTO history (id, date, changes) VALUES (?, ?, ?)')

    stmt.run(historyEntry.id, historyEntry.date, JSON.stringify(historyEntry.changes))

    res.status(201).json({
      message: 'Histórico adicionado com sucesso',
      id: historyEntry.id,
    })
  } catch (error) {
    console.error('Erro ao adicionar histórico:', error)
    res.status(500).json({ message: 'Erro ao adicionar histórico' })
  }
}
