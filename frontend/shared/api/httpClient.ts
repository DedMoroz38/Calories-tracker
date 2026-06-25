/**
 * shared/api/httpClient.ts
 *
 * One shared base HTTP client used by every repository in the app.
 * - Reads the Bearer token from localStorage on each request.
 * - Sends/receives JSON.
 * - Unwraps success envelope `{ data: T }`.
 * - Throws an `ApiError` on HTTP errors, extracting `{ message }` from body.
 */

import { API_V1, TOKEN_KEY } from "@/shared/config";

// ─── Types ────────────────────────────────────────────────────────────────

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

type Method = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

interface RequestOptions {
  method?: Method;
  body?: unknown;
  /** Additional headers (merged on top of defaults) */
  headers?: Record<string, string>;
  /** Skip prepending API_V1 (for absolute URLs) */
  raw?: boolean;
  /** Skip adding the Bearer token header */
  noAuth?: boolean;
}

// ─── Token helpers ────────────────────────────────────────────────────────

/** Returns the stored JWT, or null when running on the server. */
export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
  if (typeof window !== "undefined") {
    localStorage.setItem(TOKEN_KEY, token);
  }
}

export function clearToken(): void {
  if (typeof window !== "undefined") {
    localStorage.removeItem(TOKEN_KEY);
  }
}

// ─── Core fetch wrapper ───────────────────────────────────────────────────

/**
 * Makes a JSON request against the API and returns the unwrapped `data` field.
 *
 * @param path  URL path (relative to API_V1, or an absolute URL when `raw` is true)
 * @param opts  Method, body, headers, flags
 */
export async function request<T>(
  path: string,
  opts: RequestOptions = {},
): Promise<T> {
  const { method = "GET", body, headers: extraHeaders = {}, raw = false, noAuth = false } = opts;

  const url = raw ? path : `${API_V1}${path}`;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...extraHeaders,
  };

  if (!noAuth) {
    const token = getToken();
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }
  }

  const res = await fetch(url, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  // Try to parse the body regardless of status so we can surface `message`.
  let payload: Record<string, unknown> = {};
  try {
    payload = await res.json();
  } catch {
    // Non-JSON response — leave payload empty.
  }

  if (!res.ok) {
    const message =
      typeof payload["message"] === "string"
        ? payload["message"]
        : `HTTP ${res.status}`;
    throw new ApiError(res.status, message);
  }

  // Unwrap success envelope: { data: T }
  if ("data" in payload) {
    return payload["data"] as T;
  }

  // Fallback: return the whole payload if there's no `data` wrapper.
  return payload as unknown as T;
}

// ─── Multipart upload ─────────────────────────────────────────────────────

/**
 * Uploads a single file as multipart/form-data under the field name `file`.
 * Unlike `request`, the Content-Type header is left unset so the browser adds
 * the correct multipart boundary. Returns the unwrapped `data` field.
 */
export async function uploadFile<T>(path: string, file: File): Promise<T> {
  const form = new FormData();
  form.append("file", file);

  const headers: Record<string, string> = {};
  const token = getToken();
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(`${API_V1}${path}`, {
    method: "POST",
    headers,
    body: form,
  });

  let payload: Record<string, unknown> = {};
  try {
    payload = await res.json();
  } catch {
    // Non-JSON response — leave payload empty.
  }

  if (!res.ok) {
    const message =
      typeof payload["message"] === "string" ? payload["message"] : `HTTP ${res.status}`;
    throw new ApiError(res.status, message);
  }

  return "data" in payload ? (payload["data"] as T) : (payload as unknown as T);
}

// ─── Convenience aliases ──────────────────────────────────────────────────

export const http = {
  get: <T>(path: string, opts?: Omit<RequestOptions, "method" | "body">) =>
    request<T>(path, { ...opts, method: "GET" }),

  post: <T>(path: string, body?: unknown, opts?: Omit<RequestOptions, "method" | "body">) =>
    request<T>(path, { ...opts, method: "POST", body }),

  put: <T>(path: string, body?: unknown, opts?: Omit<RequestOptions, "method" | "body">) =>
    request<T>(path, { ...opts, method: "PUT", body }),

  patch: <T>(path: string, body?: unknown, opts?: Omit<RequestOptions, "method" | "body">) =>
    request<T>(path, { ...opts, method: "PATCH", body }),

  delete: <T>(path: string, opts?: Omit<RequestOptions, "method" | "body">) =>
    request<T>(path, { ...opts, method: "DELETE" }),
};
