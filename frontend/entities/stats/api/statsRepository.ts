import { http } from "@/shared/api/httpClient";
import type { BarDatum, LineDatum } from "@/shared/ui";

export type StatsRange = "week" | "month" | "year";

export interface StatsData {
  range:           StatsRange;
  calories_per_day: BarDatum[];
  weight_trend:    LineDatum[];
  avg_per_day:     number;
  goal:            number;
  days_under_goal: number;
  streak:          number;
}

/** GET /stats?range=week|month|year — sends tz so the streak matches the home screen. */
export async function getStats(range: StatsRange = "week"): Promise<StatsData> {
  const tz = new Date().getTimezoneOffset();
  return http.get<StatsData>(`/stats?range=${range}&tz=${tz}`);
}
