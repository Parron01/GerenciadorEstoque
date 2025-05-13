import BetterSqlite3 from 'better-sqlite3'
import path from 'path'
import fs from 'fs'
import { fileURLToPath } from 'url'
import bcrypt from 'bcrypt'
import dotenv from 'dotenv'

// Carregar variáveis de ambiente
dotenv.config()

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
  // Tabela de usuários
  db.exec(`
    CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      username TEXT UNIQUE NOT NULL,
      password TEXT NOT NULL,
      created_at TEXT DEFAULT CURRENT_TIMESTAMP
    );
  `)

  // Tabela de produtos - mantendo a estrutura atual para compatibilidade com o front-end
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

  // Verificar se já existe um usuário admin
  const adminUser = db
    .prepare('SELECT * FROM users WHERE username = ?')
    .get(process.env.ADMIN_USERNAME)

  if (!adminUser) {
    // Criar usuário admin com base no .env
    const adminUsername = process.env.ADMIN_USERNAME
    const adminPassword = process.env.ADMIN_PASSWORD

    if (adminUsername && adminPassword) {
      try {
        // Hash da senha
        const saltRounds = 10
        const passwordHash = bcrypt.hashSync(adminPassword, saltRounds)

        // Inserir usuário admin
        db.prepare('INSERT INTO users (username, password) VALUES (?, ?)').run(
          adminUsername,
          passwordHash,
        )
        console.log(`Usuário admin criado com sucesso: ${adminUsername}`)
      } catch (error) {
        console.error('Erro ao criar usuário admin:', error)
      }
    } else {
      console.error('ADMIN_USERNAME ou ADMIN_PASSWORD não definidos no arquivo .env')
    }
  }

  // Inserir dados padrão se tabela de produtos estiver vazia
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

    try {
      const insert = db.prepare(
        'INSERT INTO products (id, name, unit, quantity) VALUES (?, ?, ?, ?)',
      )

      db.transaction(() => {
        defaultProducts.forEach((product) => {
          insert.run(product.id, product.name, product.unit, product.quantity)
        })
      })()

      console.log('Dados padrão inseridos no banco de dados')
    } catch (error) {
      console.error('Erro ao inserir produtos padrão:', error)
    }
  }
}

// Inicializar banco de dados com tratamento de erro
try {
  initDatabase()
  console.log('Banco de dados inicializado com sucesso!')
} catch (error) {
  console.error('Erro ao inicializar banco de dados:', error)
}

export default db
