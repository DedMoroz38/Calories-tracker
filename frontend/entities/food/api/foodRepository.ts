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

/** GET /foods?date=YYYY-MM-DD */
export async function getFoods(date?: string): Promise<FoodEntry[]> {
  const qs = date ? `?date=${date}` : "";
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
