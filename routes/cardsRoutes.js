import express from "express";
import cards from "../data/cards.js";

const router = express.Router();


router.get("/", (req, res) => {
  res.status(200).json({
    success: true,
    count: cards.length,
    data: cards
  });
});

export default router;
