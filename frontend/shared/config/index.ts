/** Backend base URL, e.g. http://localhost:3000 */
export const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:3000";

/** Versioned prefix for all REST endpoints */
export const API_V1 = `${API_BASE}/api/v1`;

/** localStorage key for the JWT */
export const TOKEN_KEY = "maxy_token";
