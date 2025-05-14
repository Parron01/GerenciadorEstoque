import fs from "fs";
import path from "path";
import db from "../config/database.js";

export function createBackup() {
  // Criar diretório de backup se não existir
  const backupDir = path.resolve(__dirname, "../backups");
  if (!fs.existsSync(backupDir)) {
    fs.mkdirSync(backupDir, { recursive: true });
  }

  // Gerar nome do arquivo de backup com data atual
  const date = new Date();
  const timestamp = `${date.getFullYear()}-${(date.getMonth() + 1)
    .toString()
    .padStart(2, "0")}-${date.getDate().toString().padStart(2, "0")}`;
  const backupPath = path.join(
    backupDir,
    `inventory-backup-${timestamp}.sqlite`
  );

  // Caminho do banco de dados fonte
  const dbSource = path.resolve(__dirname, "../database/inventory.sqlite");

  try {
    // Fazer backup (fechar a conexão temporariamente)
    const isOpen = db.open;
    if (isOpen) db.close();

    fs.copyFileSync(dbSource, backupPath);

    // Reabrir a conexão se estava aberta antes
    if (isOpen) db.pragma("journal_mode = WAL"); // Reopen logic or ensure connection is active

    console.log(`Backup criado com sucesso: ${backupPath}`);

    // Limpar backups antigos
    cleanupOldBackups(backupDir);

    return backupPath;
  } catch (error) {
    console.error("Erro ao criar backup:", error);

    // Garantir que a conexão seja reaberta em caso de erro
    if (!db.open) {
      db.pragma("journal_mode = WAL");
    }

    throw error;
  }
}

function cleanupOldBackups(backupDir: string) {
  // Manter apenas os 10 backups mais recentes
  const files = fs
    .readdirSync(backupDir)
    .filter((file) => file.startsWith("inventory-backup-"))
    .map((file) => ({
      name: file,
      path: path.join(backupDir, file),
      mtime: fs.statSync(path.join(backupDir, file)).mtime.getTime(),
    }))
    .sort((a, b) => b.mtime - a.mtime);

  if (files.length > 10) {
    files.slice(10).forEach((file) => {
      fs.unlinkSync(file.path);
      console.log(`Backup antigo removido: ${file.name}`);
    });
  }
}
