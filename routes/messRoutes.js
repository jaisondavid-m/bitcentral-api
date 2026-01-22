import express from 'express';
import { BOYS_DATA, GIRLS_DATA, HostelType } from '../data/messData.js';
import { formatDateToDataKey, getActiveMeal } from '../utils/dateUtils.js';

const router = express.Router();

router.get('/', (req, res) => {
  const { hostel, date } = req.query;

  if (!hostel || !date) {
    return res.status(400).json({ message: 'hostel and date required' });
  }

  const selectedDate = new Date(date);
  const dateKey = formatDateToDataKey(selectedDate);

  const data =
    hostel === HostelType.BOYS ? BOYS_DATA : GIRLS_DATA;

  const messData = data.mess.find(d => d.date === dateKey);

  res.json({
    date: dateKey,
    hostel,
    activeMeal: getActiveMeal(new Date()),
    data: messData || null
  });
});

export default router;
