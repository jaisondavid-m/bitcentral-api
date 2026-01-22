import { MEAL_TIMINGS } from '../data/messData.js';

export const timeToMinutes = (timeStr) => {
  const [h, m] = timeStr.split(':').map(Number);
  return h * 60 + m;
};

export const formatDateToDataKey = (date) => {
  const d = String(date.getDate()).padStart(2, '0');
  const m = String(date.getMonth() + 1).padStart(2, '0');
  const y = String(date.getFullYear()).slice(-2);
  return `${d}-${m}-${y}`;
};

export const getActiveMeal = (now) => {
  const mins = now.getHours() * 60 + now.getMinutes();

  if (mins < timeToMinutes(MEAL_TIMINGS.breakfast.end)) return 'breakfast';
  if (mins < timeToMinutes(MEAL_TIMINGS.lunch.end)) return 'lunch';
  if (mins < timeToMinutes(MEAL_TIMINGS.tea.end)) return 'tea';
  if (mins < timeToMinutes(MEAL_TIMINGS.dinner.end)) return 'dinner';

  return 'next_day_breakfast';
};