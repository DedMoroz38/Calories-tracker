"use client";

import { useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { updateProfile, type UpdateProfilePayload } from "@/entities/user/api/profileRepository";

export type Direction = "lose" | "maintain" | "gain";

interface OnboardingState {
  step: 1 | 2 | 3;
  // Step 1
  currentWeight: string;
  goalWeight:    string;
  direction:     Direction;
  // Step 2
  calorieGoal: string;
  // Step 3 (optional)
  carbsGoal:   string;
  fatGoal:     string;
  proteinGoal: string;
  // Submit state
  loading: boolean;
  error:   string | null;
}

const DEFAULT: OnboardingState = {
  step: 1,
  currentWeight: "",
  goalWeight:    "",
  direction:     "lose",
  calorieGoal:  "",
  carbsGoal:    "",
  fatGoal:      "",
  proteinGoal:  "",
  loading: false,
  error:   null,
};

export function useOnboarding() {
  const router = useRouter();
  const [state, setState] = useState<OnboardingState>(DEFAULT);

  const set = useCallback(<K extends keyof OnboardingState>(key: K, value: OnboardingState[K]) => {
    setState((prev) => ({ ...prev, [key]: value }));
  }, []);

  const nextStep = useCallback(() => {
    setState((prev) => ({
      ...prev,
      step: Math.min(3, prev.step + 1) as 1 | 2 | 3,
      error: null,
    }));
  }, []);

  const prevStep = useCallback(() => {
    setState((prev) => ({
      ...prev,
      step: Math.max(1, prev.step - 1) as 1 | 2 | 3,
    }));
  }, []);

  const submit = useCallback(async () => {
    setState((prev) => ({ ...prev, loading: true, error: null }));
    try {
      const payload: UpdateProfilePayload = {
        calorie_goal:   Number(state.calorieGoal),
        current_weight: state.currentWeight ? Number(state.currentWeight) : undefined,
        goal_weight:    state.goalWeight    ? Number(state.goalWeight)    : undefined,
        direction:      state.direction,
        carbs_goal:     state.carbsGoal   ? Number(state.carbsGoal)   : undefined,
        fat_goal:       state.fatGoal     ? Number(state.fatGoal)     : undefined,
        protein_goal:   state.proteinGoal ? Number(state.proteinGoal) : undefined,
      };
      await updateProfile(payload);
      router.replace("/home");
    } catch (err) {
      setState((prev) => ({ ...prev, loading: false, error: (err as Error).message }));
    }
  }, [state, router]);

  return { state, set, nextStep, prevStep, submit };
}
