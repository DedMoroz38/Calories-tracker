"use client";

/**
 * Thin client-side session store — reads/writes the JWT from localStorage
 * and keeps a decoded payload in module-level state so we don't parse on
 * every call.  No React context; components pull from here directly.
 */

import { getToken, clearToken, setToken } from "@/shared/api/httpClient";

interface TokenPayload {
  sub:  string;   // user id as string
  exp:  number;   // unix timestamp
  iat:  number;
}

function decodeJwt(token: string): TokenPayload | null {
  try {
    const [, payload] = token.split(".");
    const json = atob(payload.replace(/-/g, "+").replace(/_/g, "/"));
    return JSON.parse(json) as TokenPayload;
  } catch {
    return null;
  }
}

export function isTokenValid(): boolean {
  const token = getToken();
  if (!token) return false;
  const payload = decodeJwt(token);
  if (!payload) return false;
  // 10-second grace window
  return payload.exp * 1000 > Date.now() - 10_000;
}

export function saveToken(token: string): void {
  setToken(token);
}

export function destroySession(): void {
  clearToken();
}
