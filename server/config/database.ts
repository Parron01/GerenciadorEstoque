import BetterSqlite3 from 'better-sqlite3'
import path from 'path'
import fs from 'fs'
import { fileURLToPath } from 'url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

// Garantir que o diretório database existe
const dbDir = path.resolve(__dirname, '../database')
if (!fs.existsSync(dbDir)) {
  fs.mkdirSync(dbDir, { recursive: true })
}

// Caminho para o arquivo do banco de dados
const dbPath = path.join(dbDir, 'inventory.sqlite')
const db = BetterSqlite3(dbPath)

// Inicializar as tabelas caso não existam
function initDatabase() {
  // Tabela de produtos - mantendo a estrutura atual para compatibilidade com o front-end
  db.exec(`
    CREATE TABLE IF NOT EXISTS products (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      unit TEXT NOT NULL CHECK(unit IN ('L', 'kg')),
      quantity REAL NOT NULL DEFAULT 0
    );
  `)

  // Tabela de histórico - expandida para incluir mais detalhes
  db.exec(`
    CREATE TABLE IF NOT EXISTS history (
      id TEXT PRIMARY KEY,
      date TEXT NOT NULL,
      changes TEXT NOT NULL
    );
  `)

  // Inserir dados padrão se tabela estiver vazia
  const count = db.prepare('SELECT COUNT(*) as count FROM products').get().count
  if (count === 0) {
    const defaultProducts = [
      { id: '1', name: 'Alade', unit: 'L', quantity: 210 },
      { id: '2', name: 'Curbix', unit: 'L', quantity: 71 },
      { id: '3', name: 'Magnum', unit: 'kg', quantity: 110 },
      { id: '4', name: 'Instivo', unit: 'L', quantity: 3 },
      { id: '5', name: 'Kasumin', unit: 'L', quantity: 50 },
      { id: '6', name: 'Priori', unit: 'L', quantity: 33 },
    ]

    const insert = db.prepare('INSERT INTO products (id, name, unit, quantity) VALUES (?, ?, ?, ?)')

    db.transaction(() => {
      defaultProducts.forEach((product) => {
        insert.run(product.id, product.name, product.unit, product.quantity)
      })
    })()

    console.log('Dados padrão inseridos no banco de dados')
  }
}

initDatabase()

export default db
