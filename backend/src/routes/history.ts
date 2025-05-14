import express from "express";
import { authenticateToken } from "../middleware/auth.js";
import {
  getAllHistory,
  createHistoryEntry,
} from "../controllers/historyController.js";

const router = express.Router();

// Rotas protegidas
router.get("/", authenticateToken, getAllHistory);
router.post("/", authenticateToken, createHistoryEntry);

export default router;
