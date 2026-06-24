"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { isTokenValid } from "@/entities/user/model/session";

interface RouteGuardProps {
  children: React.ReactNode;
  /** When true, redirect authenticated users away (e.g. login page). */
  redirectIfAuthed?: boolean;
}

/**
 * Client-side route guard.
 * - On protected pages (default): redirects to /login when no valid token.
 * - On public pages (redirectIfAuthed=true): redirects to /home when already authed.
 */
export function RouteGuard({ children, redirectIfAuthed = false }: RouteGuardProps) {
  const router = useRouter();

  useEffect(() => {
    const authed = isTokenValid();
    if (redirectIfAuthed && authed) {
      router.replace("/home");
    } else if (!redirectIfAuthed && !authed) {
      router.replace("/login");
    }
  }, [router, redirectIfAuthed]);

  return <>{children}</>;
}
