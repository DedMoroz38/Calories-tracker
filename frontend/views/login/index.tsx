"use client";

import { useEffect, useRef } from "react";
import { Card } from "@/shared/ui";
import { useAuth } from "@/features/auth/model/useAuth";
import { RouteGuard } from "@/features/auth/ui/RouteGuard";
import styles from "./login.module.css";

/** Telegram Login Widget injects a `window.onTelegramAuth` callback. */
declare global {
  interface Window {
    onTelegramAuth: (user: Record<string, unknown>) => void;
  }
}

function LoginScreen() {
  const { error, telegramLogin } = useAuth();
  const telegramRef = useRef<HTMLDivElement>(null);

  /* Mount the official Telegram Login Widget script */
  useEffect(() => {
    const container = telegramRef.current;
    if (!container) return;

    // The widget callback is global; (re)assigning it is idempotent.
    window.onTelegramAuth = (user) => {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      telegramLogin(user as any);
    };

    // Inject exactly once. telegram-widget.js executes only on the first load
    // of its <script> tag and keeps a global init guard, so wiping + re-appending
    // (e.g. under React Strict Mode's double-invoke) leaves the widget blank.
    if (container.querySelector("script, iframe")) return;

    const script = document.createElement("script");
    script.src = "https://telegram.org/js/telegram-widget.js?22";
    script.async = true;
    script.setAttribute("data-telegram-login", "MaxyCalorieTrackerBot");
    script.setAttribute("data-size", "large");
    script.setAttribute("data-onauth", "onTelegramAuth(user)");
    script.setAttribute("data-request-access", "write");
    script.setAttribute("data-userpic", "false");
    container.appendChild(script);
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <main className={styles.root}>
      <div className={styles.content}>
        {/* Logo / hero */}
        <div className={styles.hero}>
          <div className={styles.logoMark}>M</div>
          <h1 className={styles.title}>Maxy Sports</h1>
          <p className={styles.sub}>Track calories and macros with your Telegram account</p>
        </div>

        <Card className={styles.card}>
          <div className={styles.telegramWrap}>
            <div ref={telegramRef} />
          </div>

          {error && <p className={styles.error}>{error}</p>}
        </Card>

        <p className={styles.footer}>
          By signing in you agree to the&nbsp;
          <a href="#" className={styles.link}>Terms of Service</a>
        </p>
      </div>
    </main>
  );
}

export function LoginPage() {
  return (
    /* Redirect to /home if already authenticated */
    <RouteGuard redirectIfAuthed>
      <LoginScreen />
    </RouteGuard>
  );
}
