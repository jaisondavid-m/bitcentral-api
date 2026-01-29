import express from "express";
import cors from "cors";
import subjects from "./data/subjects.js";
import messRoutes from './routes/messRoutes.js';
import cardsRoutes from "./routes/cardsRoutes.js";

const app = express();
const PORT = process.env.PORT || 5000;

app.use(cors());
app.use(express.json());


app.get("/", (req, res) => {
  res.send("Subjects API running 🚀");
});

app.get("/api/subjects", (req, res) => {
  res.status(200).json({
    success: true,
    count: subjects.length,
    data: subjects
  });
});
// app.use('/api/mess', messRoutes);
app.use('/api/mess', (req, res, next) => {
  res.set('Cache-Control', 'no-store');
  next();
}, messRoutes);

app.use("/api/cards", cardsRoutes);

app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});