import express from 'express'
import cors from 'cors'
import dotenv from 'dotenv'
import path from 'path'
import cron from 'node-cron'

import { createBackup } from './utils/backup.ts'

import authRoutes from './routes/auth.ts'
import productRoutes from './routes/products.ts'
import historyRoutes from './routes/history.ts'

// Carregar variáveis de ambiente
dotenv.config()
const app = express()
const PORT = process.env.PORT || 3000

// Definir NODE_ENV como 'development' para facilitar o desenvolvimento
process.env.NODE_ENV = process.env.NODE_ENV || 'development'
console.log(`Servidor rodando em modo: ${process.env.NODE_ENV}`)

// Configurar middleware
app.use(cors())
app.use(express.json())

// Configurar rotas da API
app.use('/api/auth', authRoutes)
app.use('/api/products', productRoutes)
app.use('/api/history', historyRoutes)

// Servir arquivos estáticos em produção
if (process.env.NODE_ENV === 'production') {
  app.use(express.static(path.join(__dirname, '../dist')))

  app.get('*', (req, res) => {
    res.sendFile(path.join(__dirname, '../dist/index.html'))
  })
}

// Iniciar o servidor
app.listen(PORT, () => {
  console.log(`Servidor rodando na porta ${PORT}`)
})

// Configurar backup semanal (todo domingo às 3:00)
cron.schedule('0 3 * * 0', () => {
  console.log('Executando backup semanal...')
  createBackup()
})
