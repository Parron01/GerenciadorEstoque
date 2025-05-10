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
  // Tabela de produtos
  db.exec(`
    CREATE TABLE IF NOT EXISTS products (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      unit TEXT NOT NULL CHECK(unit IN ('L', 'kg')),
      quantity REAL NOT NULL DEFAULT 0
    );
  `)

  // Tabela de histórico
  db.exec(`
    CREATE TABLE IF NOT EXISTS history (
      id TEXT PRIMARY KEY,
      date TEXT NOT NULL,
      changes TEXT NOT NULL
    );
  `)
}

initDatabase()

export default db
