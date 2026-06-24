import { http } from "@/shared/api/httpClient";

export interface WeightEntry {
  id:          number;
  weight:      number;
  recorded_at: string;
}

/** POST /weights */
export async function logWeight(weight: number, recorded_at?: string): Promise<WeightEntry> {
  return http.post<WeightEntry>("/weights", { weight, recorded_at });
}

/** GET /weights */
export async function getWeights(): Promise<WeightEntry[]> {
  return http.get<WeightEntry[]>("/weights");
}
