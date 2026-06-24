import { http } from "@/shared/api/httpClient";
import type { DayState } from "@/shared/ui";

export interface Macros {
  calories: number;
  carbs:    number;
  fat:      number;
  protein:  number;
}

export interface WeekDay {
  date:   string;
  letter: string;
  state:  DayState;
  /** True for today's dot once the calorie goal is met — shows the streak flame. */
  fire:   boolean;
}

export interface Summary {
  date:   string;
  totals: Macros;
  goals:  Macros;
  weight: { current: number; goal: number };
  streak: number;
  /** True when today's calorie goal is met, so the streak already counts today. */
  streak_active_today: boolean;
  week:   WeekDay[];
}

/**
 * GET /summary. Sends the client's timezone offset so the streak is anchored to
 * the user's local 3am-to-3am day.
 */
export async function getSummary(date?: string): Promise<Summary> {
  const params = new URLSearchParams({ tz: String(new Date().getTimezoneOffset()) });
  if (date) params.set("date", date);
  return http.get<Summary>(`/summary?${params.toString()}`);
}
