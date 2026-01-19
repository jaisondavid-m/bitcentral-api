import express from "express";
import cors from "cors";
import subjects from "./subjects.js";

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

app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
