import { http } from "@/shared/api/httpClient";

export interface FoodEntry {
  id:         number;
  name:       string;
  calories:   number;
  carbs:      number;
  fat:        number;
  protein:    number;
  created_at: string;
}

export interface RecentDish {
  name:     string;
  calories: number;
  protein:  number;
  carbs:    number;
  fat:      number;
}

export interface CreateFoodPayload {
  name:      string;
  calories:  number;
  carbs?:    number;
  fat?:      number;
  protein?:  number;
}

/** A half-open time window [from, to) of absolute instants. */
export interface TimeRange {
  from: Date;
  to:   Date;
}

/**
 * The local calendar day containing `d` (defaults to now), expressed as a
 * [from, to) window of absolute instants. Boundaries are local midnight, so
 * "today" is whatever the user's device considers today — the server stays
 * timezone-agnostic and just filters on the instants we send.
 */
export function localDayRange(d: Date = new Date()): TimeRange {
  const from = new Date(d.getFullYear(), d.getMonth(), d.getDate());
  const to   = new Date(from);
  to.setDate(to.getDate() + 1);
  return { from, to };
}

/** GET /foods?from=<RFC3339>&to=<RFC3339> — omit range to get all entries. */
export async function getFoods(range?: TimeRange): Promise<FoodEntry[]> {
  const qs = range
    ? `?from=${encodeURIComponent(range.from.toISOString())}&to=${encodeURIComponent(range.to.toISOString())}`
    : "";
  return http.get<FoodEntry[]>(`/foods${qs}`);
}

/** POST /foods */
export async function createFood(payload: CreateFoodPayload): Promise<FoodEntry> {
  return http.post<FoodEntry>("/foods", payload);
}

/** DELETE /foods/:id */
export async function deleteFood(id: number): Promise<void> {
  return http.delete<void>(`/foods/${id}`);
}

/** GET /foods/recent */
export async function getRecentFoods(): Promise<RecentDish[]> {
  return http.get<RecentDish[]>("/foods/recent");
}
