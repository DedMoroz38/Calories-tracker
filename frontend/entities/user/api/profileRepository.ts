/**
 * entities/user/api/profileRepository.ts
 *
 * Profile / goals endpoints.
 */

import { http } from "@/shared/api/httpClient";

export interface Profile {
  calorie_goal:  number;
  carbs_goal:    number;
  fat_goal:      number;
  protein_goal:  number;
  current_weight: number;
  goal_weight:   number;
  direction:     "lose" | "maintain" | "gain";
  onboarded:     boolean;
  first_name:    string;
  username:      string;
  avatar_url:    string;
}

export interface UpdateProfilePayload {
  calorie_goal:   number;
  carbs_goal?:    number;
  fat_goal?:      number;
  protein_goal?:  number;
  current_weight?: number;
  goal_weight?:   number;
  direction?:     "lose" | "maintain" | "gain";
}

export async function getProfile(): Promise<Profile> {
  return http.get<Profile>("/profile");
}

export async function updateProfile(payload: UpdateProfilePayload): Promise<Profile> {
  return http.put<Profile>("/profile", payload);
}
