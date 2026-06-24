"use client";

import { useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { loginWithTelegram, getMe, type TelegramAuthPayload } from "@/entities/user/api/authRepository";
import { saveToken, isTokenValid } from "@/entities/user/model/session";

interface AuthState {
  loading: boolean;
  error:   string | null;
}

export function useAuth() {
  const router = useRouter();
  const [state, setState] = useState<AuthState>({ loading: false, error: null });

  const handleSuccess = useCallback(
    async (token: string) => {
      saveToken(token);
      const user = await getMe();
      if (user.onboarded) {
        router.replace("/home");
      } else {
        router.replace("/onboarding");
      }
    },
    [router],
  );

  const telegramLogin = useCallback(
    async (payload: TelegramAuthPayload) => {
      setState({ loading: true, error: null });
      try {
        const { token } = await loginWithTelegram(payload);
        await handleSuccess(token);
      } catch (err) {
        setState({ loading: false, error: (err as Error).message });
      }
    },
    [handleSuccess],
  );

  return {
    ...state,
    telegramLogin,
    isAuthenticated: isTokenValid,
  };
}
