/**
 * entities/user/api/authRepository.ts
 *
 * All API calls related to authentication.
 * Uses the shared httpClient — no raw fetch here.
 */

import { http } from "@/shared/api/httpClient";
import type { User } from "../model/types";

export interface TelegramAuthPayload {
  id:         number;
  first_name: string;
  last_name:  string;
  username:   string;
  photo_url:  string;
  auth_date:  number;
  hash:       string;
}

export interface AuthResponse {
  token: string;
}

/** POST /auth/telegram — verify Telegram Login Widget payload */
export async function loginWithTelegram(payload: TelegramAuthPayload): Promise<AuthResponse> {
  return http.post<AuthResponse>("/auth/telegram", payload, { noAuth: true });
}

/** GET /auth/me — current user info including `onboarded` flag */
export async function getMe(): Promise<User> {
  return http.get<User>("/auth/me");
}
