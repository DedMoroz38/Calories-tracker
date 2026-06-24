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
}

export interface Summary {
  date:   string;
  totals: Macros;
  goals:  Macros;
  weight: { current: number; goal: number };
  streak: number;
  week:   WeekDay[];
}

/** GET /summary?date=YYYY-MM-DD */
export async function getSummary(date?: string): Promise<Summary> {
  const qs = date ? `?date=${date}` : "";
  return http.get<Summary>(`/summary${qs}`);
}
